# Dockerfile のベースイメージを debian12 に更新する

- Priority: Medium
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/change-dockerfile-base-image
- Polished: {YYYY-MM-DD}

## 目的

EOL の Debian 10 ベースイメージを Debian 12 に更新し、セキュリティ修正が提供される状態にする。

## 優先度根拠

Debian 10 (buster) は 2024-06-30 に LTS サポートが終了しており、セキュリティ修正が提供されない。

## 現状

`Dockerfile:1`:

```dockerfile
FROM gcr.io/distroless/static-debian10
```

## 完了条件

- ベースイメージが `gcr.io/distroless/static-debian12` に更新されている
- Docker ビルドが成功する

## 解決方法

`Dockerfile:1` を `FROM gcr.io/distroless/static-debian12` に変更する。
