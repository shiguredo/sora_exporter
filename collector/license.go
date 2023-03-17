package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	licenseMetrics = LicenseMetrics{
		licenseInfo:           newDescWithLabel("license_info", "sora license info.", []string{"expired_at", "product_name", "serial_code", "type"}),
		licenseMaxConnections: newDesc("license_max_connections", "sora license file's max_connections."),
		licenseMaxNodes:       newDesc("license_max_nodes", "sora license file's max_nodes."),
	}
)

type LicenseMetrics struct {
	licenseInfo           *prometheus.Desc
	licenseMaxConnections *prometheus.Desc
	licenseMaxNodes       *prometheus.Desc
}

func (m *LicenseMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.licenseInfo
	ch <- m.licenseMaxConnections
	ch <- m.licenseMaxNodes
}

func (m *LicenseMetrics) Collect(ch chan<- prometheus.Metric, info soraLicenseInfo) {
	ch <- newGauge(m.licenseInfo, 1, info.ExpiredAt, info.ProductName, info.SerialCode, info.Type)
	ch <- newGauge(m.licenseMaxConnections, float64(info.MaxConnections))
	if info.MaxNodes != nil {
		ch <- newGauge(m.licenseMaxNodes, float64(*info.MaxNodes))
	}
}
