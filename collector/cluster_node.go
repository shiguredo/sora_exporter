package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	soraClusterMetrics = SoraClusterMetrics{
		clusterNode:     newDescWithLabel("cluster_node", "The sora server known cluster node.", []string{"node_name", "mode"}),
		raftState:       newDescWithLabel("cluster_raft_state", "The current Raft state. The state name is indicated by the label 'state'. The value of this metric is always set to 1.", []string{"state"}),
		raftTerm:        newDesc("cluster_raft_term", "The current Raft term."),
		raftCommitIndex: newDesc("cluster_raft_commit_index", "The latest committed Raft log index."),
	}
)

type SoraClusterMetrics struct {
	clusterNode     *prometheus.Desc
	raftState       *prometheus.Desc
	raftTerm        *prometheus.Desc
	raftCommitIndex *prometheus.Desc
}

func (m *SoraClusterMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.clusterNode
	ch <- m.raftState
	ch <- m.raftTerm
	ch <- m.raftCommitIndex
}

func (m *SoraClusterMetrics) Collect(ch chan<- prometheus.Metric, nodeList []soraClusterNode, report soraClusterReport) {
	for _, node := range nodeList {
		if node.ClusterNodeName != nil {
			ch <- newGauge(m.clusterNode, 1, *node.ClusterNodeName, *node.Mode)
		} else {
			ch <- newGauge(m.clusterNode, 1, *node.NodeName, *node.Mode)
		}
	}
	if report.RaftState == "" {
		ch <- newGauge(m.raftState, 1.0, "undefined")
	} else {
		ch <- newGauge(m.raftState, 1.0, report.RaftState)
	}
	ch <- newCounter(m.raftTerm, float64(report.RaftTerm))
	ch <- newCounter(m.raftCommitIndex, float64(report.RaftCommitIndex))
}
