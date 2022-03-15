package main

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
	mutex                   sync.RWMutex
	logger                  log.Logger
	timeout                 time.Duration
	URI                     string
	skipSslVerify           bool
	enableSoraClientMetrics bool
	enableSoraErrorMetrics  bool
	enableErlangVmMetrics   bool

	soraVersionInfo *prometheus.Desc
	ConnectionMetrics
	ClientMetrics
	ErrorMetrics
	ErlangVmMetrics
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

func NewCollector(
	uri string, skipSslVerify bool, timeout time.Duration, logger log.Logger,
	enableSoraClientMetrics bool, enableSoraErrorMetrics bool, enableErlangVmMetrics bool) *Collector {
	return &Collector{
		URI:           uri,
		timeout:       timeout,
		skipSslVerify: skipSslVerify,
		logger:        logger,

		enableSoraClientMetrics: enableSoraClientMetrics,
		enableSoraErrorMetrics:  enableSoraErrorMetrics,
		enableErlangVmMetrics:   enableErlangVmMetrics,

		soraVersionInfo:   newDescWithLabel("sora_version_info", "sora version info.", []string{"version"}),
		ConnectionMetrics: connectionMetrics,
		ClientMetrics:     clientMetrics,
		ErrorMetrics:      errorMetrics,
		ErlangVmMetrics:   erlangVmMetrics,
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
		level.Error(c.logger).Log("msg", "failed to request to sora", "err", err)
		return
	}
	defer resp.Body.Close()

	var report soraGetStatsReport
	if err := json.NewDecoder(resp.Body).Decode(&report); err != nil {
		level.Error(c.logger).Log("msg", "failed to decode response body from sora", "err", err)
		return
	}

	ch <- newInfo(c.soraVersionInfo, report.SoraVersion)
	c.ConnectionMetrics.Collect(ch, report.soraConnectionReport)

	if c.enableSoraClientMetrics {
		c.ClientMetrics.Collect(ch, report.soraClientReport)
	}
	if c.enableSoraErrorMetrics {
		c.ErrorMetrics.Collect(ch, report.soraErrorReport)
	}
	if c.enableErlangVmMetrics {
		c.ErlangVmMetrics.Collect(ch, report.erlangVmReport)
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.soraVersionInfo
	c.ConnectionMetrics.Describe(ch)

	if c.enableSoraClientMetrics {
		c.ClientMetrics.Describe(ch)
	}
	if c.enableSoraErrorMetrics {
		c.ErrorMetrics.Describe(ch)
	}
	if c.enableErlangVmMetrics {
		c.ErlangVmMetrics.Describe(ch)
	}
}
