# averageSetupTimeMsec の整数除算によりメトリクスが常に 0 になる

- Priority: High
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/fix-average-setup-time-integer-division
- Polished: {YYYY-MM-DD}

## 目的

`sora_average_setup_time_seconds` メトリクスが整数除算により常に 0（または切り捨てられた値）を返すバグを修正する。

## 優先度根拠

メトリクスの値が常に誤っており、Prometheus による監視・アラートが機能しない。テストフィクスチャにもバグが焼き込まれており、テストが通ることでバグが正当化されている状態。

## 現状

`collector/connection.go:53` で msec→sec 変換が整数除算で行われている:

```go
ch <- newGauge(m.averageSetupTimeSec, float64(report.AverageSetupTimeMsec/1000))
```

`AverageSetupTimeMsec` は `int64` 型（`collector/sora_api.go:26`）のため、`372/1000 = 0` となる。正しい値は `0.372`。

一方、`collector/webhook.go:69` では `float64(report.TotalAuthWebhookResponseTimeMs) / 1000.0` と正しく浮動小数点除算を使っており、コード内で変換手法が不統一。

テストフィクスチャ `test/maximum.metrics:14` は `sora_average_setup_time_seconds 0` と期待しており、バグに追従している。`test/minimum.metrics` の入力 12000ms は割り切れるためバグが隠蔽されている。

## 完了条件

- `collector/connection.go:53` が `float64(report.AverageSetupTimeMsec) / 1000.0` に修正されている
- テストフィクスチャの期待値が正しい値（`0.372`）に更新されている
- 割り切れない値（例: 1500ms → 1.5s）のテストケースが追加されている
- 全テストが通る

## 解決方法

1. `collector/connection.go:53` を `float64(report.AverageSetupTimeMsec) / 1000.0` に修正する
2. `test/maximum.metrics` の `sora_average_setup_time_seconds` の期待値を `0.372` に更新する
3. `test/sora_up_license_failed.metrics`、`test/sora_up_cluster_nodes_failed.metrics` 等、同じテストデータを使うフィクスチャも同様に更新する
