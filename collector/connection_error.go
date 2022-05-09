package collector

import "github.com/prometheus/client_golang/prometheus"

// この統計情報はアンドキュメントです
var (
	soraConnectionErrorMetrics = SoraConnectionErrorMetrics{
		connectionError: newDescWithLabel("connection_error_total", "The total number of connection error.", []string{"reason"}),
	}
)

type SoraConnectionErrorMetrics struct {
	connectionError *prometheus.Desc
}

func (m *SoraConnectionErrorMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.connectionError
}

func (m *SoraConnectionErrorMetrics) Collect(ch chan<- prometheus.Metric, report soraConnectionErrorReport) {
	ch <- newCounter(m.connectionError, float64(report.SdpGenerationError), "sdp")
	ch <- newCounter(m.connectionError, float64(report.SignalingError), "signaling")
}
