# CHANGES

## develop

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
