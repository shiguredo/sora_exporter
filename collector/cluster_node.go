package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	soraClusterMetrics = SoraClusterMetrics{
		clusterNode:     newDescWithLabel("cluster_node", "The sora server known cluster node.", []string{"node_name", "mode"}),
		raftState:       newDescWithLabel("cluster_raft_state", "The current Raft state. The state name is indicated by the label 'state'. The value of this metric is always set to 1.", []string{"state"}),
		raftTerm:        newDesc("cluster_raft_term", "The current Raft term."),
		raftCommitIndex: newDesc("cluster_raft_commit_index", "The latest committed Raft log index."),

		clusterRelayReceivedBytesTotal:   newDescWithLabel("cluster_relay_received_bytes_total", "The total number of bytes received by the cluster relay.", []string{"node_name"}),
		clusterRelaySentBytesTotal:       newDescWithLabel("cluster_relay_sent_bytes_total", "The total number of bytes sent by the cluster relay.", []string{"node_name"}),
		clusterRelayReceivedPacketsTotal: newDescWithLabel("cluster_relay_received_packets_total", "The total number of packets received by the cluster relay.", []string{"node_name"}),
		clusterRelaySentPacketsTotal:     newDescWithLabel("cluster_relay_sent_packets_total", "The total number of packets sent by the cluster relay.", []string{"node_name"}),
	}
)

type SoraClusterMetrics struct {
	clusterNode     *prometheus.Desc
	raftState       *prometheus.Desc
	raftTerm        *prometheus.Desc
	raftCommitIndex *prometheus.Desc

	clusterRelayReceivedBytesTotal   *prometheus.Desc
	clusterRelaySentBytesTotal       *prometheus.Desc
	clusterRelayReceivedPacketsTotal *prometheus.Desc
	clusterRelaySentPacketsTotal     *prometheus.Desc
}

func (m *SoraClusterMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.clusterNode
	ch <- m.raftState
	ch <- m.raftTerm
	ch <- m.raftCommitIndex
	ch <- m.clusterRelayReceivedBytesTotal
	ch <- m.clusterRelaySentBytesTotal
	ch <- m.clusterRelayReceivedPacketsTotal
	ch <- m.clusterRelaySentPacketsTotal
}

func (m *SoraClusterMetrics) Collect(ch chan<- prometheus.Metric, nodeList []soraClusterNode, report soraClusterReport, clusterRelaies []soraClusterRelay) {
	for _, node := range nodeList {
		value := 0.0
		if node.Connected {
			value = 1.0
		}
		if node.ClusterNodeName != "" {
			ch <- newGauge(m.clusterNode, value, node.ClusterNodeName, node.Mode)
		} else {
			ch <- newGauge(m.clusterNode, value, node.NodeName, node.Mode)
		}
	}
	ch <- newGauge(m.raftState, 1.0, report.RaftState)
	ch <- newCounter(m.raftTerm, float64(report.RaftTerm))
	ch <- newCounter(m.raftCommitIndex, float64(report.RaftCommitIndex))

	for _, relayNode := range clusterRelaies {
		ch <- newCounter(m.clusterRelayReceivedBytesTotal, float64(relayNode.TotalReceivedByteSize), relayNode.NodeName)
		ch <- newCounter(m.clusterRelaySentBytesTotal, float64(relayNode.TotalSentByteSize), relayNode.NodeName)
		ch <- newCounter(m.clusterRelayReceivedPacketsTotal, float64(relayNode.TotalReceived), relayNode.NodeName)
		ch <- newCounter(m.clusterRelaySentPacketsTotal, float64(relayNode.TotalSent), relayNode.NodeName)
	}
}
