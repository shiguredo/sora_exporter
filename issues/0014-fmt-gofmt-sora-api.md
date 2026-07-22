# collector/sora_api.go の gofmt 非準拠を修正する

- Priority: Low
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/fix-gofmt-sora-api
- Polished: {YYYY-MM-DD}

## 目的

`collector/sora_api.go` の gofmt 非準拠（アラインメントの不備）を修正する。

## 優先度根拠

CI の `go fmt` チェックで検出されるフォーマットの問題。

## 現状

`collector/sora_api.go:53-60` の `soraWebhookReport` 構造体の webhook レスポンスタイム関連フィールドで、フィールド名・型・構造体タグ間のアラインメントに余分なスペースがある。

## 完了条件

- `gofmt -l .` の出力に `collector/sora_api.go` が含まれない
- 全テストが通る

## 解決方法

`gofmt -w collector/sora_api.go` を実行する。
