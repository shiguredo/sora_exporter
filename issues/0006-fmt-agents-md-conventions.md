# AGENTS.md の言語規約違反を修正する

- Priority: Medium
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/fix-agents-md-conventions
- Polished: {YYYY-MM-DD}

## 目的

AGENTS.md が定める言語規約（ログは英語、コメントは日本語、テストログは日本語、全角半角スペース）に準拠する。

## 優先度根拠

プロジェクト規約違反であり、コードベースの一貫性を損なう。

## 現状

以下の違反がある:

### ログメッセージが日本語（AGENTS.md: 「ログメッセージは全て英語にすること」）

- `main.go:176`: `logger.Warn("Sora Exporter は root ユーザーで実行されています。...")`

### コメントが英語（AGENTS.md: 「コメントは全て日本語にすること」）

- `main.go:79-80`: `// exporterMetricsRegistry is a separate registry for the metrics about // the exporter itself.`
- `main.go:124`: `// ServeHTTP implements http.Handler.`
- `main.go:153-154`: `// Note that we have to use h.exporterMetricsRegistry here to // use the same promhttp metrics for all expositions.`
- `collector/collector.go:16`: `// for testing`
- `collector/collector.go:70`: `// for testing`
- `collector/collector.go:80`: `// same as node expoter's node_time_seconds`（"expoter" typo も含む）

### テストのログメッセージが英語（AGENTS.md: 「テストのログメッセージは全て日本語にすること」）

- `main_test.go:233-234`: `t.Fatal(fmt.Errorf("The fixture file can't open %q: %w", fixture, err))`
- `main_test.go:236`: `t.Fatal("Unexpect metrics returned:", err)`（"Unexpect" typo も含む）

### 全角と半角の間に半角スペースがない

- `collector/license.go:50`: `// 期限翌月の 1 日 0 時 0 分0 秒の 1 秒前`（"分" と "0" の間にスペースなし）

## 完了条件

- 全てのログメッセージが英語である
- 全てのコメントが日本語である
- 全てのテストログメッセージが日本語である
- 全角と半角の間に半角スペースがある
- 全テストが通る

## 解決方法

各ファイルを編集し、規約に準拠するよう修正する。typo（"expoter"、"Unexpect"）も併せて修正する。
