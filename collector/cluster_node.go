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

		clusterRelayReceivedBytes:   newDescWithLabel("cluster_relay_received_bytes", "The total number of bytes received by the cluster relay.", []string{"node_name"}),
		clusterRelaySentBytes:       newDescWithLabel("cluster_relay_sent_bytes", "The total number of bytes sent by the cluster relay.", []string{"node_name"}),
		clusterRelayReceivedPackets: newDescWithLabel("cluster_relay_received_packets", "The total number of packets received by the cluster relay.", []string{"node_name"}),
		clusterRelaySentPackets:     newDescWithLabel("cluster_relay_sent_packets", "The total number of packets sent by the cluster relay.", []string{"node_name"}),
	}
)

type SoraClusterMetrics struct {
	clusterNode     *prometheus.Desc
	raftState       *prometheus.Desc
	raftTerm        *prometheus.Desc
	raftCommitIndex *prometheus.Desc

	clusterRelayReceivedBytes   *prometheus.Desc
	clusterRelaySentBytes       *prometheus.Desc
	clusterRelayReceivedPackets *prometheus.Desc
	clusterRelaySentPackets     *prometheus.Desc
}

func (m *SoraClusterMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.clusterNode
	ch <- m.raftState
	ch <- m.raftTerm
	ch <- m.raftCommitIndex
	ch <- m.clusterRelayReceivedBytes
	ch <- m.clusterRelaySentBytes
	ch <- m.clusterRelayReceivedPackets
	ch <- m.clusterRelaySentPackets
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
		ch <- newCounter(m.clusterRelayReceivedBytes, float64(relayNode.TotalReceivedByteSize), relayNode.NodeName)
		ch <- newCounter(m.clusterRelaySentBytes, float64(relayNode.TotalSentByteSize), relayNode.NodeName)
		ch <- newCounter(m.clusterRelayReceivedPackets, float64(relayNode.TotalReceived), relayNode.NodeName)
		ch <- newCounter(m.clusterRelaySentPackets, float64(relayNode.TotalSent), relayNode.NodeName)
	}
}
