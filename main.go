package main

import (
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"os/user"

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

	kingpin "gopkg.in/alecthomas/kingpin.v2"
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
	soraGetStatsReportURL = kingpin.Flag(
		"sora.get-stats-report-url",
		"URL on which to scrape Sora GetStatsReport API",
	).Default("http://127.0.0.1:3000/").String()
	soraTimeout = kingpin.Flag(
		"sora.timeout",
		"Timeout for trying to get stats from Sora GetStatsReport API URL",
	).Default("5s").Duration()
	disableExporterMetrics = kingpin.Flag(
		"web.disable-exporter-metrics",
		"Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).",
	).Bool()
	enableSoraClientMetrics = kingpin.Flag(
		"sora.client-metrics",
		"Include metrics about Sora client connection stats.",
	).Bool()
	enableSoraErrorMetrics = kingpin.Flag(
		"sora.error-metrics",
		"Include metrics about Sora connect error stats.",
	).Bool()
	enableErlangVmMetrics = kingpin.Flag(
		"sora.erlang-vm-metrics",
		"Include metrics about Erlang VM stats.",
	).Bool()
	maxRequests = kingpin.Flag(
		"web.max-requests",
		"Maximum number of parallel scrape requests. Use 0 to disable.",
	).Default("40").Int()
)

type handler struct {
	unfilteredHandler http.Handler
	// exporterMetricsRegistry is a separate registry for the metrics about
	// the exporter itself.
	exporterMetricsRegistry *prometheus.Registry
	includeExporterMetrics  bool
	maxRequests             int
	logger                  log.Logger
}

func newHandler(includeExporterMetrics bool, maxRequests int, logger log.Logger) *handler {
	h := &handler{
		exporterMetricsRegistry: prometheus.NewRegistry(),
		includeExporterMetrics:  includeExporterMetrics,
		maxRequests:             maxRequests,
		logger:                  logger,
	}
	if h.includeExporterMetrics {
		h.exporterMetricsRegistry.MustRegister(
			promcollectors.NewProcessCollector(promcollectors.ProcessCollectorOpts{}),
			promcollectors.NewGoCollector(),
		)
	}
	if innerHandler, err := h.innerHandler(); err != nil {
		panic(fmt.Sprintf("Couldn't create metrics handler: %s", err))
	} else {
		h.unfilteredHandler = innerHandler
	}
	return h
}

// ServeHTTP implements http.Handler.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filters := r.URL.Query()["collect[]"]
	level.Debug(h.logger).Log("msg", "collect query:", "filters", filters)

	if len(filters) == 0 {
		// No filters, use the prepared unfiltered handler.
		h.unfilteredHandler.ServeHTTP(w, r)
		return
	}
	// To serve filtered metrics, we create a filtering handler on the fly.
	filteredHandler, err := h.innerHandler(filters...)
	if err != nil {
		level.Warn(h.logger).Log("msg", "Couldn't create filtered metrics handler:", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Couldn't create filtered metrics handler: %s", err)))
		return
	}
	filteredHandler.ServeHTTP(w, r)
}

func (h *handler) innerHandler(filters ...string) (http.Handler, error) {
	r := prometheus.NewRegistry()
	r.MustRegister(version.NewCollector("sora_exporter"))
	r.MustRegister(collector.New(
		*soraGetStatsReportURL,
		*soraTimeout,
		h.logger,
		*enableSoraClientMetrics,
		*enableSoraErrorMetrics,
		*enableErlangVmMetrics))
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
	return handler, nil
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
	http.Handle(*metricsPath, newHandler(!*disableExporterMetrics, *maxRequests, logger))

	level.Info(logger).Log("msg", "Listening on", "address", *listenAddress)
	server := &http.Server{Addr: *listenAddress}
	if err := web.ListenAndServe(server, "", logger); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
}
