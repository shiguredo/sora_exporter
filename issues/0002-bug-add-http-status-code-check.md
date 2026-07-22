# Sora API 呼び出しで HTTP ステータスコードがチェックされていない

- Priority: High
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/fix-http-status-code-check
- Polished: {YYYY-MM-DD}

## 目的

Sora API 呼び出し時に HTTP ステータスコードを検証し、非 200 応答をエラーとして処理する。

## 優先度根拠

リバースプロキシやロードバランサーが 4xx/5xx を返しても、ボディが有効な JSON であれば正常応答として処理してしまう。誤ったメトリクスが Prometheus に登録される可能性がある。

## 現状

`collector/collector.go` の 3 つの fetch 関数（`fetchGetStatsReport`:118-131、`fetchGetLicense`:133-155、`fetchListClusterNodes`:157-179）すべてで `resp.StatusCode` を確認せず、直接 JSON デコードしている。

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

## 完了条件

- 3 つの fetch 関数すべてで `resp.StatusCode` が 2xx 以外の場合にエラーを返す
- 非 200 応答のテストケースが追加されている
- 全テストが通る

## 解決方法

各 fetch 関数で `client.Do(req)` の後にステータスコードチェックを追加する:

```go
if resp.StatusCode != http.StatusOK {
    c.logger.Error("unexpected status code from Sora API", "status", resp.StatusCode)
    return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
```
