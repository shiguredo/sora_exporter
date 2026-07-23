# Collect メソッドのネットワーク周り改善

- Priority: Medium
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/fix-collect-network-handling
- Polished: 2026-07-23

## 目的

`Collect()` メソッドのネットワーク周りの問題を修正し、リソースリーク・ブロッキング・DoS リスクを解消する。

## 優先度根拠

高頻度スクレイプ時に FD/ソケットが蓄積するリソースリーク、Sora 遅延時の全スクレイプブロッキング、巨大レスポンスによるメモリ消費のリスクがある。

## 現状

`collector/collector.go:95-114` の `Collect()` メソッドに以下の問題が集中している:

1. **HTTP Transport/Client の新規生成**（:102-106）: `Collect()` 呼び出しごとに `http.Transport` と `http.Client` を新規生成。Transport 内部のコネクションプールとアイドルコネクション管理 goroutine が蓄積する
2. **単一 context タイムアウトの共有**（:99）: 1 つの `context.WithTimeout` を 3 回の逐次 API 呼び出しで共有。最初の呼び出しがタイムアウトの大部分を消費すると残りは即座に `context deadline exceeded` で失敗する。`--sora.timeout=5s` は各 5 秒ではなく 3 回合計 5 秒
3. **mutex 保持中のネットワーク I/O**（:96-97）: `c.mutex.Lock()` が Collect 冒頭で取得され、3 回の HTTP リクエストを含む関数終了時まで保持される。Sora が遅延している場合、並行するスクレイプリクエストがすべてブロックされる
4. **レスポンスボディのサイズ制限なし**（:193, :215, :237）: `json.NewDecoder(resp.Body).Decode(&report)` がレスポンスボディを制限なく読み込む

## 設計方針

- HTTP Transport/Client は `NewCollector` 時に 1 回だけ生成し、`Collector` のフィールドとして保持する。`skipSslVerify` の値は `NewCollector` 時に確定するため、Transport の再生成は不要
- context は各 fetch 関数内で個別に `context.WithTimeout` を生成する。`--sora.timeout` の意味が「各 API 呼び出しのタイムアウト」に変わるため、CHANGES.md への記載が必要
- mutex は `Collect` 内のフィールド変更がないため、ネットワーク I/O 前に解放する。具体的には fetch 呼び出しを mutex の外に出し、メトリクス出力部分のみ mutex で保護する
- レスポンスボディは `io.LimitReader` で制限する。上限値は定数として定義する（例: 10MB）

## 完了条件

- HTTP Transport/Client が `NewCollector` 時に 1 回だけ生成され再利用される
- 各 API 呼び出しに個別の context が生成される
- mutex がネットワーク I/O の外で管理される
- レスポンスボディに `io.LimitReader` によるサイズ制限がある
- CHANGES.md にタイムアウトの意味変更が記載されている
- `go test -race -v .` が通る

## 解決方法

1. `Collector` 構造体に `client *http.Client` フィールドを追加し、`NewCollector` で生成する
2. 各 fetch 関数の引数から `client *http.Client` を削除し、`c.client` を使用する
3. 各 fetch 関数内で `context.WithTimeout(context.Background(), c.timeout)` を生成する
4. `Collect` 内の fetch 呼び出しを mutex の外に出す
5. 各 fetch 関数で `resp.Body = http.MaxBytesReader` または `io.LimitReader` によるサイズ制限を追加する
6. `go test -race -v .` で全テストの通過を確認する
