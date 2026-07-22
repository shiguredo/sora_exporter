# CI ワークフローの不整合を解消する

- Priority: Low
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/change-unify-ci-workflows
- Polished: {YYYY-MM-DD}

## 目的

CI ワークフロー間の不整合（runs-on、テストコマンド）を解消する。

## 優先度根拠

CI の一貫性の問題。テストの race detector 有無がローカルと CI で異なる。

## 現状

### runs-on の不統一

- `.github/workflows/ci.yml:15`: `ubuntu-slim`
- `.github/workflows/release.yml:11`: `ubuntu-slim`
- `.github/workflows/test.yml:18`: `ubuntu-24.04`

### テストコマンドの不整合

- `Makefile:7`: `go test -race -v .`（race detector 有効）
- `.github/workflows/test.yml:29`: `go test -v`（race detector なし）

## 完了条件

- 全ワークフローの `runs-on` が統一されている
- CI のテストコマンドに `-race` フラグが追加されている

## 解決方法

`test.yml` の `runs-on` を `ubuntu-slim` に統一し、テストコマンドを `go test -race -v .` に変更する。
