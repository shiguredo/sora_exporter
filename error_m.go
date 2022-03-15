package main

import "github.com/prometheus/client_golang/prometheus"

var (
	errorMetrics = ErrorMetrics{
		sdpGenerationError: newDesc("sdp_generation_error_total", "The total number of SDP genration error."),
		signalingError:     newDesc("signaling_error_total", "The total number of signaling error."),
	}
)

type ErrorMetrics struct {
	sdpGenerationError *prometheus.Desc
	signalingError     *prometheus.Desc
}

func (m *ErrorMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.sdpGenerationError
	ch <- m.signalingError
}

func (m *ErrorMetrics) Collect(ch chan<- prometheus.Metric, report soraErrorReport) {
	ch <- newCounter(m.sdpGenerationError, float64(report.SdpGenerationError))
	ch <- newCounter(m.signalingError, float64(report.SignalingError))
}
