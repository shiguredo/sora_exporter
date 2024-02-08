package collector

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// for testing
	freezedTimeSeconds = float64(time.Date(2024, 1, 7, 17, 41, 31, 312389, time.UTC).UnixNano()) / 1e9
)

type Collector struct {
	mutex                            sync.RWMutex
	logger                           log.Logger
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
	Logger                           log.Logger
	EnableSoraClientMetrics          bool
	EnableSoraConnectionErrorMetrics bool
	EnableErlangVMMetrics            bool
	EnableSoraClusterMetrics         bool
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type SoraListClusterNodesRequest struct {
	IncludeAllKnownNodes bool `json:"include_all_known_nodes"`
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URI, nil)
	if err != nil {
		level.Error(c.logger).Log("msg", "failed to create request to sora", "err", err)
		ch <- newGauge(c.soraUp, 0)
		return
	}
	req.Header.Set("x-sora-target", "Sora_20171010.GetStatsReport")

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: c.skipSslVerify}}
	client := http.Client{
		Timeout:   c.timeout,
		Transport: tr,
	}

	resp, err := client.Do(req)
	if err != nil {
		level.Error(c.logger).Log("msg", "failed to request to Sora GetStatsReport API", "err", err)
		ch <- newGauge(c.soraUp, 0)
		return
	}
	defer resp.Body.Close()

	var report soraGetStatsReport
	if err := json.NewDecoder(resp.Body).Decode(&report); err != nil {
		level.Error(c.logger).Log("msg", "failed to decode response body from Sora GetStatsReport API", "err", err)
		ch <- newGauge(c.soraUp, 0)
		return
	}

	var nodeList []soraClusterNode
	if c.EnableSoraClusterMetrics {
		requestParams := SoraListClusterNodesRequest{
			IncludeAllKnownNodes: true,
		}
		encodedParams, err := json.Marshal(requestParams)
		if err != nil {
			level.Error(c.logger).Log("msg", "failed to encode Sora ListClusterNodes API request parameters", "err", err)
			ch <- newGauge(c.soraUp, 0)
			return
		}

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, c.URI, bytes.NewBuffer(encodedParams))
		if err != nil {
			level.Error(c.logger).Log("msg", "failed to create request to sora", "err", err)
			ch <- newGauge(c.soraUp, 0)
			return
		}
		req.Header.Set("x-sora-target", "Sora_20211215.ListClusterNodes")

		nodeResp, err := client.Do(req)
		if err != nil {
			level.Error(c.logger).Log("msg", "failed to request to Sora ListClusterNodes API", "err", err)
			ch <- newGauge(c.soraUp, 0)
			return
		}
		defer nodeResp.Body.Close()

		if err := json.NewDecoder(nodeResp.Body).Decode(&nodeList); err != nil {
			level.Error(c.logger).Log("msg", "failed to decode response body from Sora ListClusterNodes API", "err", err)
			ch <- newGauge(c.soraUp, 0)
			return
		}
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodPost, c.URI, nil)
	if err != nil {
		level.Error(c.logger).Log("msg", "failed to create request to sora", "err", err)
		ch <- newGauge(c.soraUp, 0)
		return
	}
	req.Header.Set("x-sora-target", "Sora_20171218.GetLicense")

	licenseResp, err := client.Do(req)
	if err != nil {
		level.Error(c.logger).Log("msg", "failed to request to Sora GetLicense API", "err", err)
		ch <- newGauge(c.soraUp, 0)
		return
	}
	defer licenseResp.Body.Close()

	var licenseInfo soraLicenseInfo
	if err := json.NewDecoder(licenseResp.Body).Decode(&licenseInfo); err != nil {
		level.Error(c.logger).Log("msg", "failed to decode response body from Sora GetLicense API", "err", err)
		ch <- newGauge(c.soraUp, 0)
		return
	}

	ch <- newGauge(c.soraUp, 1)
	ch <- newGauge(c.soraVersionInfo, 1, report.SoraVersion)

	if c.freezeTimeSeconds {
		ch <- newGauge(c.soraTimeSeconds, freezedTimeSeconds)
	} else {
		nowSec := float64(time.Now().UnixNano()) / 1e9
		ch <- newGauge(c.soraTimeSeconds, nowSec)
	}

	c.LicenseMetrics.Collect(ch, licenseInfo)
	c.ConnectionMetrics.Collect(ch, report.soraConnectionReport)
	c.WebhookMetrics.Collect(ch, report.soraWebhookReport)

	if c.enableSoraClientMetrics {
		c.ClientMetrics.Collect(ch, report.SoraClientReport)
	}
	if c.enableSoraConnectionErrorMetrics {
		c.SoraConnectionErrorMetrics.Collect(ch, report.SoraConnectionErrorReport)
	}
	if c.enableErlangVMMetrics {
		c.ErlangVMMetrics.Collect(ch, report.ErlangVMReport)
	}
	if c.EnableSoraClusterMetrics {
		c.SoraClusterMetrics.Collect(ch, nodeList, report.ClusterReport)
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.soraUp
	ch <- c.soraVersionInfo
	ch <- c.soraTimeSeconds
	c.LicenseMetrics.Describe(ch)
	c.ConnectionMetrics.Describe(ch)
	c.WebhookMetrics.Describe(ch)

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
