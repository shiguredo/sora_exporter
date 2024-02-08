package main

import (
	stdlog "log"
	"net/http"
	"os"
	"os/user"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/shiguredo/sora_exporter/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"

	promcollectors "github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/prometheus/exporter-toolkit/web"

	kingpin "github.com/alecthomas/kingpin/v2"
)

var (
	listenAddress = kingpin.Flag(
		"web.listen-address",
		"Address on which to expose metrics and web interface.",
	).Default(":9490").String()
	metricsPath = kingpin.Flag(
		"web.telemetry-path",
		"Path under which to expose metrics.",
	).Default("/metrics").String()
	soraAPIURL = kingpin.Flag(
		"sora.api-url",
		"URL on which to scrape Sora API",
	).Default("http://127.0.0.1:3000/").String()
	soraTimeout = kingpin.Flag(
		"sora.timeout",
		"Timeout for trying to get stats from Sora API URL",
	).Default("5s").Duration()
	disableExporterMetrics = kingpin.Flag(
		"web.disable-exporter-metrics",
		"Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).",
	).Bool()
	// この統計情報はアンドキュメントです
	enableSoraClientMetrics = kingpin.Flag(
		"sora.client-metrics",
		"Include metrics about Sora client connection stats.",
	).Bool()
	// この統計情報はアンドキュメントです
	enableSoraConnectionErrorMetrics = kingpin.Flag(
		"sora.connection-error-metrics",
		"Include metrics about Sora connection error stats.",
	).Bool()
	// この統計情報はアンドキュメントです
	enableErlangVMMetrics = kingpin.Flag(
		"sora.erlang-vm-metrics",
		"Include metrics about Erlang VM stats.",
	).Bool()
	enableSoraClusterMetrics = kingpin.Flag(
		"sora.cluster-metrics",
		"Include metrics about Sora cluster stats.",
	).Bool()
	soraSkipSslVeirfy = kingpin.Flag(
		"sora.skip-ssl-verify",
		"Flag that enables SSL certificate verification for the Sora URL",
	).Bool()
	maxRequests = kingpin.Flag(
		"web.max-requests",
		"Maximum number of parallel scrape requests. Use 0 to disable.",
	).Default("40").Int()
)

type handler struct {
	soraMetricsHandler http.Handler
	// exporterMetricsRegistry is a separate registry for the metrics about
	// the exporter itself.
	exporterMetricsRegistry          *prometheus.Registry
	includeExporterMetrics           bool
	maxRequests                      int
	logger                           log.Logger
	soraAPIURL                       string
	soraSkipSslVeirfy                bool
	soraTimeout                      time.Duration
	soraFreezeTimeSeconds            bool
	enableSoraClientMetrics          bool
	enableSoraConnectionErrorMetrics bool
	enableErlangVMMetrics            bool
	enableSoraClusterMetrics         bool
}

func newHandler(
	includeExporterMetrics bool, maxRequests int, logger log.Logger,
	soraAPIURL string, soraSkipSslVeirfy bool, soraTimeout time.Duration, soraFreezeTimeSeconds bool,
	enableSoraClientMetrics bool, enableSoraConnectionErrorMetrics bool, enableErlangVMMetrics bool, enableSoraClusterMetrics bool) *handler {

	h := &handler{
		exporterMetricsRegistry:          prometheus.NewRegistry(),
		includeExporterMetrics:           includeExporterMetrics,
		maxRequests:                      maxRequests,
		logger:                           logger,
		soraAPIURL:                       soraAPIURL,
		soraSkipSslVeirfy:                soraSkipSslVeirfy,
		soraTimeout:                      soraTimeout,
		soraFreezeTimeSeconds:            soraFreezeTimeSeconds,
		enableSoraClientMetrics:          enableSoraClientMetrics,
		enableSoraConnectionErrorMetrics: enableSoraConnectionErrorMetrics,
		enableErlangVMMetrics:            enableErlangVMMetrics,
		enableSoraClusterMetrics:         enableSoraClusterMetrics,
	}
	if h.includeExporterMetrics {
		h.exporterMetricsRegistry.MustRegister(
			promcollectors.NewProcessCollector(promcollectors.ProcessCollectorOpts{}),
			promcollectors.NewGoCollector(),
		)
	}
	h.soraMetricsHandler = h.innerHandler()
	return h
}

// ServeHTTP implements http.Handler.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.soraMetricsHandler.ServeHTTP(w, r)
}

func (h *handler) innerHandler() http.Handler {
	r := prometheus.NewRegistry()
	r.MustRegister(version.NewCollector("sora_exporter"))
	r.MustRegister(collector.NewCollector(&collector.CollectorOptions{
		URI:                              h.soraAPIURL,
		SkipSslVerify:                    h.soraSkipSslVeirfy,
		Timeout:                          h.soraTimeout,
		FreezeTimeSeconds:                h.soraFreezeTimeSeconds,
		Logger:                           h.logger,
		EnableSoraClientMetrics:          h.enableSoraClientMetrics,
		EnableSoraConnectionErrorMetrics: h.enableSoraConnectionErrorMetrics,
		EnableErlangVMMetrics:            h.enableErlangVMMetrics,
		EnableSoraClusterMetrics:         h.enableSoraClusterMetrics,
	}))
	handler := promhttp.HandlerFor(
		prometheus.Gatherers{h.exporterMetricsRegistry, r},
		promhttp.HandlerOpts{
			ErrorLog:            stdlog.New(log.NewStdlibAdapter(level.Error(h.logger)), "", 0),
			ErrorHandling:       promhttp.ContinueOnError,
			MaxRequestsInFlight: h.maxRequests,
			Registry:            h.exporterMetricsRegistry,
		},
	)
	if h.includeExporterMetrics {
		// Note that we have to use h.exporterMetricsRegistry here to
		// use the same promhttp metrics for all expositions.
		handler = promhttp.InstrumentMetricHandler(
			h.exporterMetricsRegistry, handler,
		)
	}
	return handler
}

func main() {
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("sora_exporter"))
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting sora_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext())

	// root 権限で起動してたら warning を出す
	if user, err := user.Current(); err == nil && user.Uid == "0" {
		level.Warn(logger).Log("msg", "Sora Exporter is running as root user. This exporter is designed to run as unpriviledged user, root is not required.")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Sora Exporter</title></head>
			<body>
			<h1>Sora Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})
	soraHandler := newHandler(
		!*disableExporterMetrics, *maxRequests, logger,
		*soraAPIURL, *soraSkipSslVeirfy, *soraTimeout, false,
		*enableSoraClientMetrics, *enableSoraConnectionErrorMetrics, *enableErlangVMMetrics, *enableSoraClusterMetrics)
	http.Handle(*metricsPath, soraHandler)

	level.Info(logger).Log("msg", "Listening on", "address", *listenAddress)
	server := &http.Server{}
	webSystemdSocket := false
	webConfigFile := ""
	webFlagConfig := &web.FlagConfig{
		WebListenAddresses: &[]string{*listenAddress},
		WebSystemdSocket:   &webSystemdSocket,
		WebConfigFile:      &webConfigFile,
	}
	if err := web.ListenAndServe(server, webFlagConfig, logger); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
}
