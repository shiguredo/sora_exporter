package main

import "github.com/prometheus/client_golang/prometheus"

var (
	clientMetrics = ClientMetrics{
		totalFailedSoraClientTypeSoraAndroidSdk:             newDesc("sora_client_type_sora_android_sdk_failed_total", "The total number of failed connections for Sora Android SDK."),
		totalFailedSoraClientTypeSoraIosSdk:                 newDesc("sora_client_type_sora_ios_sdk_failed_total", "The total number of failed connections for Sora IOS SDK."),
		totalFailedSoraClientTypeSoraJsSdk:                  newDesc("sora_client_type_sora_js_sdk_failed_total", "The total number of failed connections for Sora JavaScript SDK."),
		totalFailedSoraClientTypeSoraUnitySdk:               newDesc("sora_client_type_sora_unity_sdk_failed_total", "The total number of failed connections for Sora Unity SDK."),
		totalFailedSoraClientTypeUnknown:                    newDesc("sora_client_type_unknown_failed_total", "The total number of failed connections for WebRTC native client Momo."),
		totalFailedSoraClientTypeWebrtcNativeClientMomo:     newDesc("sora_client_type_webrtc_native_client_monmo_failed_total", "The total number of failed connections for WebRTC native client Momo."),
		totalSuccessfulSoraClientTypeSoraAndroidSdk:         newDesc("sora_client_type_sora_android_sdk_successful_total", "The total number of successful connections for Sora Android SDK."),
		totalSuccessfulSoraClientTypeSoraIosSdk:             newDesc("sora_client_type_sora_ios_sdk_successful_total", "The total number of successful connections for Sora IOS SDK."),
		totalSuccessfulSoraClientTypeSoraJsSdk:              newDesc("sora_client_type_sora_js_sdk_successful_total", "The total number of successful connections for Sora JavaScript SDK."),
		totalSuccessfulSoraClientTypeSoraUnitySdk:           newDesc("sora_client_type_sora_unity_sdk_successful_total", "The total number of successful connections for Sora Unity SDK."),
		totalSuccessfulSoraClientTypeUnknown:                newDesc("sora_client_type_known_successful_total", "The total number of successful connections for unknown client"),
		totalSuccessfulSoraClientTypeWebrtcNativeClientMomo: newDesc("sora_client_type_webrtc_native_client_momo_successful_total", "The total number of successful connections for WebRTC native client Momo."),
	}
)

type ClientMetrics struct {
	totalFailedSoraClientTypeSoraAndroidSdk             *prometheus.Desc
	totalFailedSoraClientTypeSoraIosSdk                 *prometheus.Desc
	totalFailedSoraClientTypeSoraJsSdk                  *prometheus.Desc
	totalFailedSoraClientTypeSoraUnitySdk               *prometheus.Desc
	totalFailedSoraClientTypeUnknown                    *prometheus.Desc
	totalFailedSoraClientTypeWebrtcNativeClientMomo     *prometheus.Desc
	totalSuccessfulSoraClientTypeSoraAndroidSdk         *prometheus.Desc
	totalSuccessfulSoraClientTypeSoraIosSdk             *prometheus.Desc
	totalSuccessfulSoraClientTypeSoraJsSdk              *prometheus.Desc
	totalSuccessfulSoraClientTypeSoraUnitySdk           *prometheus.Desc
	totalSuccessfulSoraClientTypeUnknown                *prometheus.Desc
	totalSuccessfulSoraClientTypeWebrtcNativeClientMomo *prometheus.Desc
}

func (m *ClientMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.totalFailedSoraClientTypeSoraAndroidSdk
	ch <- m.totalFailedSoraClientTypeSoraIosSdk
	ch <- m.totalFailedSoraClientTypeSoraJsSdk
	ch <- m.totalFailedSoraClientTypeSoraUnitySdk
	ch <- m.totalFailedSoraClientTypeUnknown
	ch <- m.totalFailedSoraClientTypeWebrtcNativeClientMomo
	ch <- m.totalSuccessfulSoraClientTypeSoraAndroidSdk
	ch <- m.totalSuccessfulSoraClientTypeSoraIosSdk
	ch <- m.totalSuccessfulSoraClientTypeSoraJsSdk
	ch <- m.totalSuccessfulSoraClientTypeSoraUnitySdk
	ch <- m.totalSuccessfulSoraClientTypeUnknown
	ch <- m.totalSuccessfulSoraClientTypeWebrtcNativeClientMomo
}

func (m *ClientMetrics) Collect(ch chan<- prometheus.Metric, report soraClientReport) {
	ch <- newCounter(m.totalFailedSoraClientTypeSoraAndroidSdk, float64(report.TotalFailedSoraClientType.SoraAndroidSdk))
	ch <- newCounter(m.totalFailedSoraClientTypeSoraIosSdk, float64(report.TotalFailedSoraClientType.SoraIosSdk))
	ch <- newCounter(m.totalFailedSoraClientTypeSoraJsSdk, float64(report.TotalFailedSoraClientType.SoraJsSdk))
	ch <- newCounter(m.totalFailedSoraClientTypeSoraUnitySdk, float64(report.TotalFailedSoraClientType.SoraUnitySdk))
	ch <- newCounter(m.totalFailedSoraClientTypeUnknown, float64(report.TotalFailedSoraClientType.Unknown))
	ch <- newCounter(m.totalFailedSoraClientTypeWebrtcNativeClientMomo, float64(report.TotalFailedSoraClientType.WebrtcNativeClientMomo))
	ch <- newCounter(m.totalSuccessfulSoraClientTypeSoraAndroidSdk, float64(report.TotalSuccessfulSoraClientType.SoraAndroidSdk))
	ch <- newCounter(m.totalSuccessfulSoraClientTypeSoraIosSdk, float64(report.TotalSuccessfulSoraClientType.SoraIosSdk))
	ch <- newCounter(m.totalSuccessfulSoraClientTypeSoraJsSdk, float64(report.TotalSuccessfulSoraClientType.SoraJsSdk))
	ch <- newCounter(m.totalSuccessfulSoraClientTypeSoraUnitySdk, float64(report.TotalSuccessfulSoraClientType.SoraUnitySdk))
	ch <- newCounter(m.totalSuccessfulSoraClientTypeUnknown, float64(report.TotalSuccessfulSoraClientType.Unknown))
	ch <- newCounter(m.totalSuccessfulSoraClientTypeWebrtcNativeClientMomo, float64(report.TotalSuccessfulSoraClientType.WebrtcNativeClientMomo))
}
