# Sora exporter の利用方法

## インストール

Sora exporter は [Releases](https://github.com/shiguredo/sora_exporter/releases) より環境に合わせてビルド済みバイナリをダウンロードして利用できます。

```sh
curl -L https://github.com/shiguredo/sora_exporter/releases/download/{VERSION}/sora_exporter_linux_{CPU_ARCH}-{VERSION}.gz -o sora_exporter.gz
gzip -d sora_exporter.gz
chmod +x sora_exporter
```

## 起動方法

Sora exporter はコマンドライン引数で Sora の管理 API エンドポイントを指定して起動できます。

ここでは Sora が同一ホスト上に起動しており、Sora の API が 127.0.0.1 の 3000 番ポートで待ち受けている場合の例を示します。

```sh
./sora_exporter
```

実行すると次のようなログが出力され、Sora exporter が起動します。

```
time=2025-12-25T09:31:17.942Z level=INFO source=main.go:171 msg="Starting sora_exporter" version="(version=2025.2.0, branch=main, revision=46ecaf8)"
time=2025-12-25T09:31:17.944Z level=INFO source=main.go:172 msg="Build context" build_context="(go=go1.25.5, platform=linux/amd64, user=shiguredo, date=2025-12-25T06:54:50Z, tags=unknown)"
time=2025-12-25T09:31:17.946Z level=INFO source=main.go:194 msg="Listening on" address=:9490
time=2025-12-25T09:31:17.951Z level=INFO source=tls_config.go:354 msg="Listening on" address=[::]:9490
time=2025-12-25T09:31:17.952Z level=INFO source=tls_config.go:357 msg="TLS is disabled." http2=false address=[::]:9490
```

Sora exporter はデフォルトでポート 9490 番で待ち受けます。Prometheus サーバーはこのポートに対してメトリクスの取得リクエストを送信します。

Sora exporter を起動したホスト上で curl コマンドを実行して、メトリクスが取得できることを確認します。

```
$ curl -s http://127.0.0.1:9490/metrics
# HELP promhttp_metric_handler_errors_total Total number of internal errors encountered by the promhttp metric handler.
# TYPE promhttp_metric_handler_errors_total counter
promhttp_metric_handler_errors_total{cause="encoding"} 0
promhttp_metric_handler_errors_total{cause="gathering"} 0
# HELP sora_auth_webhook_total The total number of auth webhook.
# TYPE sora_auth_webhook_total counter
sora_auth_webhook_total{state="failed"} 1
sora_auth_webhook_total{state="successful"} 9621
# HELP sora_average_duration_seconds The average connection duration in seconds.
# TYPE sora_average_duration_seconds gauge
sora_average_duration_seconds 44
...
```

以上のようにメトリクスが取得できれば、Sora exporter の起動は成功しています。

Sora exporter が接続する Sora API のデフォルト URL は `http://127.0.0.1:3000/` です。

Sora が同一ホスト上で、API ポートをデフォルト設定の 3000 番ポートで待ち受けている場合は、引数指定なしでも Sora API に接続できます。

もし、Sora API が別のホストやポートで待ち受けている場合は、`--sora.api-url` 引数でエンドポイント URL を指定します。

```sh
./sora_exporter --sora.api-url="http://192.0.2.1:3000/"
```

これで Sora exporter が Sora API を呼び出す先は `http://192.0.2.1:3000/` になります。

Sora exporter 自身はデフォルトでポート番号 9490 番で待ち受けますが、`--web.listen-address` 引数で変更できます。

この起動オプションでは、ポート番号だけではなく、IP アドレスも指定できます。

```sh
./sora_exporter --web.listen-address="127.0.0.1:9000"
```

上記の指定例では、Sora exporter は Listen するアドレスは 127.0.0.1 のみで 9000 番ポートで待ち受けるようになります。

プライベートネットワークからのみアクセスしたい場合などにはこのオプションを利用してください。

その他の利用可能なオプションは `-h` 引数で確認できます

```sh
./sora_exporter -h
```

## Sora のメトリクス収集

Sora exporter は以下のような Sora のメトリクスを収集できます。

- 接続数
- ウェブフックの成功、失敗数
- Sora の転送バイト数、パケット数
- Sora のライセンス情報

これに加えて Sora をクラスター構成で利用している場合に、起動オプションを追加することでクラスターのメトリクスも収集できます。

クラスターのメトリクスを収集するには `--sora.cluster-metrics` を起動時に追加します。

```sh
./sora_exporter --sora.cluster-metrics
```

このオプションを追加することで、Sora クラスターが認識している他のノードの状態や、各ノード間の転送バイト数などのメトリクスも収集できるようになります。

## Sora の起動状態の確認について

### sora_up メトリクス

Sora exporter では Sora の起動状態を表すメトリクス `sora_up` を提供しています。

Sora exporter はリクエストを受け取ると、次の Sora API を呼び出します。

- `Sora_20171010.GetStatsReport`
- `Sora_20171218.GetLicense`

すべての API 呼び出しが成功した場合は `sora_up` メトリクスに `1` を返します。

いずれかの API 呼び出しが失敗した場合は `sora_up` を `0` を返します。この値が返された場合は、Sora exporter から Sora への接続に問題が発生している可能性があります。

また、`sora_up` メトリクスが `0` を返したとしても、Sora API から必要な情報が取得できた部分だけ、レスポンスのメトリクスに反映されます。

例えば `Sora_20171218.GetLicense` API の呼び出しが失敗した場合は、ライセンス情報のメトリクスのみ欠落しますが、それ以外のメトリクスはレスポンスに反映されます。

### sora_cluster_up メトリクス

Sora クラスターのメトリクス収集を有効にしている場合は `sora_cluster_up` メトリクスも提供されます。

この `sora_cluster_up` メトリクスは次の Sora API を呼び出し、成功した場合に `1` を設定して返します。

- `Sora_20211215.ListClusterNodes`

`sora_cluster_up` メトリクスは `0` を返したとしても、Sora クラスター全体が必ずしもアクセス不可となっているわけではありませんが Sora クラスターの状態を確認する必要があると理解してください。

また、 `Sora_20211215.ListClusterNodes` API の呼び出しに失敗し、`sora_cluster_up` メトリクスが `0` を返した場合でも、 `Sora_20171010.GetStatsReport` API に含まれるクラスター関連のメトリクスはレスポンスに反映されます。
