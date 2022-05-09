package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	soraClusterMetrics = SoraClusterMetrics{
		clusterNode: newDescWithLabel("cluster_node", "The sora server known cluster node.", []string{"node_name", "mode"}),
	}
)

type SoraClusterMetrics struct {
	clusterNode *prometheus.Desc
}

func (m *SoraClusterMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.clusterNode
}

func (m *SoraClusterMetrics) Collect(ch chan<- prometheus.Metric, nodeList []soraClusterNode) {
	for _, node := range nodeList {
		if node.ClusterNodeName != nil {
			ch <- newGauge(m.clusterNode, 1, *node.ClusterNodeName, *node.Mode)
		} else {
			ch <- newGauge(m.clusterNode, 1, *node.NodeName, *node.Mode)
		}
	}
}
