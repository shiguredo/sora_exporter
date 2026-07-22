# sora_ongoing_connections_total の _total サフィックスを削除する

- Priority: Low
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/change-rename-ongoing-connections
- Polished: {YYYY-MM-DD}

## 目的

gauge メトリクス `sora_ongoing_connections_total` から `_total` サフィックスを削除し、Prometheus 命名規則に準拠する。

## 優先度根拠

Prometheus の命名規則では `_total` サフィックスは counter のみに使用すべき。gauge に `_total` が付いているとユーザーが counter 前提のクエリを書く誤用を誘発する。ただし破壊的変更になるため優先度は Low。

## 現状

`collector/connection.go:4`:

```go
totalOngoingConnections: newDesc("ongoing_connections_total", "The total number of ongoing connections."),
```

`collector/connection.go:50` で `newGauge` として登録されている。

## 完了条件

- メトリクス名が `sora_ongoing_connections` に変更されている
- テストフィクスチャが更新されている
- CHANGES.md に破壊的変更として記載されている
- 全テストが通る

## 解決方法

`newDesc("ongoing_connections_total", ...)` を `newDesc("ongoing_connections", ...)` に変更し、テストフィクスチャと CHANGES.md を更新する。
