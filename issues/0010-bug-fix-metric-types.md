# Erlang VM・Raft メトリクスの gauge/counter 不整合を修正する

- Priority: Medium
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/fix-metric-types
- Polished: 2026-07-23

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

- `collector/cluster_node.go:119`: `raftTerm` — Raft の term は論理エポック番号（現在の選挙期）であり累積カウントではない。ヘルプテキストも "The **current** Raft term." と現在値であることを示す。`raftCommitIndex`（:120）は単調増加かつ `increase()` で「期間中のコミットエントリ数」が得られるため counter のまま据え置く

## 設計方針

Erlang VM メトリクスはアンドキュメント（`--sora.erlang-vm-metrics` フラグ有効時のみ出力）のため、TYPE 変更の後方互換への影響は限定的。`raftTerm` の TYPE 変更は `--sora.cluster-metrics` 有効時に影響する。

## 完了条件

- 4 件の Erlang VM メトリクスが `newCounter` に変更されている
- `raftTerm` が `newGauge` に変更されている
- 影響を受けるテストフィクスチャの `# TYPE` 行が更新されている
- `go test -race -v .` が通る

## 解決方法

1. `collector/erlang_vm.go:124,125,132,137` の `newGauge` を `newCounter` に変更する
2. `collector/cluster_node.go:119` の `newCounter` を `newGauge` に変更する
3. 以下のテストフィクスチャで `# TYPE` 行を更新する:
   - `test/maximum.metrics`（erlang_vm 4 件: gauge→counter、raft_term: counter→gauge）
   - `test/sora_erlang_vm_enabled.metrics`（erlang_vm 4 件: gauge→counter）
   - `test/sora_cluster_metrics_enabled.metrics`（raft_term: counter→gauge）
   - `test/sora_up_cluster_nodes_failed.metrics`（erlang_vm 4 件: gauge→counter）
   - `test/sora_up_license_failed.metrics`（erlang_vm 4 件: gauge→counter、raft_term: counter→gauge）
4. `go test -race -v .` で全テストの通過を確認する
