// TODO: rename file
package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	soraClusterMetrics = SoraClusterMetrics{
		clusterNode: newDescWithLabel("cluster_node", "The sora server known cluster node.", []string{"node_name", "mode"}),
		raftRole:        newDescWithLabel("cluster_raft_role", "The current Raft role. The role name is indicated by the label 'role'. The value of this metric is always set to 1.", []string{"role"}),
		raftTerm:        newDesc("cluster_raft_term", "The current Raft term."),
		raftCommitIndex: newDesc("cluster_raft_commit_index", "The latest committed Raft log index."),
	}
)

type SoraClusterMetrics struct {
	clusterNode     *prometheus.Desc
	raftRole        *prometheus.Desc
	raftTerm        *prometheus.Desc
	raftCommitIndex *prometheus.Desc
}

func (m *SoraClusterMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.clusterNode
	ch <- m.raftRole
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
	ch <- newGauge(m.raftRole, 1.0, report.RaftRole)
	ch <- newCounter(m.raftTerm, float64(report.RaftTerm))
	ch <- newCounter(m.raftCommitIndex, float64(report.RaftCommitIndex))
}
