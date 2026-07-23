# expiredAtToSecondSinceEpoch がパース失敗時に 0 を返しログも出さない

- Priority: Medium
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/fix-expired-at-parse-error
- Polished: 2026-07-23

## 目的

ライセンス期限のパース失敗時に誤ったメトリクス（1970 年失効）が出力されることを防ぎ、エラーを可視化する。

## 優先度根拠

`expired_at` が `"2024-12"` 形式でない場合、`sora_license_expired_at_timestamp_seconds` が `0`（1970-01-01T00:00:00Z）を返す。「ライセンスが 1970 年に失効した」という誤ったメトリクスが出力され、誤ったアラートを引き起こす可能性がある。

## 現状

`collector/license.go:43-53` の `expiredAtToSecondSinceEpoch` がパース失敗時に `0` を返す:

```go
func expiredAtToSecondSinceEpoch(expiredAt string) float64 {
    expiredAtTime, err := time.Parse("2006-01", expiredAt)
    if err != nil {
        return 0
    }
    // ...
}
```

呼び出し元の `LicenseMetrics.Collect`（`license.go:36-41`）は戻り値をそのままメトリクスとして出力するため、パース失敗時に `sora_license_expired_at_timestamp_seconds 0` が出力される。エラーのログ出力もない。

`LicenseMetrics` 構造体には logger フィールドがない。logger は `Collector` 構造体が持っている（`collector.go:23`）。

## 設計方針

`expiredAtToSecondSinceEpoch` の戻り値を `(float64, error)` に変更する。`LicenseMetrics.Collect` でエラーを受け取り、`sora_license_expired_at_timestamp_seconds` の出力をスキップする。ログ出力は `Collector.Collect`（`collector.go:152-154`）から `LicenseMetrics.Collect` を呼び出している箇所で行う（`LicenseMetrics` に logger を追加しない）。

## 完了条件

- `expiredAtToSecondSinceEpoch` が `(float64, error)` を返す
- パース失敗時に `sora_license_expired_at_timestamp_seconds` メトリクスの出力がスキップされる
- パース失敗時にエラーログが出力される
- `go test -race -v .` が通る

## 解決方法

1. `expiredAtToSecondSinceEpoch` の戻り値を `(float64, error)` に変更し、パース失敗時に `(0, err)` を返す
2. `LicenseMetrics.Collect` でエラーを受け取り、エラー時はメトリクスをスキップして error を返す
3. `Collector.Collect`（`collector.go:152-154`）で `LicenseMetrics.Collect` のエラーをログ出力する
4. `go test -race -v .` で全テストの通過を確認する
