package collector

import "github.com/prometheus/client_golang/prometheus"

// この統計情報はアンドキュメントです
var (
	clientMetrics = ClientMetrics{
		totalSoraClientConnections: newDescWithLabel("client_type_total", "The total number of connections by Sora client types", []string{"client", "state"}),
	}
)

type ClientMetrics struct {
	totalSoraClientConnections *prometheus.Desc
}

func (m *ClientMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.totalSoraClientConnections
}

func (m *ClientMetrics) Collect(ch chan<- prometheus.Metric, report soraClientReport) {
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalFailedSoraClientType.SoraAndroidSdk), "android_sdk", "failed")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalFailedSoraClientType.SoraCSdk), "c_sdk", "failed")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalFailedSoraClientType.SoraCppSdk), "cpp_sdk", "failed")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalFailedSoraClientType.SoraFlutterSdk), "flutter_sdk", "failed")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalFailedSoraClientType.SoraIosSdk), "ios_sdk", "failed")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalFailedSoraClientType.SoraJsSdk), "js_sdk", "failed")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalFailedSoraClientType.SoraUnitySdk), "unity_sdk", "failed")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalFailedSoraClientType.ObsStudioWhip), "obs_studio_whip", "failed")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalFailedSoraClientType.ObsStudioWhep), "obs_studio_whep", "failed")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalFailedSoraClientType.SoraPythonSdk), "python_sdk", "failed")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalFailedSoraClientType.Unknown), "unknown", "failed")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalFailedSoraClientType.WebrtcLoadTestingToolZakuro), "zakuro", "failed")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalFailedSoraClientType.WebrtcNativeClientMomo), "momo", "failed")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalSuccessfulSoraClientType.SoraAndroidSdk), "android_sdk", "successful")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalSuccessfulSoraClientType.SoraCSdk), "c_sdk", "successful")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalSuccessfulSoraClientType.SoraCppSdk), "cpp_sdk", "successful")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalSuccessfulSoraClientType.SoraFlutterSdk), "flutter_sdk", "successful")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalSuccessfulSoraClientType.SoraIosSdk), "ios_sdk", "successful")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalSuccessfulSoraClientType.SoraJsSdk), "js_sdk", "successful")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalSuccessfulSoraClientType.SoraUnitySdk), "unity_sdk", "successful")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalSuccessfulSoraClientType.ObsStudioWhip), "obs_studio_whip", "successful")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalSuccessfulSoraClientType.ObsStudioWhep), "obs_studio_whep", "successful")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalSuccessfulSoraClientType.SoraPythonSdk), "python_sdk", "successful")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalSuccessfulSoraClientType.Unknown), "unknown", "successful")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalSuccessfulSoraClientType.WebrtcLoadTestingToolZakuro), "zakuro", "successful")
	ch <- newCounter(m.totalSoraClientConnections, float64(report.TotalSuccessfulSoraClientType.WebrtcNativeClientMomo), "momo", "successful")
}
