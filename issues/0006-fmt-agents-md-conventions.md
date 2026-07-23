# AGENTS.md の言語規約違反を修正する

- Priority: Medium
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/fix-agents-md-conventions
- Polished: 2026-07-23

## 目的

AGENTS.md が定める言語規約（ログは英語、コメントは日本語、テストログは日本語、全角半角スペース）に準拠する。

## 優先度根拠

プロジェクト規約違反であり、コードベースの一貫性を損なう。

## 現状

以下の違反がある:

### ログメッセージが日本語（AGENTS.md: 「ログメッセージは全て英語にすること」）

- `main.go:176`: `logger.Warn("Sora Exporter は root ユーザーで実行されています。このエクスポーターは特権を必要としません。root で実行する必要はありません。")`

### コメントが英語（AGENTS.md: 「コメントは全て日本語にすること」）

- `main.go:79-80`: `// exporterMetricsRegistry is a separate registry for the metrics about // the exporter itself.`
- `main.go:124`: `// ServeHTTP implements http.Handler.`
- `main.go:153-154`: `// Note that we have to use h.exporterMetricsRegistry here to // use the same promhttp metrics for all expositions.`
- `collector/collector.go:16`: `// for testing`
- `collector/collector.go:70`: `// for testing`
- `collector/collector.go:80`: `// same as node expoter's node_time_seconds`（"expoter" typo も含む）

### テストのログメッセージが英語（AGENTS.md: 「テストのログメッセージは全て日本語にすること」）

- `main_test.go:366`: `t.Fatal(fmt.Errorf("The fixture file can't open %q: %w", fixture, err))`
- `main_test.go:369`: `t.Fatal("Unexpect metrics returned:", err)`（"Unexpect" typo も含む）

### 全角と半角の間に半角スペースがない

- `collector/license.go:50`: `// 期限翌月の 1 日 0 時 0 分0 秒の 1 秒前`（"分" と "0" の間にスペースなし）

## 設計方針

各違反箇所を以下の方針で修正する:

- ログメッセージ: 英語に翻訳する（例: `"Sora Exporter is running as root. This exporter does not require privileged access."`）
- コメント: 日本語に翻訳する。`// ServeHTTP implements http.Handler.` は Go の慣例的なコメントだが、AGENTS.md の規約を優先して日本語化する
- テストログ: 日本語に翻訳する（例: `t.Fatalf("フィクスチャファイル %q を開けません: %v", fixture, err)`）
- typo も併せて修正する（"expoter"→"exporter"、"Unexpect"→削除して日本語化）

## 完了条件

- 全てのログメッセージが英語である
- 全てのコメントが日本語である
- 全てのテストログメッセージが日本語である
- 全角と半角の間に半角スペースがある
- `go test -race -v .` が通る

## 解決方法

1. `main.go:176` のログメッセージを英語に変更する
2. `main.go:79-80,124,153-154` のコメントを日本語に翻訳する
3. `collector/collector.go:16,70,80` のコメントを日本語に翻訳する（"expoter" typo も修正）
4. `main_test.go:366,369` のテストログを日本語に変更する（"Unexpect" typo も修正）
5. `collector/license.go:50` の全角半角スペースを修正する
6. `go test -race -v .` で全テストの通過を確認する
