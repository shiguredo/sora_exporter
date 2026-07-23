# Webhook ヒストグラムのバケットデータが検証されていない

- Priority: Medium
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/fix-webhook-histogram-validation
- Polished: 2026-07-23

## 目的

Sora API から取得したヒストグラムのバケットデータを検証し、不正なデータが黙ってエクスポートされることを防ぐ。

## 優先度根拠

不正なヒストグラムが Prometheus に登録され、ダッシュボード・アラートに悪影響が出る可能性がある。`prometheus.MustNewConstHistogram` はバケット境界の昇順・count 整合性を検証しないため、不正データがそのままエクスポートされる。

## 現状

`collector/webhook.go:80-86` の `newWebhookResponseTimeHistogram` は Sora API のバケットデータを検証せずに `prometheus.MustNewConstHistogram` に渡している:

```go
func newWebhookResponseTimeHistogram(d *prometheus.Desc, count uint64, sum float64, buckets []webhookResponseTimeBucket) prometheus.Metric {
    bucketMap := make(map[float64]uint64, len(buckets))
    for _, b := range buckets {
        bucketMap[float64(b.UpperBound)/1000.0] = uint64(b.Count)
    }
    return prometheus.MustNewConstHistogram(d, count, sum, bucketMap)
}
```

`webhookResponseTimeBucket` の `UpperBound`・`Count` は `int64` 型（`collector/sora_api.go:33-34`）。

以下のケースで不正なヒストグラムがエクスポートされる:

1. バケットの cumulative count が `count`（successful+failed）を超える場合 — session/event/stats webhook には `TotalIgnoredXxxWebhook` フィールドが存在する（`sora_api.go:49,52,55`）が、count の計算（`webhook.go:59,64,69,74`）には ignored が含まれない。Sora API が ignored をバケットに含めている場合、最大バケットの count が `count` を超える
2. バケット count が単調増加しない場合 — Prometheus のヒストグラムは cumulative count の単調増加を要求する
3. int64 の負値が uint64 変換で巨大値にラップアラウンドする場合 — `uint64(b.Count)` が約 1.8×10^19 になる

## 設計方針

`newWebhookResponseTimeHistogram` の戻り値を `(prometheus.Metric, error)` に変更し、呼び出し側（`webhook.go:58-77` の 4 箇所）でエラー時にログ出力してスキップする。`ch <- nil` はチャネルに nil を送ることになるため避ける。

## 完了条件

- バケットデータの検証（upper_bound の昇順、count の単調増加、負値チェック、最大バケット count <= count）が追加されている
- 不正なデータの場合はログ出力してメトリクスをスキップする
- `go test -race -v .` が通る

## 解決方法

1. `newWebhookResponseTimeHistogram` の戻り値を `(prometheus.Metric, error)` に変更する
2. 関数内で以下を検証する:
   - `b.UpperBound` と `b.Count` が負でないこと
   - バケット境界（`UpperBound`）が昇順であること
   - バケットの cumulative count が単調増加であること
   - 最大バケットの count が `count` 以下であること
3. 呼び出し側 4 箇所（`webhook.go:58-77`）でエラー時に `slog.Error` でログ出力してスキップする
4. `go test -race -v .` で全テストの通過を確認する
