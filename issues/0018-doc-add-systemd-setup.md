# systemd での起動方法をドキュメントに追加する

- Priority: Low
- Created: 2026-09-10
- Completed: {YYYY-MM-DD}
- Model: OpenCode
- Branch: feature/add-systemd-setup
- Polished: {YYYY-MM-DD}

## 目的

Sora exporter を systemd のサービスとして起動するためのユニットファイルの例と、起動・停止・有効化の手順を提供する。

Sora 本体の運用では systemd の適用が強く推奨されている。Sora exporter は Sora と同一ホストで常駐させることが多く、手動起動の説明だけでは実運用に不足している。

## 優先度根拠

ドキュメントと設定例の追加であり、Sora exporter の動作には影響しない。常駐運用に必要となる情報のため Low とする。

## 現状

- docs/README.md はバイナリのダウンロード、手動起動、起動オプション、メトリクスの収集方法を説明しているが、systemd でサービスとして起動する方法を説明していない
- リポジトリに systemd のユニットファイルが含まれていない
- Sora 本体は sora-doc の SYSTEMD.rst で `sora.service` のユニットファイル例と `systemctl` の操作を説明しており、Sora exporter には同等の情報がない

## 完了条件

- systemd のユニットファイルの例がリポジトリに追加されている
- docs/README.md に systemd での起動・停止・有効化の手順が追加されている

## 解決方法

- Suzu の `script/suzu.service` と同様に `script/sora_exporter.service` を追加する
- ユニットファイルには `ExecStart`、`Restart`、`User`、`Group` などを設定し、Sora 本体のユニットファイルと同様に利用者が環境に合わせてパスを置き換える前提の例とする
- docs/README.md にユニットファイルの配置手順と `systemctl enable` / `start` / `stop` の操作を追加する
