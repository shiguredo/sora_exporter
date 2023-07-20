# CHANGES

## develop

- [ADD] `sora_client` の `obs_studio_whip` に対応する
  - @voluntas
- [CHANGE] 依存パッケージを更新する
  - `prometheus/client-golang` 1.14.0 => 1.16.0
  - `prometheus/common` 0.41.0 => 0.44.0
  - `prometheus/exporter-toolkit` 0.9.1 => 0.10.0

## 2023.2.0

- [CHANGE] kingpin の更新に伴うパッケージ名の変更
  - kingpin のバージョンは `v2.2.6` から `v2.3.2` に更新
  - `gopkg.in/alecthomas/kingpin.v2` から `github.com/alecthomas/kingpin/v2` に変更
  - kingpin に依存している関連パッケージの更新
  - @tnamao
- [CHANGE] Sora exporter がスクレイピングする Sora API のオプション名を変更する
  - コマンドライン引数の `--sora.get-stats-report-url` を `--sora.api-url` に変更する
  - 破壊的変更になるため、バージョンアップの際に注意してください
  - @tnamao
- [ADD] Sora のライセンス情報を返すメトリクスを追加する
  - `sora_license_info` ライセンスのテキスト情報
  - `sora_license_max_connections` ライセンスの同時接続数
  - `sora_license_max_nodes` クラスターライセンスに含まれる最大ノード数
    - GetLicense API のレスポンスに `max_nodes` が含まれる場合のみにメトリクスを返す
  - @tnamao

## 2023.1.0

- [CHANGE] Go 1.20 に上げる
  - @tnamao
- [CHANGE] staticcheck を 2023.1.1 に上げる
  - @tnamao
- [CHANGE] staticcheck-action を 1.3.0 に上げる
  - @tnamao

## 2022.6.1

- [CHANGE] リリース用 Github Actions のワークフローを修正
  - @tnamao

## 2022.6.0

- [ADD] クラスタ機能で使用している Raft 関連のメトリクスを追加
  - 以下の三つのメトリクスを追加する
    - [counter] sora_cluster_raft_commit_index ${INDEX}
    - [counter] sora_cluster_raft_term ${TERM}
    - [gauge] sora_cluster_raft_state { state = "${STATE_NAME}" } 1
  - これらは、以下の条件が満たされた時だけ、結果に含まれる
    - Sora のクラスタ機能が有効になっている
    - sora_exporter が `--sora.cluster-metrics` オプション付きで起動されている
  - @sile

## 2022.5.0

- [ADD] Sora の接続クライアントメトリクスに `flutter_sdk` を追加する
  - @tnamao
- [CHANGE] exporter-toolkit の変更に追従する
  - @tnamao

## 2022.4.0

- [ADD] Sora の Webhook メトリクスに対応する
  - @tnamao
- [CHANGE] Go 1.19 に上げる
  - @tnamao

## 2022.3.0

- [ADD] Sora の接続クライアントメトリクスに `cpp_sdk` と `zakuro` を追加する
  - @tnamao
- [ADD] [Staticcheck](https://staticcheck.io/) の静的解析に対応する
  - @tnamao

## 2022.2.0

- [CHANGE] Sora に関するメトリクスの接頭辞を `sora_exporter` から `sora` に変更する
  - @tnamao
- [CHANGE] 接続数のメトリクス名を `sora_connections_total` に変更し、状態はラベルに変更する
  - @tnamao
- [CHANGE] セッション数のメトリクス名を `sora_session_total` に変更し、状態はラベルに変更する
  - @tnamao
- [CHANGE] 接続エラーのメトリクス名を `sora_connection_error_total` に変更し、エラー理由はラベルに変更する
  - @tnamao
- [CHANGE] 接続クライアントのメトリクス名を `sora_client_type_total` に変更し、クライアント種別と接続結果はラベルに変更する
  - @tnamao
- [CHANGE] `total_received_invalid_turn_tcp_packet` を `received_invalid_turn_tcp_packet_total` に変更する
  - @tnamao
- [ADD] Sora への接続可否を判定するための `sora_up` メトリクスを追加する
  - @tnamao
- [ADD] Sora が認識しているクラスターノードをメトリクスで返す `sora_cluster_node` を追加する
  - @tnamao
- [CHANGE] GitHub リリースのタイトルにタグ名を設定する
  - @tnamao
- [CHANGE] リリース時のアーカイブファイルに対象の OS 名を含める
  - @tnamao

## 2022.1.0

- [CHANGE] `total_ongoing_connections` を counter => gauge に変更する
  - @tnamao
- [ADD] Sora の標準統計情報の `total_received_invalid_turn_tcp_packet` `total_session_created` `total_session_destroyed` に対応する
  - @tnamao
- [ADD] アンドキュメントの Sora 接続クライアントの統計情報に対応する
  - @tnamao
- [ADD] アンドキュメントの Sora 接続エラーの統計情報に対応する
  - @tnamao
- [ADD] アンドキュメントの Erlang VM の統計情報に対応する
  - @tnamao
- [ADD] スクレイプ先 URL の SSL 検証を行わないオプションを追加する
  - @tnamao
- [ADD] `/metrics` を同時呼び出しされた時に Sora の GetStatReport の呼び出しは一度に呼び出さないように mutex で同時実行を抑制する
  - @tnamao
- [ADD] テストコードを追加する
  - @tnamao
