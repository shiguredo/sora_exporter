package collector

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// for testing
	freezedTimeSeconds = float64(time.Date(2024, 1, 7, 17, 41, 31, 312389, time.UTC).UnixNano()) / 1e9
)

type Collector struct {
	mutex                            sync.RWMutex
	logger                           *slog.Logger
	timeout                          time.Duration
	URI                              string
	skipSslVerify                    bool
	freezeTimeSeconds                bool
	enableSoraClientMetrics          bool
	enableSoraConnectionErrorMetrics bool
	enableErlangVMMetrics            bool
	EnableSoraClusterMetrics         bool

	soraUp          *prometheus.Desc
	soraVersionInfo *prometheus.Desc
	soraTimeSeconds *prometheus.Desc

	ConnectionMetrics
	WebhookMetrics
	SrtpMetrics
	SctpMetrics
	ClientMetrics
	SoraConnectionErrorMetrics
	ErlangVMMetrics
	SoraClusterMetrics
	LicenseMetrics
}

type CollectorOptions struct {
	URI                              string
	SkipSslVerify                    bool
	Timeout                          time.Duration
	FreezeTimeSeconds                bool
	Logger                           *slog.Logger
	EnableSoraClientMetrics          bool
	EnableSoraConnectionErrorMetrics bool
	EnableErlangVMMetrics            bool
	EnableSoraClusterMetrics         bool
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

func NewCollector(options *CollectorOptions) *Collector {
	return &Collector{
		URI:           options.URI,
		timeout:       options.Timeout,
		skipSslVerify: options.SkipSslVerify,
		logger:        options.Logger,

		// for testing
		freezeTimeSeconds: options.FreezeTimeSeconds,

		enableSoraClientMetrics:          options.EnableSoraClientMetrics,
		enableSoraConnectionErrorMetrics: options.EnableSoraConnectionErrorMetrics,
		enableErlangVMMetrics:            options.EnableErlangVMMetrics,
		EnableSoraClusterMetrics:         options.EnableSoraClusterMetrics,

		soraUp:          newDesc("up", "Whether the last scrape of metrics from Sora was able to connect to the server (1 for yes, 0 for no)."),
		soraVersionInfo: newDescWithLabel("version_info", "sora version info.", []string{"version"}),
		// same as node expoter's node_time_seconds
		soraTimeSeconds: newDesc("time_seconds", "System time in seconds since epoch."),

		ConnectionMetrics:          connectionMetrics,
		WebhookMetrics:             webhookMetrics,
		SrtpMetrics:                srtpMetrics,
		SctpMetrics:                sctpMetrics,
		ClientMetrics:              clientMetrics,
		SoraConnectionErrorMetrics: soraConnectionErrorMetrics,
		ErlangVMMetrics:            erlangVMMetrics,
		SoraClusterMetrics:         soraClusterMetrics,
		LicenseMetrics:             licenseMetrics,
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: c.skipSslVerify}}
	client := &http.Client{
		Timeout:   c.timeout,
		Transport: tr,
	}

	report, errFetchGetStatsReport := c.fetchGetStatsReport(ctx, client)
	licenseInfo, errFetchGetLicense := c.fetchGetLicense(ctx, client)
	// exporter の起動オプションで `--sora.cluster-metrics` を有効にしている場合はクラスター情報を収集する
	var nodeList *[]soraClusterNode
	var errFetchListClusterNodes error
	if c.EnableSoraClusterMetrics {
		nodeList, errFetchListClusterNodes = c.fetchListClusterNodes(ctx, client)
	}

	if errFetchGetStatsReport == nil && errFetchGetLicense == nil {
		// GetStatsReport と GetLicense の両方が成功した場合のみ Sora は up とみなす
		ch <- newGauge(c.soraUp, 1)
	} else {
		ch <- newGauge(c.soraUp, 0)
	}

	if c.freezeTimeSeconds {
		ch <- newGauge(c.soraTimeSeconds, freezedTimeSeconds)
	} else {
		nowSec := float64(time.Now().UnixNano()) / 1e9
		ch <- newGauge(c.soraTimeSeconds, nowSec)
	}

	if report != nil {
		ch <- newGauge(c.soraVersionInfo, 1, report.SoraVersion)
		c.ConnectionMetrics.Collect(ch, report.soraConnectionReport)
		c.WebhookMetrics.Collect(ch, report.soraWebhookReport)
		c.SrtpMetrics.Collect(ch, report.soraSrtpReport)
		c.SctpMetrics.Collect(ch, report.soraSctpReport)

		if c.enableSoraClientMetrics {
			c.ClientMetrics.Collect(ch, report.SoraClientReport)
		}
		if c.enableSoraConnectionErrorMetrics {
			c.SoraConnectionErrorMetrics.Collect(ch, report.SoraConnectionErrorReport)
		}
		if c.enableErlangVMMetrics {
			c.ErlangVMMetrics.Collect(ch, report.ErlangVMReport)
		}
	}

	if licenseInfo != nil {
		c.LicenseMetrics.Collect(ch, licenseInfo)
	}

	// exporter 起動時に EnableSoraClusterMetrics が有効であればクラスター情報を収集する
	if c.EnableSoraClusterMetrics {
		if errFetchListClusterNodes == nil {
			// クラスター API の呼び出しが成功した場合は Sora クラスターは up とみなす
			ch <- newGauge(c.soraClusterUp, 1)
			if nodeList != nil {
				c.SoraClusterMetrics.CollectClusterNodes(ch, nodeList)
			}
		} else {
			ch <- newGauge(c.soraClusterUp, 0)
		}
		// ListClusterNodes の呼び出しが失敗しても、GetStatsReport が成功していればレポートは存在する可能性があり
		// その場合は GetStatsReport の結果からクラスター情報を収集する
		if report != nil {
			c.SoraClusterMetrics.CollectClusterReport(ch, report.ClusterReport, report.ClusterRelay)
		}
	}
}

func (c *Collector) fetchGetStatsReport(ctx context.Context, client *http.Client) (*soraGetStatsReport, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URI, nil)
	if err != nil {
		c.logger.Error("failed to create request to sora", "err", err)
		return nil, err
	}
	req.Header.Set("x-sora-target", "Sora_20171010.GetStatsReport")

	resp, err := client.Do(req)
	if err != nil {
		c.logger.Error("failed to request to Sora GetStatsReport API", "err", err)
		return nil, err
	}
	defer resp.Body.Close()

	var report soraGetStatsReport
	if err := json.NewDecoder(resp.Body).Decode(&report); err != nil {
		c.logger.Error("failed to decode response body from Sora GetStatsReport API", "err", err)
		return nil, err
	}

	return &report, nil
}

func (c *Collector) fetchGetLicense(ctx context.Context, client *http.Client) (*soraLicenseInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URI, nil)
	if err != nil {
		c.logger.Error("failed to create request to sora", "err", err)
		return nil, err
	}
	req.Header.Set("x-sora-target", "Sora_20171218.GetLicense")

	licenseResp, err := client.Do(req)
	if err != nil {
		c.logger.Error("failed to request to Sora GetLicense API", "err", err)
		return nil, err
	}
	defer licenseResp.Body.Close()

	var licenseInfo soraLicenseInfo
	if err := json.NewDecoder(licenseResp.Body).Decode(&licenseInfo); err != nil {
		c.logger.Error("failed to decode response body from Sora GetLicense API", "err", err)
		return nil, err
	}

	return &licenseInfo, nil
}

func (c *Collector) fetchListClusterNodes(ctx context.Context, client *http.Client) (*[]soraClusterNode, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URI, nil)
	if err != nil {
		c.logger.Error("failed to create request to sora", "err", err.Error())
		return nil, err
	}
	req.Header.Set("x-sora-target", "Sora_20211215.ListClusterNodes")

	nodeResp, err := client.Do(req)
	if err != nil {
		c.logger.Error("failed to request to Sora ListClusterNodes API", "err", err)
		return nil, err
	}
	defer nodeResp.Body.Close()

	var nodeList []soraClusterNode
	if err := json.NewDecoder(nodeResp.Body).Decode(&nodeList); err != nil {
		c.logger.Error("failed to decode response body from Sora ListClusterNodes API", "err", err)
		return nil, err
	}

	return &nodeList, nil
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.soraUp
	ch <- c.soraVersionInfo
	ch <- c.soraTimeSeconds
	c.LicenseMetrics.Describe(ch)
	c.ConnectionMetrics.Describe(ch)
	c.WebhookMetrics.Describe(ch)
	c.SrtpMetrics.Describe(ch)
	c.SctpMetrics.Describe(ch)

	if c.enableSoraClientMetrics {
		c.ClientMetrics.Describe(ch)
	}
	if c.enableSoraConnectionErrorMetrics {
		c.SoraConnectionErrorMetrics.Describe(ch)
	}
	if c.enableErlangVMMetrics {
		c.ErlangVMMetrics.Describe(ch)
	}
	if c.EnableSoraClusterMetrics {
		ch <- c.soraClusterUp
		c.SoraClusterMetrics.Describe(ch)
	}
}

func newDesc(name, help string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName("sora", "", name), help, nil, nil)
}

func newDescWithLabel(name, help string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName("sora", "", name), help, labels, nil)
}

func newGauge(d *prometheus.Desc, v float64, labelValues ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(d, prometheus.GaugeValue, v, labelValues...)
}

func newCounter(d *prometheus.Desc, v float64, labelValues ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(d, prometheus.CounterValue, v, labelValues...)
}
