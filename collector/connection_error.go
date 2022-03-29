package collector

import "github.com/prometheus/client_golang/prometheus"

// この統計情報はアンドキュメントです
var (
	soraConnectionErrorMetrics = SoraConnectionErrorMetrics{
		sdpGenerationError: newDesc("sdp_generation_error_total", "The total number of SDP genration error."),
		signalingError:     newDesc("signaling_error_total", "The total number of signaling error."),
	}
)

type SoraConnectionErrorMetrics struct {
	sdpGenerationError *prometheus.Desc
	signalingError     *prometheus.Desc
}

func (m *SoraConnectionErrorMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.sdpGenerationError
	ch <- m.signalingError
}

func (m *SoraConnectionErrorMetrics) Collect(ch chan<- prometheus.Metric, report soraConnectionErrorReport) {
	ch <- newCounter(m.sdpGenerationError, float64(report.SdpGenerationError))
	ch <- newCounter(m.signalingError, float64(report.SignalingError))
}
