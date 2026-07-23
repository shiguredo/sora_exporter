# averageSetupTimeMsec の整数除算によりメトリクスが常に 0 になる

- Priority: High
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/fix-average-setup-time-integer-division
- Polished: 2026-07-23

## 目的

`sora_average_setup_time_seconds` メトリクスが整数除算により常に 0（または切り捨てられた値）を返すバグを修正する。

## 優先度根拠

メトリクスの値が常に誤っており、Prometheus による監視・アラートが機能しない。テストフィクスチャにもバグが焼き込まれており、テストが通ることでバグが正当化されている状態。

## 現状

`collector/connection.go:53` で msec→sec 変換が整数除算で行われている:

```go
ch <- newGauge(m.averageSetupTimeSec, float64(report.AverageSetupTimeMsec/1000))
```

`AverageSetupTimeMsec` は `int64` 型（`collector/sora_api.go:27`）のため、`372/1000 = 0` となる。正しい値は `0.372`。

一方、`collector/webhook.go:60` では `float64(report.TotalAuthWebhookResponseTimeMs) / 1000.0` と正しく浮動小数点除算を使っており、コード内で変換手法が不統一。

### テストフィクスチャへの影響

テストデータ（`main_test.go:23`）の `"average_setup_time_msec": 372` に対し、以下の 7 フィクスチャで期待値が `0`（バグに追従した値）になっている。正しい期待値は `0.372`:

- `test/maximum.metrics:19`
- `test/sora_client_enabled.metrics:19`
- `test/sora_cluster_metrics_enabled.metrics:19`
- `test/sora_connection_error_enabled.metrics:19`
- `test/sora_erlang_vm_enabled.metrics:19`
- `test/sora_up_cluster_nodes_failed.metrics:19`
- `test/sora_up_license_failed.metrics:19`

以下のフィクスチャは更新不要:

- `test/minimum.metrics:10` — 入力 `12000` ms は `12000/1000 = 12` と割り切れるため整数除算でも結果が一致する
- `test/invalid_config.metrics` — 全 API 失敗のため `sora_average_setup_time_seconds` 自体が出力されない
- `test/sora_up_stats_report_failed.metrics` — GetStatsReport 失敗のため出力されない
- `test/sora_up_cluster_nodes_only.metrics` — GetStatsReport 失敗のため出力されない

## 設計方針

`collector/webhook.go:60` の既存パターン（`float64(value) / 1000.0`）に揃え、int64 を先に float64 に変換してから浮動小数点除算を行う。

## 完了条件

- `collector/connection.go:53` が `float64(report.AverageSetupTimeMsec) / 1000.0` に修正されている
- 上記 7 フィクスチャの期待値が `0.372` に更新されている
- `go test -race -v .` が通る

## 解決方法

1. `collector/connection.go:53` を以下に修正する:
   ```go
   ch <- newGauge(m.averageSetupTimeSec, float64(report.AverageSetupTimeMsec)/1000.0)
   ```
2. 上記 7 フィクスチャの `sora_average_setup_time_seconds 0` を `sora_average_setup_time_seconds 0.372` に更新する
3. `go test -race -v .` で全テストの通過を確認する
