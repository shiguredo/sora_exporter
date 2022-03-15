package main

import "github.com/prometheus/client_golang/prometheus"

var (
	connectionMetrics = ConnectionMetrics{
		totalConnectionCreated:     newDesc("connections_created_total", "The total number of connections created."),
		totalConnectionUpdated:     newDesc("connections_updated_total", "The total number of connections updated."),
		totalConnectionDestroyed:   newDesc("connections_destroyed_total", "The total number of connections destryed."),
		totalSuccessfulConnections: newDesc("successfull_connections_total", "The total number of successfull connections."),
		totalOngoingConnections:    newDesc("ongoing_connections_total", "The total number of ongoing connections."),
		totalFailedConnections:     newDesc("failed_connections_total", "The total number of failed connections."),
		totalDurationSec:           newDesc("duration_seconds_total", "The total duration of connections."),
		totalTurnUdpConnections:    newDesc("turn_udp_connections_total", "The total number of connections with TURN-UDP."),
		totalTurnTcpConnections:    newDesc("turn_tcp_connections_total", "The total number of connections with TURN-TCP."),
		averageDurationSec:         newDesc("average_duration_seconds", "The average connection duration in seconds."),
		averageSetupTimeSec:        newDesc("average_setup_time_seconds", "The average setup time in seconds."),
	}
)

type ConnectionMetrics struct {
	totalConnectionCreated     *prometheus.Desc
	totalConnectionUpdated     *prometheus.Desc
	totalConnectionDestroyed   *prometheus.Desc
	totalSuccessfulConnections *prometheus.Desc
	totalOngoingConnections    *prometheus.Desc
	totalFailedConnections     *prometheus.Desc
	totalDurationSec           *prometheus.Desc
	totalTurnUdpConnections    *prometheus.Desc
	totalTurnTcpConnections    *prometheus.Desc
	averageDurationSec         *prometheus.Desc
	averageSetupTimeSec        *prometheus.Desc
}

func (m *ConnectionMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.totalConnectionCreated
	ch <- m.totalConnectionUpdated
	ch <- m.totalConnectionUpdated
	ch <- m.totalSuccessfulConnections
	ch <- m.totalOngoingConnections
	ch <- m.totalFailedConnections
	ch <- m.totalDurationSec
	ch <- m.totalTurnUdpConnections
	ch <- m.totalTurnTcpConnections
	ch <- m.averageDurationSec
	ch <- m.averageSetupTimeSec
}

func (m *ConnectionMetrics) Collect(ch chan<- prometheus.Metric, report soraConnectionReport) {
	ch <- newCounter(m.totalConnectionCreated, float64(report.TotalConnectionCreated))
	ch <- newCounter(m.totalConnectionUpdated, float64(report.TotalConnectionUpdated))
	ch <- newCounter(m.totalConnectionDestroyed, float64(report.TotalConnectionDestroyed))
	ch <- newCounter(m.totalSuccessfulConnections, float64(report.TotalSuccessfulConnections))
	ch <- newGauge(m.totalOngoingConnections, float64(report.TotalOngoingConnections))
	ch <- newCounter(m.totalFailedConnections, float64(report.TotalFailedConnections))
	ch <- newCounter(m.totalDurationSec, float64(report.TotalDurationSec))
	ch <- newCounter(m.totalTurnUdpConnections, float64(report.TotalTurnUdpConnections))
	ch <- newCounter(m.totalTurnTcpConnections, float64(report.TotalTurnTcpConnections))
	ch <- newGauge(m.averageDurationSec, float64(report.AverageDurationSec))
	ch <- newGauge(m.averageSetupTimeSec, float64(report.AverageSetupTimeMsec/1000))
}
