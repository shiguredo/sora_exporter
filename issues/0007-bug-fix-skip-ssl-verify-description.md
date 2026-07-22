# sora.skip-ssl-verify フラグの説明文が意味と正反対

- Priority: Medium
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/fix-skip-ssl-verify-description
- Polished: {YYYY-MM-DD}

## 目的

`--sora.skip-ssl-verify` フラグの説明文がフラグ名と矛盾しているのを修正する。

## 優先度根拠

セキュリティに関わるオプションの説明文が正反対の意味を持っており、ユーザーの誤解を招く。

## 現状

`main.go:67-71`:

```go
soraSkipSslVerify = kingpin.Flag(
    "sora.skip-ssl-verify",
    "Flag that enables SSL certificate verification for the Sora URL",
).Bool()
```

フラグ名は `skip-ssl-verify`（SSL 検証をスキップ）だが、説明文は "Flag that **enables** SSL certificate verification" と書かれている。

## 完了条件

- 説明文がフラグ名と整合する意味になっている
- 全テストが通る

## 解決方法

説明文を `"Flag that disables SSL certificate verification for the Sora URL"` に修正する。
