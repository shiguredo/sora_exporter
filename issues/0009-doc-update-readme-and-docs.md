# README・docs・CHANGES.md の記載を修正する

- Priority: Low
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/fix-docs
- Polished: {YYYY-MM-DD}

## 目的

README.md、docs/README.md、CHANGES.md の記載不備・typo・古い情報を修正する。

## 優先度根拠

ドキュメントの正確性の問題であり、コードの動作には影響しない。

## 現状

### README.md の著作権表記が古い

`README.md:39-40`: `Copyright 2022-2023` のまま。2024 年〜2026 年もリリースが継続されている。

### docs/README.md に --sora.skip-ssl-verify が未記載

セキュリティ関連オプションだがドキュメントに記載がない。`// この統計情報はアンドキュメントです` コメントもなく、意図的な非公開対象ではない。

### CHANGES.md の typo

- 2024.4.0（:280,283）: `sora_cluster_relay_recived_bytes_total` / `recived_packets_total`（正しくは `received`）。実際のコードは正しい綴り
- 2024.7.0（:223）: `Sora expoter`（正しくは `Sora exporter`）
- 2025.2.0（:133）: `srtp_decrpyted_bytes_total`（`sora_` プレフィックス省略）

### CHANGES.md の表記不統一

- 2024.4.0 以前にリリース日記載なし（2024.5.0 以降は記載あり）
- 2022.6.0 は「クラスタ」、2024.3.0 以降は「クラスター」と長音符号が不統一

## 完了条件

- README.md の著作権表記が `2022-2026` に更新されている
- docs/README.md に `--sora.skip-ssl-verify` の記載がある
- CHANGES.md の typo が修正されている

## 解決方法

各ファイルを編集して修正する。CHANGES.md の過去セクションのリリース日・表記統一は任意。
