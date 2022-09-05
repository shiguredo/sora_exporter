package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	webhookMetrics = WebhookMetrics{
		totalSuccessfulAuthWebhook: newDescWithLabel("successful_auth_webhook_total", "The total number of successful auth webhook.", []string{"state"}),
		totalAuthWebhook:           newDescWithLabel("auth_webhook_total", "The total number of auth webhook.", []string{"state"}),
		totalSessionWebhook:        newDescWithLabel("session_webhook_total", "The total number of session webhook.", []string{"state"}),
		totalEventWebhook:          newDescWithLabel("event_webhook_total", "The total number of event webhook.", []string{"state"}),
	}
)

type WebhookMetrics struct {
	totalSuccessfulAuthWebhook *prometheus.Desc
	totalAuthWebhook           *prometheus.Desc
	totalSessionWebhook        *prometheus.Desc
	totalEventWebhook          *prometheus.Desc
}

func (m *WebhookMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.totalSuccessfulAuthWebhook
	ch <- m.totalAuthWebhook
	ch <- m.totalSessionWebhook
	ch <- m.totalEventWebhook
}

func (m *WebhookMetrics) Collect(ch chan<- prometheus.Metric, report soraWebhookReport) {
	ch <- newCounter(m.totalSuccessfulAuthWebhook, float64(report.TotalAuthWebhookAllowed), "allowed")
	ch <- newCounter(m.totalSuccessfulAuthWebhook, float64(report.TotalAuthWebhookDenied), "denied")
	ch <- newCounter(m.totalAuthWebhook, float64(report.TotalSuccessfulAuthWebhook), "successful")
	ch <- newCounter(m.totalAuthWebhook, float64(report.TotalFailedAuthWebhook), "failed")
	ch <- newCounter(m.totalSessionWebhook, float64(report.TotalSuccessfulSessionWebhook), "successful")
	ch <- newCounter(m.totalSessionWebhook, float64(report.TotalFailedSessionWebhook), "failed")
	ch <- newCounter(m.totalEventWebhook, float64(report.TotalSuccessfulEventWebhook), "successful")
	ch <- newCounter(m.totalEventWebhook, float64(report.TotalFailedEventWebhook), "failed")
}
