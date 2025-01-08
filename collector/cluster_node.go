package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	soraClusterMetrics = SoraClusterMetrics{
		clusterNode:     newDescWithLabel("cluster_node", "The sora server known cluster node.", []string{"node_name", "mode", "node_type"}),
		raftState:       newDescWithLabel("cluster_raft_state", "The current Raft state. The state name is indicated by the label 'state'. The value of this metric is always set to 1.", []string{"state"}),
		raftTerm:        newDesc("cluster_raft_term", "The current Raft term."),
		raftCommitIndex: newDesc("cluster_raft_commit_index", "The latest committed Raft log index."),

		clusterRelayReceivedBytesTotal:   newDescWithLabel("cluster_relay_received_bytes_total", "The total number of bytes received by the cluster relay.", []string{"node_name"}),
		clusterRelaySentBytesTotal:       newDescWithLabel("cluster_relay_sent_bytes_total", "The total number of bytes sent by the cluster relay.", []string{"node_name"}),
		clusterRelayReceivedPacketsTotal: newDescWithLabel("cluster_relay_received_packets_total", "The total number of packets received by the cluster relay.", []string{"node_name"}),
		clusterRelaySentPacketsTotal:     newDescWithLabel("cluster_relay_sent_packets_total", "The total number of packets sent by the cluster relay.", []string{"node_name"}),

		clusterRelayPlumtreeSentGossipTotal:   newDescWithLabel("cluster_relay_plumtree_sent_gossip_total", "The total number of Plumtree GOSSIP messages sent by the cluster relay.", []string{"node_name"}),
		clusterRelayPlumtreeReceivedGossipTotal:   newDescWithLabel("cluster_relay_plumtree_received_gossip_total", "The total number of Plumtree GOSSIP messages received by the cluster relay.", []string{"node_name"}),
		clusterRelayPlumtreeReceivedGossipHopTotal:   newDescWithLabel("cluster_relay_plumtree_received_gossip_hop_total", "The total number of hop count of Plumtree GOSSIP messages received by the cluster relay.", []string{"node_name"}),
		clusterRelayPlumtreeSentIhaveTotal:   newDescWithLabel("cluster_relay_plumtree_sent_ihave_total", "The total number of Plumtree IHAVE messages sent by the cluster relay.", []string{"node_name"}),
		clusterRelayPlumtreeReceivedIhaveTotal:   newDescWithLabel("cluster_relay_plumtree_received_ihave_total", "The total number of Plumtree IHAVE messages received by the cluster relay.", []string{"node_name"}),
		clusterRelayPlumtreeSentGraftTotal:   newDescWithLabel("cluster_relay_plumtree_sent_graft_total", "The total number of Plumtree GRAFT messages sent by the cluster relay.", []string{"node_name"}),
		clusterRelayPlumtreeReceivedGraftTotal:   newDescWithLabel("cluster_relay_plumtree_received_graft_total", "The total number of Plumtree GRAFT messages received by the cluster relay.", []string{"node_name"}),
		clusterRelayPlumtreeSentPruneTotal:   newDescWithLabel("cluster_relay_plumtree_sent_prune_total", "The total number of Plumtree PRUNE messages sent by the cluster relay.", []string{"node_name"}),
		clusterRelayPlumtreeReceivedPruneTotal:   newDescWithLabel("cluster_relay_plumtree_received_prune_total", "The total number of Plumtree PRUNE messages received by the cluster relay.", []string{"node_name"}),
		clusterRelayPlumtreeGraftMissTotal:   newDescWithLabel("cluster_relay_plumtree_graft_miss_total", "The total number of Plumtree GRAFT messages missed by the cluster relay.", []string{"node_name"}),
		clusterRelayPlumtreeSkippedSendTotal:   newDescWithLabel("cluster_relay_plumtree_skipped_send_total", "The total number of Plumtree messages whose sending was skipped by the cluster relay.", []string{"node_name"}),
		clusterRelayPlumtreeIgnoredTotal:   newDescWithLabel("cluster_relay_plumtree_ignored_total", "The total number of Plumtree messages received but ignored by the cluster relay.", []string{"node_name"}),
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

	clusterRelayPlumtreeSentGossipTotal        *prometheus.Desc
	clusterRelayPlumtreeReceivedGossipTotal    *prometheus.Desc
	clusterRelayPlumtreeReceivedGossipHopTotal *prometheus.Desc
	clusterRelayPlumtreeSentIhaveTotal         *prometheus.Desc
	clusterRelayPlumtreeReceivedIhaveTotal     *prometheus.Desc
	clusterRelayPlumtreeSentGraftTotal         *prometheus.Desc
	clusterRelayPlumtreeReceivedGraftTotal     *prometheus.Desc
	clusterRelayPlumtreeSentPruneTotal         *prometheus.Desc
	clusterRelayPlumtreeReceivedPruneTotal     *prometheus.Desc
	clusterRelayPlumtreeGraftMissTotal         *prometheus.Desc
	clusterRelayPlumtreeSkippedSendTotal       *prometheus.Desc
	clusterRelayPlumtreeIgnoredTotal           *prometheus.Desc
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
	ch <- m.clusterRelayPlumtreeSentGossipTotal
	ch <- m.clusterRelayPlumtreeReceivedGossipTotal
	ch <- m.clusterRelayPlumtreeReceivedGossipHopTotal
	ch <- m.clusterRelayPlumtreeSentIhaveTotal
	ch <- m.clusterRelayPlumtreeReceivedIhaveTotal
	ch <- m.clusterRelayPlumtreeSentGraftTotal
	ch <- m.clusterRelayPlumtreeReceivedGraftTotal
	ch <- m.clusterRelayPlumtreeSentPruneTotal
	ch <- m.clusterRelayPlumtreeReceivedPruneTotal
	ch <- m.clusterRelayPlumtreeGraftMissTotal
	ch <- m.clusterRelayPlumtreeSkippedSendTotal
	ch <- m.clusterRelayPlumtreeIgnoredTotal
}

func (m *SoraClusterMetrics) Collect(ch chan<- prometheus.Metric, nodeList []soraClusterNode, report soraClusterReport, clusterRelaies []soraClusterRelay) {
	for _, node := range nodeList {
		value := 0.0
		if node.Connected {
			value = 1.0
		}
		// connected の状態によらず、基本の状態は regular
		// temporary_node が true の場合だけ temporary になる
		nodeType := "regular"
		if node.TemporaryNode {
			nodeType = "temporary"
		}
		if node.ClusterNodeName != "" {
			ch <- newGauge(m.clusterNode, value, node.ClusterNodeName, node.Mode, nodeType)
		} else {
			ch <- newGauge(m.clusterNode, value, node.NodeName, node.Mode, nodeType)
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

		ch <- newCounter(m.clusterRelayPlumtreeSentGossipTotal, float64(relayNode.Plumtree.TotalSentGossip), relayNode.NodeName)
		ch <- newCounter(m.clusterRelayPlumtreeReceivedGossipTotal, float64(relayNode.Plumtree.TotalReceivedGossip), relayNode.NodeName)
		ch <- newCounter(m.clusterRelayPlumtreeReceivedGossipHopTotal, float64(relayNode.Plumtree.TotalReceivedGossipHop), relayNode.NodeName)
		ch <- newCounter(m.clusterRelayPlumtreeSentIhaveTotal, float64(relayNode.Plumtree.TotalSentIhave), relayNode.NodeName)
		ch <- newCounter(m.clusterRelayPlumtreeReceivedIhaveTotal, float64(relayNode.Plumtree.TotalReceivedIhave), relayNode.NodeName)
		ch <- newCounter(m.clusterRelayPlumtreeSentGraftTotal, float64(relayNode.Plumtree.TotalSentGraft), relayNode.NodeName)
		ch <- newCounter(m.clusterRelayPlumtreeReceivedGraftTotal, float64(relayNode.Plumtree.TotalReceivedGraft), relayNode.NodeName)
		ch <- newCounter(m.clusterRelayPlumtreeSentPruneTotal, float64(relayNode.Plumtree.TotalSentPrune), relayNode.NodeName)
		ch <- newCounter(m.clusterRelayPlumtreeReceivedPruneTotal, float64(relayNode.Plumtree.TotalReceivedPrune), relayNode.NodeName)
		ch <- newCounter(m.clusterRelayPlumtreeGraftMissTotal, float64(relayNode.Plumtree.TotalGraftMiss), relayNode.NodeName)
		ch <- newCounter(m.clusterRelayPlumtreeSkippedSendTotal, float64(relayNode.Plumtree.TotalSkippedSend), relayNode.NodeName)
		ch <- newCounter(m.clusterRelayPlumtreeIgnoredTotal, float64(relayNode.Plumtree.TotalIgnored), relayNode.NodeName)
	}
}
