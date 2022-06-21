package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	connectionMetrics = ConnectionMetrics{
		totalConnection:                   newDescWithLabel("connections_total", "The total number of connections created.", []string{"state"}),
		totalOngoingConnections:           newDesc("ongoing_connections_total", "The total number of ongoing connections."),
		totalDurationSec:                  newDesc("duration_seconds_total", "The total duration of connections."),
		totalTurnConnections:              newDescWithLabel("turn_connections_total", "The total number of connections with TURN.", []string{"proto"}),
		averageDurationSec:                newDesc("average_duration_seconds", "The average connection duration in seconds."),
		averageSetupTimeSec:               newDesc("average_setup_time_seconds", "The average setup time in seconds."),
		totalSession:                      newDescWithLabel("session_total", "The total number of session.", []string{"state"}),
		totalReceivedInvalidTurnTCPPacket: newDesc("received_invalid_turn_tcp_packet_total", "The total number of invalid packets with TURN-TCP"),
	}
)

type ConnectionMetrics struct {
	totalConnection                   *prometheus.Desc
	totalOngoingConnections           *prometheus.Desc
	totalDurationSec                  *prometheus.Desc
	totalTurnConnections              *prometheus.Desc
	averageDurationSec                *prometheus.Desc
	averageSetupTimeSec               *prometheus.Desc
	totalSession                      *prometheus.Desc
	totalReceivedInvalidTurnTCPPacket *prometheus.Desc
}

func (m *ConnectionMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.totalConnection
	ch <- m.totalOngoingConnections
	ch <- m.totalDurationSec
	ch <- m.totalTurnConnections
	ch <- m.averageDurationSec
	ch <- m.averageSetupTimeSec
	ch <- m.totalSession
	ch <- m.totalReceivedInvalidTurnTCPPacket
}

func (m *ConnectionMetrics) Collect(ch chan<- prometheus.Metric, report soraConnectionReport) {
	ch <- newCounter(m.totalConnection, float64(report.TotalConnectionCreated), "created")
	ch <- newCounter(m.totalConnection, float64(report.TotalConnectionUpdated), "updated")
	ch <- newCounter(m.totalConnection, float64(report.TotalConnectionDestroyed), "destroyed")
	ch <- newCounter(m.totalConnection, float64(report.TotalSuccessfulConnections), "successful")
	ch <- newCounter(m.totalConnection, float64(report.TotalFailedConnections), "failed")
	ch <- newGauge(m.totalOngoingConnections, float64(report.TotalOngoingConnections))
	ch <- newCounter(m.totalDurationSec, float64(report.TotalDurationSec))
	ch <- newCounter(m.totalTurnConnections, float64(report.TotalTurnUDPConnections), "udp")
	ch <- newCounter(m.totalTurnConnections, float64(report.TotalTurnTCPConnections), "tcp")
	ch <- newGauge(m.averageDurationSec, float64(report.AverageDurationSec))
	ch <- newGauge(m.averageSetupTimeSec, float64(report.AverageSetupTimeMsec/1000))
	ch <- newCounter(m.totalSession, float64(report.TotalSessionCreated), "created")
	ch <- newCounter(m.totalSession, float64(report.TotalSessionDestroyed), "destroyed")
	ch <- newCounter(m.totalReceivedInvalidTurnTCPPacket, float64(report.TotalReceivedInvalidTurnTCPPacket))
}
