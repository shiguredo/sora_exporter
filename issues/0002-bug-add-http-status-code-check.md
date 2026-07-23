# Sora API 呼び出しで HTTP ステータスコードがチェックされていない

- Priority: High
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/fix-http-status-code-check
- Polished: 2026-07-23

## 目的

Sora API 呼び出し時に HTTP ステータスコードを検証し、非 200 応答をエラーとして処理する。

## 優先度根拠

リバースプロキシやロードバランサーが 4xx/5xx を返しても、ボディが有効な JSON であれば正常応答として処理してしまう。誤ったメトリクスが Prometheus に登録される可能性がある。

## 現状

`collector/collector.go` の 3 つの fetch 関数すべてで `resp.StatusCode` を確認せず、直接 JSON デコードしている:

- `fetchGetStatsReport`（:176-198）
- `fetchGetLicense`（:200-222）
- `fetchListClusterNodes`（:224-246）

```go
resp, err := client.Do(req)
if err != nil {
    c.logger.Error("failed to request to Sora GetStatsReport API", "err", err)
    return nil, err
}
defer resp.Body.Close()

var report soraGetStatsReport
if err := json.NewDecoder(resp.Body).Decode(&report); err != nil {
```

### 既存テストとの関係

`TestInvalidConfig` と `TestSoraUp` は非 200 応答をテストしているが、`http.Error()` がプレーンテキストのボディを返すため、JSON デコード失敗としてエラーパスに入っている。ステータスコードチェックがなくてもテストは通る。非 200 + 有効な JSON ボディのケースはテストされていない。

## 設計方針

`client.Do(req)` の直後、`resp.Body.Close()` の前にステータスコードチェックを挿入する。3 関数とも同一のパターンで修正する。

## 完了条件

- 3 つの fetch 関数すべてで `resp.StatusCode != http.StatusOK` の場合にエラーを返す
- 非 200 + 有効な JSON ボディを返すテストケースが追加されている
- `go test -race -v .` が通る

## 解決方法

1. 各 fetch 関数の `client.Do(req)` 後に以下を追加する:
   ```go
   if resp.StatusCode != http.StatusOK {
       c.logger.Error("unexpected status code from Sora API", "status", resp.StatusCode)
       return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }
   ```
2. `main_test.go` に非 200 + 有効な JSON ボディを返すテストケースを追加する（例: `http.HandlerFunc` で `w.WriteHeader(http.StatusInternalServerError)` の後に有効な JSON を書き込む）
3. `go test -race -v .` で全テストの通過を確認する
