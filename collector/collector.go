package collector

import (
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

type Collector struct {
	mutex                            sync.RWMutex
	logger                           log.Logger
	timeout                          time.Duration
	URI                              string
	skipSslVerify                    bool
	enableSoraClientMetrics          bool
	enableSoraConnectionErrorMetrics bool
	enableErlangVmMetrics            bool
	EnableSoraClusterMetrics         bool

	soraUp          *prometheus.Desc
	soraVersionInfo *prometheus.Desc
	ConnectionMetrics
	ClientMetrics
	SoraConnectionErrorMetrics
	ErlangVmMetrics
	SoraClusterMetrics
}

type CollectorOptions struct {
	URI                              string
	SkipSslVerify                    bool
	Timeout                          time.Duration
	Logger                           log.Logger
	EnableSoraClientMetrics          bool
	EnableSoraConnectionErrorMetrics bool
	EnableErlangVmMetrics            bool
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

		enableSoraClientMetrics:          options.EnableSoraClientMetrics,
		enableSoraConnectionErrorMetrics: options.EnableSoraConnectionErrorMetrics,
		enableErlangVmMetrics:            options.EnableErlangVmMetrics,
		EnableSoraClusterMetrics:         options.EnableSoraClusterMetrics,

		soraUp:                     newDesc("up", "Whether the last scrape of metrics from Sora was able to connect to the server (1 for yes, 0 for no)."),
		soraVersionInfo:            newDescWithLabel("version_info", "sora version info.", []string{"version"}),
		ConnectionMetrics:          connectionMetrics,
		ClientMetrics:              clientMetrics,
		SoraConnectionErrorMetrics: soraConnectionErrorMetrics,
		ErlangVmMetrics:            erlangVmMetrics,
		SoraClusterMetrics:         soraClusterMetrics,
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
		req, err = http.NewRequestWithContext(ctx, http.MethodPost, c.URI, nil)
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

	ch <- newGauge(c.soraUp, 1)
	ch <- newGauge(c.soraVersionInfo, 1, report.SoraVersion)
	c.ConnectionMetrics.Collect(ch, report.soraConnectionReport)

	if c.enableSoraClientMetrics {
		c.ClientMetrics.Collect(ch, report.SoraClientReport)
	}
	if c.enableSoraConnectionErrorMetrics {
		c.SoraConnectionErrorMetrics.Collect(ch, report.SoraConnectionErrorReport)
	}
	if c.enableErlangVmMetrics {
		c.ErlangVmMetrics.Collect(ch, report.ErlangVmReport)
	}
	if c.EnableSoraClusterMetrics {
		c.SoraClusterMetrics.Collect(ch, nodeList)
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.soraUp
	ch <- c.soraVersionInfo
	c.ConnectionMetrics.Describe(ch)

	if c.enableSoraClientMetrics {
		c.ClientMetrics.Describe(ch)
	}
	if c.enableSoraConnectionErrorMetrics {
		c.SoraConnectionErrorMetrics.Describe(ch)
	}
	if c.enableErlangVmMetrics {
		c.ErlangVmMetrics.Describe(ch)
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
