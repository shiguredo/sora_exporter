package collector

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	licenseMetrics = LicenseMetrics{
		licenseInfo:                      newDescWithLabel("license_info", "sora license info.", []string{"expired_at", "product_name", "serial_code", "type"}),
		licenseMaxConnections:            newDesc("license_max_connections", "sora license file's max_connections."),
		licenseMaxNodes:                  newDesc("license_max_nodes", "sora license file's max_nodes."),
		licenseExpiredAtTimestampSeconds: newDesc("license_expired_at_timestamp_seconds", "sora license file's expire seconds since epoch."),
	}
)

type LicenseMetrics struct {
	licenseInfo                      *prometheus.Desc
	licenseMaxConnections            *prometheus.Desc
	licenseMaxNodes                  *prometheus.Desc
	licenseExpiredAtTimestampSeconds *prometheus.Desc
}

func (m *LicenseMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.licenseInfo
	ch <- m.licenseMaxConnections
	ch <- m.licenseMaxNodes
	ch <- m.licenseExpiredAtTimestampSeconds
}

func (m *LicenseMetrics) Collect(ch chan<- prometheus.Metric, info *soraLicenseInfo) {
	ch <- newGauge(m.licenseInfo, 1, info.ExpiredAt, info.ProductName, info.SerialCode, info.Type)
	ch <- newGauge(m.licenseMaxConnections, float64(info.MaxConnections))
	if info.MaxNodes != nil {
		ch <- newGauge(m.licenseMaxNodes, float64(*info.MaxNodes))
	}

	expiredAtSec := expiredAtToSecondSinceEpoch(info.ExpiredAt)
	ch <- newGauge(m.licenseExpiredAtTimestampSeconds, expiredAtSec)
}

func expiredAtToSecondSinceEpoch(expiredAt string) float64 {
	// expiredAt: "2024-12"
	expiredAtTime, err := time.Parse("2006-01", expiredAt)
	if err != nil {
		return 0
	}

	// 期限翌月の 1 日 0 時 0 分0 秒の 1 秒前
	expiredTimestamp := time.Date(expiredAtTime.Year(), expiredAtTime.Month()+1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second)
	return float64(expiredTimestamp.UnixNano()) / 1e9
}
