# collector パッケージにユニットテストを追加する

- Priority: Medium
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/add-collector-unit-tests
- Polished: {YYYY-MM-DD}

## 目的

collector パッケージの変換ロジックに対するユニットテストを追加し、境界値・エラーパスの検証を可能にする。

## 優先度根拠

全変換ロジックがエンドツーエンドの統合テストでのみ検証されており、境界値やエラーパスのテストが不足している。致命的バグ（整数除算）が検出できなかった一因。

## 現状

`collector/` パッケージに `*_test.go` ファイルが一切存在しない。以下の関数・ロジックが統合テスト経由でしかテストされていない:

- `expiredAtToSecondSinceEpoch`（license.go:42）— 不正フォーマット、空文字列、境界値（`"2025-12"` → 翌年繰り上がり）
- `newWebhookResponseTimeHistogram`（webhook.go:82）— 空バケット、upper_bound=0、非単調増加
- `CollectClusterNodes` の `ClusterNodeName`/`NodeName` フォールバック（cluster_node.go:90-99）
- `ConnectionMetrics.Collect` の各変換（connection.go:39-53）

## 完了条件

- `collector/` パッケージにユニットテストファイルが追加されている
- 上記の関数・ロジックに対する境界値・エラーパスのテストが含まれている
- 全テストが通る

## 解決方法

`collector/license_test.go`、`collector/webhook_test.go`、`collector/cluster_node_test.go` 等を新設し、テーブル駆動テストで境界値・エラーパスを網羅する。
