# Webhook ヒストグラムのバケットデータが検証されていない

- Priority: Medium
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/fix-webhook-histogram-validation
- Polished: {YYYY-MM-DD}

## 目的

Sora API から取得したヒストグラムのバケットデータを検証し、不正なデータが黙ってエクスポートされることを防ぐ。

## 優先度根拠

panic はしない（第 2 周レビューで確認済み）が、不正なヒストグラムが Prometheus に登録され、ダッシュボード・アラートに悪影響が出る可能性がある。

## 現状

`collector/webhook.go:82-86` の `newWebhookResponseTimeHistogram` は Sora API のバケットデータを検証せずに `prometheus.MustNewConstHistogram` に渡している。

以下のケースで不正なヒストグラムがエクスポートされる:

1. バケットの cumulative count が `count`（successful+failed）を超える場合（ignored webhook がバケットに含まれる場合）
2. バケット count が単調増加しない場合
3. int64 の負値が uint64 変換で巨大値にラップアラウンドする場合

```go
func newWebhookResponseTimeHistogram(d *prometheus.Desc, count uint64, sum float64, buckets []webhookResponseTimeBucket) prometheus.Metric {
    bucketMap := make(map[float64]uint64, len(buckets))
    for _, b := range buckets {
        bucketMap[float64(b.UpperBound)/1000.0] = uint64(b.Count)
    }
    return prometheus.MustNewConstHistogram(d, count, sum, bucketMap)
}
```

## 完了条件

- バケットデータの検証（単調増加、count 整合性、負値チェック）が追加されている
- 不正なデータの場合はログ出力してメトリクスをスキップする
- 全テストが通る

## 解決方法

`newWebhookResponseTimeHistogram` 内でバケットデータを検証し、不正な場合は `nil` を返して呼び出し側でスキップする。または `prometheus.NewConstHistogram`（error を返す版）を使用する。
