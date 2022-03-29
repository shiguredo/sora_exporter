package collector

import "github.com/prometheus/client_golang/prometheus"

func newDesc(name, help string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName("sora", "exporter", name), help, nil, nil)
}

func newDescWithLabel(name, help string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName("sora", "exporter", name), help, labels, nil)
}

func newGauge(d *prometheus.Desc, v float64) prometheus.Metric {
	return prometheus.MustNewConstMetric(d, prometheus.GaugeValue, v)
}

func newCounter(d *prometheus.Desc, v float64) prometheus.Metric {
	return prometheus.MustNewConstMetric(d, prometheus.CounterValue, v)
}

func newInfo(d *prometheus.Desc, labelValues ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(d, prometheus.GaugeValue, 1, labelValues...)
}
