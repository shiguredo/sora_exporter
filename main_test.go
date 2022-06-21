package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/shiguredo/sora_exporter/collector"
)

var (
	testJSONData = `{
		"average_duration_sec": 706,
		"average_setup_time_msec": 372,
		"erlang_vm": {
		  "memory": {
			"total": 1234,
			"processes": 5678,
			"processes_used": 6666,
			"system": 4444,
			"atom": 3333,
			"atom_used": 2222,
			"binary": 1111,
			"code": 777,
			"ets": 888
		  },
		  "statistics": {
			"active_tasks": [
			  0,
			  1,
			  0
			],
			"active_tasks_all": [
			  0,
			  1,
			  0,
			  0
			],
			"context_switches": 321,
			"exact_reductions": {
			  "exact_reductions_since_last_call": 213,
			  "total_exact_reductions": 432
			},
			"garbage_collection": {
			  "number_of_gcs": 912,
			  "words_reclaimed": 345
			},
			"io": {
			  "input": 769,
			  "output": 331
			},
			"reductions": {
			  "reductions_since_last_call": 443,
			  "total_reductions": 8080
			},
			"run_queue": 0,
			"run_queue_lengths": [
			  0,
			  0,
			  0
			],
			"run_queue_lengths_all": [
			  0,
			  0,
			  0,
			  0
			],
			"runtime": {
			  "time_since_last_call": 100,
			  "total_run_time": 200
			},
			"total_active_tasks": 1,
			"total_active_tasks_all": 1,
			"total_run_queue_lengths": 0,
			"total_run_queue_lengths_all": 0,
			"wall_clock": {
			  "total_wallclock_time": 994,
			  "wallclock_time_since_last_call": 9090
			}
		  }
		},
		"error": {
		  "sdp_generation_error": 3,
		  "signaling_error": 4
		},
		"sora_client": {
		  "total_failed_sora_client_type": {
			"sora_android_sdk": 1,
			"sora_ios_sdk": 2,
			"sora_js_sdk": 3,
			"sora_unity_sdk": 4,
			"unknown": 5,
			"webrtc_native_client_momo": 6
		  },
		  "total_successful_sora_client_type": {
			"sora_android_sdk": 11,
			"sora_ios_sdk": 22,
			"sora_js_sdk": 33,
			"sora_unity_sdk": 44,
			"unknown": 55,
			"webrtc_native_client_momo": 66
		  }
		},
		"total_connection_created": 2,
		"total_connection_destroyed": 2,
		"total_connection_updated": 23,
		"total_duration_sec": 1412,
		"total_failed_connections": 0,
		"total_ongoing_connections": 0,
		"total_received_invalid_turn_tcp_packet": 0,
		"total_session_created": 1,
		"total_session_destroyed": 0,
		"total_successful_connections": 2,
		"total_turn_tcp_connections": 2,
		"total_turn_udp_connections": 0,
		"version": "2022.1.0-canary.28"
	  }`
	listClusterNodesJSONData = `[
		{
		  "node_name": "node-01_canary_sora@10.211.55.42",
		  "epoch": 1,
		  "mode": "normal",
		  "cluster_signaling_url": "ws://127.0.0.1:5001/signaling",
		  "cluster_api_url": "http://127.0.0.1:3101/",
		  "member_since": "2022-05-09T07:44:52.973761Z",
		  "sora_version": "2022.1.0-canary.44",
		  "license_max_nodes": 10,
		  "license_max_connections": 100,
		  "license_serial_code": "SAMPLE-SRA-E001-202212-N10-100",
		  "license_type": "Experimental",
		  "connected": true
		},
		{
		  "node_name": "node-02_canary_sora@10.211.55.40",
		  "epoch": 1,
		  "mode": "block_new_connection",
		  "cluster_signaling_url": "ws://127.0.0.1:5002/signaling",
		  "cluster_api_url": "http://127.0.0.1:3102/",
		  "member_since": "2022-05-09T07:44:54.160763Z",
		  "sora_version": "2022.1.0-canary.44",
		  "license_max_nodes": 10,
		  "license_max_connections": 100,
		  "license_serial_code": "SAMPLE-SRA-E001-202212-N10-100",
		  "license_type": "Experimental",
		  "connected": true
		}
	  ]`
	listClusterNodesCurrentJSONData = `[
		{
		  "cluster_node_name": "node-01_canary_sora@10.211.55.42",
		  "epoch": 1,
		  "mode": "normal",
		  "member_since": "2022-05-02T15:26:44.302363Z",
		  "sora_version": "2021.2.9",
		  "license_max_connections": 100,
		  "license_serial_code": "SAMPLE-SRA-E001-202212-N10-100",
		  "license_type": "Experimental",
		  "cluster_signaling_url": "ws://127.0.0.1:5001/signaling",
		  "cluster_api_url": "http://10.1.1.4:3000/",
		  "connected": false
		},
		{
		  "cluster_node_name": "node-02_canary_sora@10.211.55.40",
		  "epoch": 1,
		  "mode": "block_new_connection",
		  "member_since": "2022-05-02T15:25:21.805078Z",
		  "sora_version": "2021.2.9",
		  "license_max_connections": 100,
		  "license_serial_code": "SAMPLE-SRA-E001-202212-N10-100",
		  "license_type": "Experimental",
		  "cluster_signaling_url": "ws://127.0.0.1:5002/signaling",
		  "cluster_api_url": "http://10.1.1.3:3000/",
		  "connected": true
		}
	  ]
	  `
)

type sora struct {
	*httptest.Server
	response                 []byte
	listClusterNodesResponse []byte
}

func newSora(response []byte, listClusterNodesResponse []byte) *sora {
	s := &sora{response: response, listClusterNodesResponse: listClusterNodesResponse}
	s.Server = httptest.NewServer(soraHandler(s))
	return s
}

func soraHandler(s *sora) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-sora-target") == "Sora_20211215.ListClusterNodes" {
			w.Write(s.listClusterNodesResponse)
			return
		}
		w.Write(s.response)
	}
}

func expectMetrics(t *testing.T, c prometheus.Collector, fixture string) {
	exp, err := os.Open(path.Join("test", fixture))
	if err != nil {
		t.Fatal(fmt.Errorf("The fixture file can't open %q: %w", fixture, err))
	}
	if err := testutil.CollectAndCompare(c, exp); err != nil {
		t.Fatal("Unexpect metrics returned:", err)
	}
}

func TestInvalidConfig(t *testing.T) {
	s := newSora([]byte("invalid config parameter"), []byte(listClusterNodesJSONData))
	defer s.Close()

	timeout, _ := time.ParseDuration("5s")
	h := collector.NewCollector(&collector.CollectorOptions{
		URI:                              s.URL,
		SkipSslVerify:                    true,
		Timeout:                          timeout,
		Logger:                           log.NewNopLogger(),
		EnableSoraClientMetrics:          true,
		EnableSoraConnectionErrorMetrics: true,
		EnableErlangVMMetrics:            true,
		EnableSoraClusterMetrics:         true,
	})
	expectMetrics(t, h, "invalid_config.metrics")
}

func TestMaximumMetrics(t *testing.T) {
	s := newSora([]byte(testJSONData), []byte(listClusterNodesJSONData))
	defer s.Close()

	timeout, _ := time.ParseDuration("5s")
	h := collector.NewCollector(&collector.CollectorOptions{
		URI:                              s.URL,
		SkipSslVerify:                    true,
		Timeout:                          timeout,
		Logger:                           log.NewNopLogger(),
		EnableSoraClientMetrics:          true,
		EnableSoraConnectionErrorMetrics: true,
		EnableErlangVMMetrics:            true,
		EnableSoraClusterMetrics:         true,
	})
	expectMetrics(t, h, "maximum.metrics")
}

func TestSoraErlangVMEnabledMetrics(t *testing.T) {
	s := newSora([]byte(testJSONData), []byte(listClusterNodesJSONData))
	defer s.Close()

	timeout, _ := time.ParseDuration("5s")
	h := collector.NewCollector(&collector.CollectorOptions{
		URI:                              s.URL,
		SkipSslVerify:                    true,
		Timeout:                          timeout,
		Logger:                           log.NewNopLogger(),
		EnableSoraClientMetrics:          false,
		EnableSoraConnectionErrorMetrics: false,
		EnableErlangVMMetrics:            true,
		EnableSoraClusterMetrics:         false,
	})
	expectMetrics(t, h, "sora_erlang_vm_enabled.metrics")
}

func TestSoraClientEnabledMetrics(t *testing.T) {
	s := newSora([]byte(testJSONData), []byte(listClusterNodesJSONData))
	defer s.Close()

	timeout, _ := time.ParseDuration("5s")
	h := collector.NewCollector(&collector.CollectorOptions{
		URI:                              s.URL,
		SkipSslVerify:                    true,
		Timeout:                          timeout,
		Logger:                           log.NewNopLogger(),
		EnableSoraClientMetrics:          true,
		EnableSoraConnectionErrorMetrics: false,
		EnableErlangVMMetrics:            false,
		EnableSoraClusterMetrics:         false,
	})
	expectMetrics(t, h, "sora_client_enabled.metrics")
}

func TestSoraConnectionErrorEnabledMetrics(t *testing.T) {
	s := newSora([]byte(testJSONData), []byte(listClusterNodesJSONData))
	defer s.Close()

	timeout, _ := time.ParseDuration("5s")
	h := collector.NewCollector(&collector.CollectorOptions{
		URI:                              s.URL,
		SkipSslVerify:                    true,
		Timeout:                          timeout,
		Logger:                           log.NewNopLogger(),
		EnableSoraClientMetrics:          false,
		EnableSoraConnectionErrorMetrics: true,
		EnableErlangVMMetrics:            false,
		EnableSoraClusterMetrics:         false,
	})
	expectMetrics(t, h, "sora_connection_error_enabled.metrics")
}

func TestMinimumMetrics(t *testing.T) {
	resp := `{
		"average_duration_sec": 706,
		"average_setup_time_msec": 12000,
		"total_connection_created": 3,
		"total_connection_destroyed": 2,
		"total_connection_updated": 23,
		"total_duration_sec": 1412,
		"total_failed_connections": 100,
		"total_ongoing_connections": 88,
		"total_received_invalid_turn_tcp_packet": 123,
		"total_session_created": 111,
		"total_session_destroyed": 222,
		"total_successful_connections": 333,
		"total_turn_tcp_connections": 444,
		"total_turn_udp_connections": 555,
		"version": "2022.1.0-canary.28"
	  }`
	s := newSora([]byte(resp), []byte(listClusterNodesJSONData))
	defer s.Close()

	timeout, _ := time.ParseDuration("5s")
	h := collector.NewCollector(&collector.CollectorOptions{
		URI:                              s.URL,
		SkipSslVerify:                    true,
		Timeout:                          timeout,
		Logger:                           log.NewNopLogger(),
		EnableSoraClientMetrics:          false,
		EnableSoraConnectionErrorMetrics: false,
		EnableErlangVMMetrics:            false,
		EnableSoraClusterMetrics:         false,
	})
	expectMetrics(t, h, "minimum.metrics")
}

func TestSoraClusterEnabledMetrics(t *testing.T) {
	s := newSora([]byte(testJSONData), []byte(listClusterNodesJSONData))
	defer s.Close()

	timeout, _ := time.ParseDuration("5s")
	h := collector.NewCollector(&collector.CollectorOptions{
		URI:                              s.URL,
		SkipSslVerify:                    true,
		Timeout:                          timeout,
		Logger:                           log.NewNopLogger(),
		EnableSoraClientMetrics:          false,
		EnableSoraConnectionErrorMetrics: false,
		EnableErlangVMMetrics:            false,
		EnableSoraClusterMetrics:         true,
	})
	expectMetrics(t, h, "sora_cluster_metrics_enabled.metrics")
}

// Sora-2021.9.x 系の JSON レスポンスデータでのテスト
func TestSoraClusterEnabledMetricsCurrentJsonData(t *testing.T) {
	s := newSora([]byte(testJSONData), []byte(listClusterNodesCurrentJSONData))
	defer s.Close()

	timeout, _ := time.ParseDuration("5s")
	h := collector.NewCollector(&collector.CollectorOptions{
		URI:                              s.URL,
		SkipSslVerify:                    true,
		Timeout:                          timeout,
		Logger:                           log.NewNopLogger(),
		EnableSoraClientMetrics:          false,
		EnableSoraConnectionErrorMetrics: false,
		EnableErlangVMMetrics:            false,
		EnableSoraClusterMetrics:         true,
	})
	expectMetrics(t, h, "sora_cluster_metrics_enabled.metrics")
}
