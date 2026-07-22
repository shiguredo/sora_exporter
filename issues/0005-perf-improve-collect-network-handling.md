# Collect メソッドのネットワーク周り改善

- Priority: Medium
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/refactor-collect-network-handling
- Polished: {YYYY-MM-DD}

## 目的

`Collect()` メソッドのネットワーク周りの問題を修正し、リソースリーク・ブロッキング・DoS リスクを解消する。

## 優先度根拠

高頻度スクレイプ時に FD/ソケットが蓄積するリソースリーク、Sora 遅延時の全スクレイプブロッキング、巨大レスポンスによるメモリ消費のリスクがある。

## 現状

`collector/collector.go:96-114` に以下の問題が集中している:

1. **HTTP Transport/Client の新規生成**（:101-105）: `Collect()` 呼び出しごとに `http.Transport` と `http.Client` を新規生成。Transport 内部のコネクションプールとアイドルコネクション管理 goroutine が蓄積する
2. **単一 context タイムアウトの共有**（:99）: 1 つの `context.WithTimeout` を 3 回の逐次 API 呼び出しで共有。最初の呼び出しがタイムアウトの大部分を消費すると残りは即座に失敗する
3. **mutex 保持中のネットワーク I/O**（:96-97）: `c.mutex.Lock()` が Collect 冒頭で取得され、3 回の HTTP リクエストを含む関数終了時まで保持される
4. **レスポンスボディのサイズ制限なし**（:127 等）: `json.NewDecoder(resp.Body).Decode(&report)` がレスポンスボディを制限なく読み込む

## 完了条件

- HTTP Transport/Client が `NewCollector` 時に 1 回だけ生成され再利用される
- 各 API 呼び出しに個別の context が生成される（またはドキュメントで合計タイムアウトであることが明示される）
- mutex がネットワーク I/O の外で管理される
- レスポンスボディに `io.LimitReader` によるサイズ制限がある
- 全テストが通る

## 解決方法

1. `http.Transport` と `http.Client` を `Collector` のフィールドとして `NewCollector` 時に生成する
2. 各 fetch 関数で個別の `context.WithTimeout` を生成する
3. mutex のスコープをメトリクス出力部分に限定する
4. `io.LimitReader(resp.Body, maxSize)` で読み取りサイズを制限する
