# コード設計の改善

- Priority: Low
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/refactor-code-design
- Polished: {YYYY-MM-DD}

## 目的

コードの設計上の問題を改善し、保守性・可読性を向上する。

## 優先度根拠

設計の問題であり、動作には影響しない。今後の機能追加時の保守性を改善する。

## 現状

### fetch 系メソッドの重複コード

`collector/collector.go:155-210` の 3 メソッド（`fetchGetStatsReport`、`fetchGetLicense`、`fetchListClusterNodes`）が同一パターンを繰り返す。ジェネリック関数に統合可能。

### newHandler の引数が 11 個

`main.go:83-86` で bool が 5 個連続する位置引数。設定構造体の導入で解決可能。

### handler 構造体のフィールドが構築後に死んでいる

`main.go:73-88` の handler 構造体は `soraMetricsHandler` 以外の全フィールドが `innerHandler()` 呼び出し以降参照されない。`http.Handler` を直接返す関数に簡素化可能。

### soraClusterUp の所属と振る舞いの不一致

`collector/cluster_node.go:10` で `SoraClusterMetrics` 構造体に定義されているが、`Describe()`/`Collect()` では親の `Collector` が直接管理している。

### freezeTimeSeconds のテスト専用コード混入

`collector/collector.go:17-19`、`main.go:82,89,155`。`main()` では常に `false` がハードコードされている。

### その他の軽微な改善

- `freezedTimeSeconds` → `frozenTimeSeconds`（命名誤り、collector.go:17）
- `clusterRelaies` → `clusterRelays`（typo、cluster_node.go:100）
- `Collector.URI` → `uri`（不要なエクスポート、collector.go:25）
- `fetchListClusterNodes` の `*[]soraClusterNode` → `[]soraClusterNode`（不要なポインタ、collector.go:193）
- `sync.RWMutex` → `sync.Mutex`（RLock 未使用、collector.go:22）
- デコード専用構造体の無意味な `omitempty` タグ（sora_api.go:6-10）

## 完了条件

- 上記の設計改善が適用されている
- 全テストが通る

## 解決方法

各項目を個別に修正する。破壊的変更を伴うもの（`Collector.URI` の非公開化等）は注意深く行う。
