package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	sctpMetrics = SctpMetrics{
		totalReceivedSctp:         newDesc("sctp_received_packets_total", "The total number of received SCTP packets."),
		totalReceivedSctpByteSize: newDesc("sctp_received_bytes_total", "The total number of received SCTP bytes."),
		totalSentSctp:             newDesc("sctp_sent_packets_total", "The total number of sent SRPT packets."),
		totalSentSctpByteSize:     newDesc("sctp_sent_bytes_total", "The total number of sent SCTP bytes."),
	}
)

type SctpMetrics struct {
	totalReceivedSctp         *prometheus.Desc
	totalReceivedSctpByteSize *prometheus.Desc
	totalSentSctp             *prometheus.Desc
	totalSentSctpByteSize     *prometheus.Desc
}

func (m *SctpMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.totalReceivedSctp
	ch <- m.totalReceivedSctpByteSize
	ch <- m.totalSentSctp
	ch <- m.totalSentSctpByteSize
}

func (m *SctpMetrics) Collect(ch chan<- prometheus.Metric, report soraSctpReport) {
	ch <- newCounter(m.totalReceivedSctp, float64(report.TotalReceivedSctp))
	ch <- newCounter(m.totalReceivedSctpByteSize, float64(report.TotalReceivedSctpByteSize))
	ch <- newCounter(m.totalSentSctp, float64(report.TotalSentSctp))
	ch <- newCounter(m.totalSentSctpByteSize, float64(report.TotalSentSctpByteSize))
}
