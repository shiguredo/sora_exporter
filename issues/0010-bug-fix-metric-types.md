# Erlang VM・Raft メトリクスの gauge/counter 不整合を修正する

- Priority: Medium
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/fix-metric-types
- Polished: {YYYY-MM-DD}

## 目的

Prometheus の gauge/counter の意味論に合わないメトリクスの型を修正する。

## 優先度根拠

メトリクスの型が誤っていると、Prometheus の `rate()`/`increase()` 関数が正しく動作せず、監視・アラートの精度が低下する。

## 現状

### gauge だが counter が妥当（4 件）

以下は VM 起動以来の累積値（単調増加）だが gauge として登録されている:

- `collector/erlang_vm.go:124`: `erlangVMGarbageCollectionNumberOfGcs` — 累積 GC 回数
- `collector/erlang_vm.go:125`: `erlangVMGarbageCollectionWordsReclaimed` — 累積回収ワード数
- `collector/erlang_vm.go:132`: `erlangVMRuntimeTotalRunTime` — 累積実行時間
- `collector/erlang_vm.go:137`: `erlangVMWallClockTotalWallclockTime` — 累積実時間

同ファイル内の同種の累積値（`erlangVMContextSwitches`:121、`erlangVMReductionsTotalReductions`:129 等）は正しく `newCounter` を使っており不整合。

### counter だが gauge が妥当（1 件）

- `collector/cluster_node.go:119`: `raftTerm` — Raft の term は論理エポック番号（現在の選挙期）であり累積カウントではない。ヘルプテキストも "The **current** Raft term." と現在値であることを示す

## 完了条件

- 4 件の Erlang VM メトリクスが `newCounter` に変更されている
- `raftTerm` が `newGauge` に変更されている
- テストフィクスチャの `# TYPE` 行が更新されている
- 全テストが通る

## 解決方法

各メトリクスの `newGauge`/`newCounter` を変更し、テストフィクスチャの `# TYPE` 行を対応させて更新する。アンドキュメントメトリクス（`--sora.erlang-vm-metrics` フラグ有効時のみ）のため後方互換への影響は限定的。
