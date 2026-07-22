# expiredAtToSecondSinceEpoch がパース失敗時に 0 を返しログも出さない

- Priority: Medium
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/fix-expired-at-parse-error
- Polished: {YYYY-MM-DD}

## 目的

ライセンス期限のパース失敗時に誤ったメトリクス（1970 年失効）が出力されることを防ぎ、エラーを可視化する。

## 優先度根拠

`expired_at` が `"2024-12"` 形式でない場合、`sora_license_expired_at_timestamp_seconds` が `0`（1970-01-01T00:00:00Z）を返す。「ライセンスが 1970 年に失効した」という誤ったメトリクスが出力され、誤ったアラートを引き起こす可能性がある。

## 現状

`collector/license.go:42-52`:

```go
func expiredAtToSecondSinceEpoch(expiredAt string) float64 {
    expiredAtTime, err := time.Parse("2006-01", expiredAt)
    if err != nil {
        return 0
    }
    // ...
}
```

パース失敗時に `0` を返すが、エラーのログ出力もなく、メトリクスの出力もスキップされない。

## 完了条件

- パース失敗時にエラーログが出力される
- パース失敗時に `sora_license_expired_at_timestamp_seconds` メトリクスの出力がスキップされる
- パース失敗のテストケースが追加されている
- 全テストが通る

## 解決方法

`LicenseMetrics.Collect` 内で `expiredAtToSecondSinceEpoch` がエラーを返すように変更し、エラー時はログ出力してメトリクスをスキップする。
