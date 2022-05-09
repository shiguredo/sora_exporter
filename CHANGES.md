# CHANGES

## develop

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
