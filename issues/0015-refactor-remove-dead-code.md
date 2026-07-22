# 死にコード・コメントアウトされたコードを削除する

- Priority: Low
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/refactor-remove-dead-code
- Polished: {YYYY-MM-DD}

## 目的

コードベース内の死にコード・コメントアウトされたコード・重複コメントを削除し、可読性を向上する。

## 優先度根拠

コードの可読性・保守性の問題。動作には影響しない。

## 現状

### コメントアウトされたコード（5 箇所・計 20 行超）

- `collector/erlang_vm.go:30-33`（var ブロック）
- `collector/erlang_vm.go:62-66`（struct、`XXX(tnamao)` コメント含む）
- `collector/erlang_vm.go:89-92`（Describe）
- `collector/erlang_vm.go:119-122`（Collect）
- `collector/sora_api.go:148-151`（struct フィールド）

array 型の Erlang VM 統計に「対応しない」という判断は確定しており、Git 履歴に存在する情報。

### 未使用の HTTPClient interface

- `collector/collector.go:63-65`: どこからも参照されていない

### 重複コメント

- `// この統計情報はアンドキュメントです` が `main.go:47,52,57`、`collector/client.go:4`、`collector/connection_error.go:5`、`collector/erlang_vm.go:5` の 6 箇所に重複

## 完了条件

- コメントアウトされたコードが全て削除されている
- 未使用の `HTTPClient` interface が削除されている
- 重複コメントが整理されている
- 全テストが通る

## 解決方法

各ファイルを編集して削除する。
