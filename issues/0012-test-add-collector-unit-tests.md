# collector パッケージにユニットテストを追加する

- Priority: Medium
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/add-collector-unit-tests
- Polished: 2026-07-23

## 目的

collector パッケージの変換ロジックに対するユニットテストを追加し、境界値・エラーパスの検証を可能にする。

## 優先度根拠

全変換ロジックがエンドツーエンドの統合テストでのみ検証されており、境界値やエラーパスのテストが不足している。致命的バグ（整数除算、0001 参照）が検出できなかった一因。

## 現状

`collector/` パッケージに `*_test.go` ファイルが一切存在しない。以下の関数・ロジックが統合テスト経由でしかテストされていない:

- `expiredAtToSecondSinceEpoch`（license.go:43）— 不正フォーマット、空文字列、境界値（`"2025-12"` → 翌年 1 月への繰り上がり）が未テスト
- `newWebhookResponseTimeHistogram`（webhook.go:80）— 空バケット、upper_bound=0、非単調増加が未テスト
- `CollectClusterNodes` の `ClusterNodeName`/`NodeName` フォールバック（cluster_node.go:97-112）— 両方空の場合が未テスト
- `ConnectionMetrics.Collect` の msec→sec 変換（connection.go:40-55）— 割り切れない値の精度が未テスト

## 設計方針

テーブル駆動テストで境界値・エラーパスを網羅する。テストファイルは対象ファイルと同じ `collector` パッケージに配置する（`collector/license_test.go` 等）。モックやスタブは使用しない（AGENTS.md 規約）。

## 完了条件

- `collector/` パッケージにユニットテストファイルが追加されている
- 上記 4 つの関数・ロジックに対する境界値・エラーパスのテストが含まれている
- `go test -race -v ./...` が通る

## 解決方法

1. `collector/license_test.go` を新設し、`expiredAtToSecondSinceEpoch` のテーブル駆動テストを追加する:
   - 正常系: `"2025-09"` → `1759276799`、`"2025-12"` → 翌年 1 月繰り上がり
   - 異常系: `""`、`"invalid"`、`"2025-13"` → エラー（0004 の修正後は error を返す）
2. `collector/webhook_test.go` を新設し、`newWebhookResponseTimeHistogram` のテストを追加する:
   - 正常系: 4 バケツの正常データ
   - 境界値: upper_bound=0、バケツ 1 個
3. `collector/cluster_node_test.go` を新設し、`CollectClusterNodes` のフォールバックのテストを追加する:
   - `ClusterNodeName` のみ、`NodeName` のみ、両方設定、両方空
4. `collector/connection_test.go` を新設し、msec→sec 変換の精度テストを追加する:
   - `372` → `0.372`、`12000` → `12.0`、`1500` → `1.5`
5. `go test -race -v ./...` で全テストの通過を確認する
