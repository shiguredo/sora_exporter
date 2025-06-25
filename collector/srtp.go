package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	srtpMetrics = SrtpMetrics{
		totalReceivedSrtp:          newDesc("srtp_received_packets_total", "The total number of received SRTP packets."),
		totalReceivedSrtpByteSize:  newDesc("srtp_received_bytes_total", "The total number of received SRTP bytes."),
		totalSentSrtp:              newDesc("srtp_sent_packets_total", "The total number of sent SRTP packets."),
		totalSentSrtpByteSize:      newDesc("srtp_sent_bytes_total", "The total number of sent SRTP bytes."),
		totalSentSrtpSfuDelayUs:    newDesc("srtp_sent_sfu_delay_us_total", "The total delay introduced by the SFU during the transfer of SRTP packets."),
		totalDecryptedSrtp:         newDesc("srtp_decrypted_packets_total", "The total number of decrpyted SRTP packets."),
		totalDecryptedSrtpByteSize: newDesc("srtp_decrpyted_bytes_total", "The total number of decrypted SRTP bytes."),
	}
)

type SrtpMetrics struct {
	totalReceivedSrtp          *prometheus.Desc
	totalReceivedSrtpByteSize  *prometheus.Desc
	totalSentSrtp              *prometheus.Desc
	totalSentSrtpByteSize      *prometheus.Desc
	totalSentSrtpSfuDelayUs    *prometheus.Desc
	totalDecryptedSrtp         *prometheus.Desc
	totalDecryptedSrtpByteSize *prometheus.Desc
}

func (m *SrtpMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.totalReceivedSrtp
	ch <- m.totalReceivedSrtpByteSize
	ch <- m.totalSentSrtp
	ch <- m.totalSentSrtpByteSize
	ch <- m.totalSentSrtpSfuDelayUs
	ch <- m.totalDecryptedSrtp
	ch <- m.totalDecryptedSrtpByteSize
}

func (m *SrtpMetrics) Collect(ch chan<- prometheus.Metric, report soraSrtpReport) {
	ch <- newCounter(m.totalReceivedSrtp, float64(report.TotalReceivedSrtp))
	ch <- newCounter(m.totalReceivedSrtpByteSize, float64(report.TotalReceivedSrtpByteSize))
	ch <- newCounter(m.totalSentSrtp, float64(report.TotalSentSrtp))
	ch <- newCounter(m.totalSentSrtpByteSize, float64(report.TotalSentSrtpByteSize))
	ch <- newCounter(m.totalSentSrtpSfuDelayUs, float64(report.TotalSentSrtpSfuDelayUs))
	ch <- newCounter(m.totalDecryptedSrtp, float64(report.TotalDecryptedSrtp))
	ch <- newCounter(m.totalDecryptedSrtpByteSize, float64(report.TotalDecryptedSrtpByteSize))
}
