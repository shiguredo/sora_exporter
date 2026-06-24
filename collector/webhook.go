package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	webhookMetrics = WebhookMetrics{
		totalSuccessfulAuthWebhook:        newDescWithLabel("successful_auth_webhook_total", "The total number of successful auth webhook.", []string{"state"}),
		totalAuthWebhook:                  newDescWithLabel("auth_webhook_total", "The total number of auth webhook.", []string{"state"}),
		totalSessionWebhook:               newDescWithLabel("session_webhook_total", "The total number of session webhook.", []string{"state"}),
		totalEventWebhook:                 newDescWithLabel("event_webhook_total", "The total number of event webhook.", []string{"state"}),
		totalStatsWebhook:                 newDescWithLabel("stats_webhook_total", "The total number of stats webhook.", []string{"state"}),
		authWebhookResponseTimeSeconds:    newDesc("auth_webhook_response_time_seconds", "Response time histogram for auth webhook requests."),
		sessionWebhookResponseTimeSeconds: newDesc("session_webhook_response_time_seconds", "Response time histogram for session webhook requests."),
		eventWebhookResponseTimeSeconds:   newDesc("event_webhook_response_time_seconds", "Response time histogram for event webhook requests."),
		statsWebhookResponseTimeSeconds:   newDesc("stats_webhook_response_time_seconds", "Response time histogram for stats webhook requests."),
	}
)

type WebhookMetrics struct {
	totalSuccessfulAuthWebhook        *prometheus.Desc
	totalAuthWebhook                  *prometheus.Desc
	totalSessionWebhook               *prometheus.Desc
	totalEventWebhook                 *prometheus.Desc
	totalStatsWebhook                 *prometheus.Desc
	authWebhookResponseTimeSeconds    *prometheus.Desc
	sessionWebhookResponseTimeSeconds *prometheus.Desc
	eventWebhookResponseTimeSeconds   *prometheus.Desc
	statsWebhookResponseTimeSeconds   *prometheus.Desc
}

func (m *WebhookMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.totalSuccessfulAuthWebhook
	ch <- m.totalAuthWebhook
	ch <- m.totalSessionWebhook
	ch <- m.totalEventWebhook
	ch <- m.totalStatsWebhook
	ch <- m.authWebhookResponseTimeSeconds
	ch <- m.sessionWebhookResponseTimeSeconds
	ch <- m.eventWebhookResponseTimeSeconds
	ch <- m.statsWebhookResponseTimeSeconds
}

func (m *WebhookMetrics) Collect(ch chan<- prometheus.Metric, report soraWebhookReport) {
	ch <- newCounter(m.totalSuccessfulAuthWebhook, float64(report.TotalAuthWebhookAllowed), "allowed")
	ch <- newCounter(m.totalSuccessfulAuthWebhook, float64(report.TotalAuthWebhookDenied), "denied")
	ch <- newCounter(m.totalAuthWebhook, float64(report.TotalSuccessfulAuthWebhook), "successful")
	ch <- newCounter(m.totalAuthWebhook, float64(report.TotalFailedAuthWebhook), "failed")
	ch <- newCounter(m.totalSessionWebhook, float64(report.TotalSuccessfulSessionWebhook), "successful")
	ch <- newCounter(m.totalSessionWebhook, float64(report.TotalFailedSessionWebhook), "failed")
	ch <- newCounter(m.totalSessionWebhook, float64(report.TotalIgnoredSessionWebhook), "ignored")
	ch <- newCounter(m.totalEventWebhook, float64(report.TotalSuccessfulEventWebhook), "successful")
	ch <- newCounter(m.totalEventWebhook, float64(report.TotalFailedEventWebhook), "failed")
	ch <- newCounter(m.totalEventWebhook, float64(report.TotalIgnoredEventWebhook), "ignored")
	ch <- newCounter(m.totalStatsWebhook, float64(report.TotalSuccessfulStatsWebhook), "successful")
	ch <- newCounter(m.totalStatsWebhook, float64(report.TotalFailedStatsWebhook), "failed")
	ch <- newCounter(m.totalStatsWebhook, float64(report.TotalIgnoredStatsWebhook), "ignored")

	if len(report.AuthWebhookResponseTimeMsBuckets) > 0 {
		count := uint64(report.TotalSuccessfulAuthWebhook + report.TotalFailedAuthWebhook)
		sum := float64(report.TotalAuthWebhookResponseTimeMs) / 1000.0
		ch <- newWebhookResponseTimeHistogram(m.authWebhookResponseTimeSeconds, count, sum, report.AuthWebhookResponseTimeMsBuckets)
	}
	if len(report.SessionWebhookResponseTimeMsBuckets) > 0 {
		count := uint64(report.TotalSuccessfulSessionWebhook + report.TotalFailedSessionWebhook)
		sum := float64(report.TotalSessionWebhookResponseTimeMs) / 1000.0
		ch <- newWebhookResponseTimeHistogram(m.sessionWebhookResponseTimeSeconds, count, sum, report.SessionWebhookResponseTimeMsBuckets)
	}
	if len(report.EventWebhookResponseTimeMsBuckets) > 0 {
		count := uint64(report.TotalSuccessfulEventWebhook + report.TotalFailedEventWebhook)
		sum := float64(report.TotalEventWebhookResponseTimeMs) / 1000.0
		ch <- newWebhookResponseTimeHistogram(m.eventWebhookResponseTimeSeconds, count, sum, report.EventWebhookResponseTimeMsBuckets)
	}
	if len(report.StatsWebhookResponseTimeMsBuckets) > 0 {
		count := uint64(report.TotalSuccessfulStatsWebhook + report.TotalFailedStatsWebhook)
		sum := float64(report.TotalStatsWebhookResponseTimeMs) / 1000.0
		ch <- newWebhookResponseTimeHistogram(m.statsWebhookResponseTimeSeconds, count, sum, report.StatsWebhookResponseTimeMsBuckets)
	}
}

func newWebhookResponseTimeHistogram(d *prometheus.Desc, count uint64, sum float64, buckets []webhookResponseTimeBucket) prometheus.Metric {
	bucketMap := make(map[float64]uint64, len(buckets))
	for _, b := range buckets {
		bucketMap[float64(b.UpperBound)/1000.0] = uint64(b.Count)
	}
	return prometheus.MustNewConstHistogram(d, count, sum, bucketMap)
}
