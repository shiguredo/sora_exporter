package collector

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	mutex                   sync.RWMutex
	httpClient              HTTPClient
	logger                  log.Logger
	timeout                 time.Duration
	URI                     string
	enableSoraClientMetrics bool
	enableSoraErrorMetrics  bool
	enableErlangVmMetrics   bool

	soraVersionInfo *prometheus.Desc

	// Sora Connection stats
	totalConnectionCreated     *prometheus.Desc
	totalConnectionUpdated     *prometheus.Desc
	totalConnectionDestroyed   *prometheus.Desc
	totalSuccessfulConnections *prometheus.Desc
	totalOngoingConnections    *prometheus.Desc
	totalFailedConnections     *prometheus.Desc
	totalDurationSec           *prometheus.Desc
	totalTurnUdpConnections    *prometheus.Desc
	totalTurnTcpConnections    *prometheus.Desc
	averageDurationSec         *prometheus.Desc
	averageSetupTimeSec        *prometheus.Desc

	// Sora Client stats
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

	// Sora Connect error stats
	sdpGenerationError *prometheus.Desc
	signalingError     *prometheus.Desc

	// Erlang VM stats
	// XXX(tnamao): Sora GetStatsReport API が array を返す値は未対応です。
	erlangVmMemoryTotal                                 *prometheus.Desc
	erlangVmMemoryProcesses                             *prometheus.Desc
	erlangVmMemoryProcessesUsed                         *prometheus.Desc
	erlangVmMemorySystem                                *prometheus.Desc
	erlangVmMemoryAtom                                  *prometheus.Desc
	erlangVmMemoryAtomUsed                              *prometheus.Desc
	erlangVmMemoryBinary                                *prometheus.Desc
	erlangVmMemoryCode                                  *prometheus.Desc
	erlangVmMemoryEts                                   *prometheus.Desc
	erlangVmContextSwitches                             *prometheus.Desc
	erlangVmExactReductionsExactReductionsSinceLastCall *prometheus.Desc
	erlangVmExactReductionsTotalExactReductions         *prometheus.Desc
	erlangVmGarbageCollectionNumberOfGcs                *prometheus.Desc
	erlangVmGarbageCollectionWordsReclaimed             *prometheus.Desc
	erlangVmIoInput                                     *prometheus.Desc
	erlangVmIoOutput                                    *prometheus.Desc
	erlangVmReductionsReductionsSinceLastCall           *prometheus.Desc
	erlangVmReductionsTotalReductions                   *prometheus.Desc
	erlangVmRunQueue                                    *prometheus.Desc
	erlangVmRuntimeTimeSinceLastCall                    *prometheus.Desc
	erlangVmRuntimeTotalRunTime                         *prometheus.Desc
	erlangVmTotalActiveTasks                            *prometheus.Desc
	erlangVmTotalActiveTasksAll                         *prometheus.Desc
	erlangVmTotalRunQueueLengths                        *prometheus.Desc
	erlangVmTotalRunQueueLengthsAll                     *prometheus.Desc
	erlangVmWallClockTotalWallclockTime                 *prometheus.Desc
	erlangVmWallClockWallclockTimeSinceLastCall         *prometheus.Desc
	// erlangVmActiveTasks                                 *prometheus.Desc
	// erlangVmActiveTasksAll                              *prometheus.Desc
	// erlangVmRunQueueLengths                             *prometheus.Desc
	// erlangVmRunQueueLengthsAll                          *prometheus.Desc
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

func New(uri string, timeout time.Duration, logger log.Logger, enableSoraClientMetrics bool, enableSoraErrorMetrics bool, enableErlangVmMetrics bool) *Collector {
	return &Collector{
		httpClient: http.DefaultClient,
		URI:        uri,
		timeout:    timeout,
		logger:     logger,

		enableSoraClientMetrics: enableSoraClientMetrics,
		enableSoraErrorMetrics:  enableSoraErrorMetrics,
		enableErlangVmMetrics:   enableErlangVmMetrics,

		soraVersionInfo: newDescWithLabel("sora_version_info", "sora version info.", []string{"version"}),

		totalConnectionCreated:     newDesc("connections_created_total", "The total number of connections created."),
		totalConnectionUpdated:     newDesc("connections_updated_total", "The total number of connections updated."),
		totalConnectionDestroyed:   newDesc("connections_destroyed_total", "The total number of connections destryed."),
		totalSuccessfulConnections: newDesc("successfull_connections_total", "The total number of successfull connections."),
		totalOngoingConnections:    newDesc("ongoing_connections_total", "The total number of ongoing connections."),
		totalFailedConnections:     newDesc("failed_connections_total", "The total number of failed connections."),
		totalDurationSec:           newDesc("duration_seconds_total", "The total duration of connections."),
		totalTurnUdpConnections:    newDesc("turn_udp_connections_total", "The total number of connections with TURN-UDP."),
		totalTurnTcpConnections:    newDesc("turn_tcp_connections_total", "The total number of connections with TURN-TCP."),
		averageDurationSec:         newDesc("average_duration_seconds", "The average connection duration in seconds."),
		averageSetupTimeSec:        newDesc("average_setup_time_seconds", "The average setup time in seconds."),

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

		sdpGenerationError: newDesc("sdp_generation_error_total", "The total number of SDP genration error."),
		signalingError:     newDesc("signaling_error_total", "The total number of signaling error."),

		erlangVmMemoryTotal:                                 newDesc("erlang_vm_memory_total", "The total amount of memory currently allocated. This is the same as the sum of the memory size for processes and system."),
		erlangVmMemoryProcesses:                             newDesc("erlang_vm_memory_processes", "The total amount of memory currently allocated for the Erlang processes."),
		erlangVmMemoryProcessesUsed:                         newDesc("erlang_vm_memory_processes_used", "The total amount of memory currently used by the Erlang processes. This is part of the memory presented as processes memory."),
		erlangVmMemorySystem:                                newDesc("erlang_vm_memory_system", "The total amount of memory currently allocated for the emulator that is not directly related to any Erlang process. Memory presented as processes is not included in this memory. instrument(3) can be used to get a more detailed breakdown of what memory is part of this type."),
		erlangVmMemoryAtom:                                  newDesc("erlang_vm_memory_atom", "The total amount of memory currently allocated for atoms. This memory is part of the memory presented as system memory."),
		erlangVmMemoryAtomUsed:                              newDesc("erlang_vm_memory_atom_used", "The total amount of memory currently used for atoms. This memory is part of the memory presented as atom memory."),
		erlangVmMemoryBinary:                                newDesc("erlang_vm_memory_binary", "The total amount of memory currently allocated for binaries. This memory is part of the memory presented as system memory."),
		erlangVmMemoryCode:                                  newDesc("erlang_vm_memory_code", "The total amount of memory currently allocated for Erlang code. This memory is part of the memory presented as system memory."),
		erlangVmMemoryEts:                                   newDesc("erlang_vm_memory_ets", "The total amount of memory currently allocated for ETS tables. This memory is part of the memory presented as system memory."),
		erlangVmContextSwitches:                             newDesc("erlang_vm_context_switches", "The total number of context switches since the system started."),
		erlangVmExactReductionsExactReductionsSinceLastCall: newDesc("erlang_vm_exact_reductions_exact_reductions_since_last_call", "The number of exact reductions since last call."),
		erlangVmExactReductionsTotalExactReductions:         newDesc("erlang_vm_exact_reductions_total_exact_reductions", "The total number of exact reductions."),
		erlangVmGarbageCollectionNumberOfGcs:                newDesc("erlang_vm_garbage_collection_number_of_gcs", "The number of information about garbage collection."),
		erlangVmGarbageCollectionWordsReclaimed:             newDesc("erlang_vm_garbage_collection_words_reclaimed", "The number of information about garbage collection word reclaimed."),
		erlangVmIoInput:                                     newDesc("erlang_vm_io_input", ""),
		erlangVmIoOutput:                                    newDesc("erlang_vm_io_output", ""),
		erlangVmReductionsReductionsSinceLastCall:           newDesc("erlang_vm_reductions_reductions_since_last_call", "The number of information about reductions."),
		erlangVmReductionsTotalReductions:                   newDesc("erlang_vm_reductions_total_reductions", "The total number of information about reductions."),
		erlangVmRunQueue:                                    newDesc("erlang_vm_run_queue", "The total length of all normal and dirty CPU run queues."),
		erlangVmRuntimeTimeSinceLastCall:                    newDesc("erlang_vm_runtime_time_since_last_call", "The number of information about runtime since last call, in milliseconds."),
		erlangVmRuntimeTotalRunTime:                         newDesc("erlang_vm_runtime_total_run_time", "The number of information about runtime, in milliseconds."),
		erlangVmTotalActiveTasks:                            newDesc("erlang_vm_total_active_tasks", "The number of active processes and ports on each run queue and its associated schedulers. (only tasks that are expected to be CPU bound are part of the result.)"),
		erlangVmTotalActiveTasksAll:                         newDesc("erlang_vm_total_active_tasks_all", "The number of active processes and ports on each run queue and its associated schedulers."),
		erlangVmTotalRunQueueLengths:                        newDesc("erlang_vm_total_run_queue_lengths", "The number of processes and ports ready to run for each run queue. (only run queues with work that is expected to be CPU bound is part of the result.)"),
		erlangVmTotalRunQueueLengthsAll:                     newDesc("erlang_vm_total_run_queue_lengths_all", "The number of processes and ports ready to run for each run queue."),
		erlangVmWallClockTotalWallclockTime:                 newDesc("erlang_vm_wall_clock_total_wallclock_time", "The number of information about wall clock."),
		erlangVmWallClockWallclockTimeSinceLastCall:         newDesc("erlang_vm_wall_clock_wallclock_time_since_last_call", "The number of information about wall clock since last call."),
		// erlangVmActiveTasks:                                 newDesc("erlang_vm_active_tasks", ""),
		// erlangVmActiveTasksAll:                              newDesc("erlang_vm_active_tasks_all", ""),
		// erlangVmRunQueueLengths:                             newDesc("erlang_vm_run_queue_lengths", ""),
		// erlangVmRunQueueLengthsAll:                          newDesc("erlang_vm_run_queue_lengths_all", ""),
	}
}

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

// TODO(tnamao): 不要？
var _ prometheus.Collector = (*Collector)(nil)

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URI, nil)
	if err != nil {
		level.Error(c.logger).Log("msg", "failed to create request to sora", "err", err)
		return
	}
	req.Header.Set("x-sora-target", "Sora_20171010.GetStatsReport")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		level.Error(c.logger).Log("msg", "failed to request to sora", "err", err)
		return
	}
	defer resp.Body.Close()

	var report soraGetStatsReport
	if err := json.NewDecoder(resp.Body).Decode(&report); err != nil {
		level.Error(c.logger).Log("msg", "failed to decode response body from sora", "err", err)
		return
	}

	ch <- newInfo(c.soraVersionInfo, report.SoraVersion)
	ch <- newCounter(c.totalConnectionCreated, float64(report.TotalConnectionCreated))
	ch <- newCounter(c.totalConnectionUpdated, float64(report.TotalConnectionUpdated))
	ch <- newCounter(c.totalConnectionDestroyed, float64(report.TotalConnectionDestroyed))
	ch <- newCounter(c.totalSuccessfulConnections, float64(report.TotalSuccessfulConnections))
	ch <- newGauge(c.totalOngoingConnections, float64(report.TotalOngoingConnections))
	ch <- newCounter(c.totalFailedConnections, float64(report.TotalFailedConnections))
	ch <- newCounter(c.totalDurationSec, float64(report.TotalDurationSec))
	ch <- newCounter(c.totalTurnUdpConnections, float64(report.TotalTurnUdpConnections))
	ch <- newCounter(c.totalTurnTcpConnections, float64(report.TotalTurnTcpConnections))
	ch <- newGauge(c.averageDurationSec, float64(report.AverageDurationSec))
	ch <- newGauge(c.averageSetupTimeSec, float64(report.AverageSetupTimeMsec/1000))

	if c.enableSoraClientMetrics {
		ch <- newCounter(c.totalFailedSoraClientTypeSoraAndroidSdk, float64(report.TotalFailedSoraClientTypeSoraAndroidSdk))
		ch <- newCounter(c.totalFailedSoraClientTypeSoraIosSdk, float64(report.TotalFailedSoraClientTypeSoraIosSdk))
		ch <- newCounter(c.totalFailedSoraClientTypeSoraJsSdk, float64(report.TotalFailedSoraClientTypeSoraJsSdk))
		ch <- newCounter(c.totalFailedSoraClientTypeSoraUnitySdk, float64(report.TotalFailedSoraClientTypeSoraUnitySdk))
		ch <- newCounter(c.totalFailedSoraClientTypeUnknown, float64(report.TotalFailedSoraClientTypeUnknown))
		ch <- newCounter(c.totalFailedSoraClientTypeWebrtcNativeClientMomo, float64(report.TotalFailedSoraClientTypeWebrtcNativeClientMomo))
		ch <- newCounter(c.totalSuccessfulSoraClientTypeSoraAndroidSdk, float64(report.TotalSuccessfulSoraClientTypeSoraAndroidSdk))
		ch <- newCounter(c.totalSuccessfulSoraClientTypeSoraIosSdk, float64(report.TotalSuccessfulSoraClientTypeSoraIosSdk))
		ch <- newCounter(c.totalSuccessfulSoraClientTypeSoraJsSdk, float64(report.TotalSuccessfulSoraClientTypeSoraJsSdk))
		ch <- newCounter(c.totalSuccessfulSoraClientTypeSoraUnitySdk, float64(report.TotalSuccessfulSoraClientTypeSoraUnitySdk))
		ch <- newCounter(c.totalSuccessfulSoraClientTypeUnknown, float64(report.TotalSuccessfulSoraClientTypeUnknown))
		ch <- newCounter(c.totalSuccessfulSoraClientTypeWebrtcNativeClientMomo, float64(report.TotalSuccessfulSoraClientTypeWebrtcNativeClientMomo))
	}

	if c.enableSoraErrorMetrics {
		ch <- newCounter(c.sdpGenerationError, float64(report.SdpGenerationError))
		ch <- newCounter(c.signalingError, float64(report.SignalingError))
	}

	if c.enableErlangVmMetrics {
		ch <- newGauge(c.erlangVmMemoryTotal, float64(report.ErlangVmMemoryTotal))
		ch <- newGauge(c.erlangVmMemoryProcesses, float64(report.ErlangVmMemoryProcesses))
		ch <- newGauge(c.erlangVmMemoryProcessesUsed, float64(report.ErlangVmMemoryProcessesUsed))
		ch <- newGauge(c.erlangVmMemorySystem, float64(report.ErlangVmMemorySystem))
		ch <- newGauge(c.erlangVmMemoryAtom, float64(report.ErlangVmMemoryAtom))
		ch <- newGauge(c.erlangVmMemoryAtomUsed, float64(report.ErlangVmMemoryAtomUsed))
		ch <- newGauge(c.erlangVmMemoryBinary, float64(report.ErlangVmMemoryBinary))
		ch <- newGauge(c.erlangVmMemoryCode, float64(report.ErlangVmMemoryCode))
		ch <- newGauge(c.erlangVmMemoryEts, float64(report.ErlangVmMemoryEts))
		ch <- newGauge(c.erlangVmContextSwitches, float64(report.ErlangVmContextSwitches))
		ch <- newGauge(c.erlangVmExactReductionsExactReductionsSinceLastCall, float64(report.ErlangVmExactReductionsExactReductionsSinceLastCall))
		ch <- newGauge(c.erlangVmExactReductionsTotalExactReductions, float64(report.ErlangVmExactReductionsTotalExactReductions))
		ch <- newGauge(c.erlangVmGarbageCollectionNumberOfGcs, float64(report.ErlangVmGarbageCollectionNumberOfGcs))
		ch <- newGauge(c.erlangVmGarbageCollectionWordsReclaimed, float64(report.ErlangVmGarbageCollectionWordsReclaimed))
		ch <- newGauge(c.erlangVmIoInput, float64(report.ErlangVmIoInput))
		ch <- newGauge(c.erlangVmIoOutput, float64(report.ErlangVmIoOutput))
		ch <- newGauge(c.erlangVmReductionsReductionsSinceLastCall, float64(report.ErlangVmReductionsReductionsSinceLastCall))
		ch <- newGauge(c.erlangVmReductionsTotalReductions, float64(report.ErlangVmReductionsTotalReductions))
		ch <- newGauge(c.erlangVmRunQueue, float64(report.ErlangVmRunQueue))
		ch <- newGauge(c.erlangVmRuntimeTimeSinceLastCall, float64(report.ErlangVmRuntimeTimeSinceLastCall))
		ch <- newGauge(c.erlangVmRuntimeTotalRunTime, float64(report.ErlangVmRuntimeTotalRunTime))
		ch <- newGauge(c.erlangVmTotalActiveTasks, float64(report.ErlangVmTotalActiveTasks))
		ch <- newGauge(c.erlangVmTotalActiveTasksAll, float64(report.ErlangVmTotalActiveTasksAll))
		ch <- newGauge(c.erlangVmTotalRunQueueLengths, float64(report.ErlangVmTotalRunQueueLengths))
		ch <- newGauge(c.erlangVmTotalRunQueueLengthsAll, float64(report.ErlangVmTotalRunQueueLengthsAll))
		ch <- newGauge(c.erlangVmWallClockTotalWallclockTime, float64(report.ErlangVmWallClockTotalWallclockTime))
		ch <- newGauge(c.erlangVmWallClockWallclockTimeSinceLastCall, float64(report.ErlangVmWallClockWallclockTimeSinceLastCall))
		// ch <- newGauge(c.erlangVmActiveTasks, float64(report.ErlangVmActiveTasks))
		// ch <- newGauge(c.erlangVmActiveTasksAll, float64(report.ErlangVmActiveTasksAll))
		// ch <- newGauge(c.erlangVmRunQueueLengths, float64(report.ErlangVmRunQueueLengths))
		// ch <- newGauge(c.erlangVmRunQueueLengthsAll, float64(report.ErlangVmRunQueueLengthsAll))
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.soraVersionInfo
	ch <- c.totalConnectionCreated
	ch <- c.totalConnectionUpdated
	ch <- c.totalConnectionUpdated
	ch <- c.totalSuccessfulConnections
	ch <- c.totalOngoingConnections
	ch <- c.totalFailedConnections
	ch <- c.totalDurationSec
	ch <- c.totalTurnUdpConnections
	ch <- c.totalTurnTcpConnections
	ch <- c.averageDurationSec
	ch <- c.averageSetupTimeSec

	if c.enableSoraClientMetrics {
		ch <- c.totalFailedSoraClientTypeSoraAndroidSdk
		ch <- c.totalFailedSoraClientTypeSoraIosSdk
		ch <- c.totalFailedSoraClientTypeSoraJsSdk
		ch <- c.totalFailedSoraClientTypeSoraUnitySdk
		ch <- c.totalFailedSoraClientTypeUnknown
		ch <- c.totalFailedSoraClientTypeWebrtcNativeClientMomo
		ch <- c.totalSuccessfulSoraClientTypeSoraAndroidSdk
		ch <- c.totalSuccessfulSoraClientTypeSoraIosSdk
		ch <- c.totalSuccessfulSoraClientTypeSoraJsSdk
		ch <- c.totalSuccessfulSoraClientTypeSoraUnitySdk
		ch <- c.totalSuccessfulSoraClientTypeUnknown
		ch <- c.totalSuccessfulSoraClientTypeWebrtcNativeClientMomo
	}

	if c.enableSoraErrorMetrics {
		ch <- c.sdpGenerationError
		ch <- c.signalingError
	}

	if c.enableErlangVmMetrics {
		ch <- c.erlangVmMemoryTotal
		ch <- c.erlangVmMemoryProcesses
		ch <- c.erlangVmMemoryProcessesUsed
		ch <- c.erlangVmMemorySystem
		ch <- c.erlangVmMemoryAtom
		ch <- c.erlangVmMemoryAtomUsed
		ch <- c.erlangVmMemoryBinary
		ch <- c.erlangVmMemoryCode
		ch <- c.erlangVmMemoryEts
		ch <- c.erlangVmContextSwitches
		ch <- c.erlangVmExactReductionsExactReductionsSinceLastCall
		ch <- c.erlangVmExactReductionsTotalExactReductions
		ch <- c.erlangVmGarbageCollectionNumberOfGcs
		ch <- c.erlangVmGarbageCollectionWordsReclaimed
		ch <- c.erlangVmIoInput
		ch <- c.erlangVmIoOutput
		ch <- c.erlangVmReductionsReductionsSinceLastCall
		ch <- c.erlangVmReductionsTotalReductions
		ch <- c.erlangVmRunQueue
		ch <- c.erlangVmRuntimeTimeSinceLastCall
		ch <- c.erlangVmRuntimeTotalRunTime
		ch <- c.erlangVmTotalActiveTasks
		ch <- c.erlangVmTotalActiveTasksAll
		ch <- c.erlangVmTotalRunQueueLengths
		ch <- c.erlangVmTotalRunQueueLengthsAll
		ch <- c.erlangVmWallClockTotalWallclockTime
		ch <- c.erlangVmWallClockWallclockTimeSinceLastCall
		// ch <- c.erlangVmActiveTasks
		// ch <- c.erlangVmActiveTasksAll
		// ch <- c.erlangVmRunQueueLengths
		// ch <- c.erlangVmRunQueueLengthsAll
	}
}

type soraGetStatsReport struct {
	SoraVersion string `json:"version"`
	soraConnectionReport
	soraClientReport
	soraErrorReport
	erlangVmReport
}

type soraConnectionReport struct {
	TotalConnectionCreated     int64 `json:"total_connection_created"`
	TotalConnectionUpdated     int64 `json:"total_connection_updated"`
	TotalConnectionDestroyed   int64 `json:"total_connection_destroyed"`
	TotalSuccessfulConnections int64 `json:"total_successful_connections"`
	TotalOngoingConnections    int64 `json:"total_ongoing_connections"`
	TotalFailedConnections     int64 `json:"total_failed_connections"`
	TotalDurationSec           int64 `json:"total_duration_sec"`
	TotalTurnUdpConnections    int64 `json:"total_turn_udp_connections"`
	TotalTurnTcpConnections    int64 `json:"total_turn_tcp_connections"`
	AverageDurationSec         int64 `json:"average_duration_sec"`
	AverageSetupTimeMsec       int64 `json:"average_setup_time_msec"`
}

type soraClientReport struct {
	TotalFailedSoraClientTypeSoraAndroidSdk             int64 `json:"sora_client.total_failed_sora_client_type.sora_android_sdk"`
	TotalFailedSoraClientTypeSoraIosSdk                 int64 `json:"sora_client.total_failed_sora_client_type.sora_ios_sdk"`
	TotalFailedSoraClientTypeSoraJsSdk                  int64 `json:"sora_client.total_failed_sora_client_type.sora_js_sdk"`
	TotalFailedSoraClientTypeSoraUnitySdk               int64 `json:"sora_client.total_failed_sora_client_type.sora_unity_sdk"`
	TotalFailedSoraClientTypeUnknown                    int64 `json:"sora_client.total_failed_sora_client_type.unknown"`
	TotalFailedSoraClientTypeWebrtcNativeClientMomo     int64 `json:"sora_client.total_failed_sora_client_type.webrtc_native_client_momo"`
	TotalSuccessfulSoraClientTypeSoraAndroidSdk         int64 `json:"sora_client.total_successful_sora_client_type.sora_android_sdk"`
	TotalSuccessfulSoraClientTypeSoraIosSdk             int64 `json:"sora_client.total_successful_sora_client_type.sora_ios_sdk"`
	TotalSuccessfulSoraClientTypeSoraJsSdk              int64 `json:"sora_client.total_successful_sora_client_type.sora_js_sdk"`
	TotalSuccessfulSoraClientTypeSoraUnitySdk           int64 `json:"sora_client.total_successful_sora_client_type.sora_unity_sdk"`
	TotalSuccessfulSoraClientTypeUnknown                int64 `json:"sora_client.total_successful_sora_client_type.unknown"`
	TotalSuccessfulSoraClientTypeWebrtcNativeClientMomo int64 `json:"sora_client.total_successful_sora_client_type.webrtc_native_client_momo"`
}

type soraErrorReport struct {
	SdpGenerationError int64 `json:"error.sdp_generation_error"`
	SignalingError     int64 `json:"error.signaling_error"`
}

type erlangVmReport struct {
	ErlangVmMemoryTotal                                 int64 `json:"erlang_vm.memory.total"`
	ErlangVmMemoryProcesses                             int64 `json:"erlang_vm.memory.processes"`
	ErlangVmMemoryProcessesUsed                         int64 `json:"erlang_vm.memory.processes_used"`
	ErlangVmMemorySystem                                int64 `json:"erlang_vm.memory.system"`
	ErlangVmMemoryAtom                                  int64 `json:"erlang_vm.memory.atom"`
	ErlangVmMemoryAtomUsed                              int64 `json:"erlang_vm.memory.atom_used"`
	ErlangVmMemoryBinary                                int64 `json:"erlang_vm.memory.binary"`
	ErlangVmMemoryCode                                  int64 `json:"erlang_vm.memory.code"`
	ErlangVmMemoryEts                                   int64 `json:"erlang_vm.memory.ets"`
	ErlangVmContextSwitches                             int64 `json:"erlang_vm.statistics.context_switches"`
	ErlangVmExactReductionsExactReductionsSinceLastCall int64 `json:"erlang_vm.statistics.exact_reductions.exact_reductions_since_last_call"`
	ErlangVmExactReductionsTotalExactReductions         int64 `json:"erlang_vm.statistics.exact_reductions.total_exact_reductions"`
	ErlangVmGarbageCollectionNumberOfGcs                int64 `json:"erlang_vm.statistics.garbage_collection.number_of_gcs"`
	ErlangVmGarbageCollectionWordsReclaimed             int64 `json:"erlang_vm.statistics.garbage_collection.words_reclaimed"`
	ErlangVmIoInput                                     int64 `json:"erlang_vm.statistics.io.input"`
	ErlangVmIoOutput                                    int64 `json:"erlang_vm.statistics.io.output"`
	ErlangVmReductionsReductionsSinceLastCall           int64 `json:"erlang_vm.statistics.reductions.reductions_since_last_call"`
	ErlangVmReductionsTotalReductions                   int64 `json:"erlang_vm.statistics.reductions.total_reductions"`
	ErlangVmRunQueue                                    int64 `json:"erlang_vm.statistics.run_queue"`
	ErlangVmRuntimeTimeSinceLastCall                    int64 `json:"erlang_vm.statistics.runtime.time_since_last_call"`
	ErlangVmRuntimeTotalRunTime                         int64 `json:"erlang_vm.statistics.runtime.total_run_time"`
	ErlangVmTotalActiveTasks                            int64 `json:"erlang_vm.statistics.total_active_tasks"`
	ErlangVmTotalActiveTasksAll                         int64 `json:"erlang_vm.statistics.total_active_tasks_all"`
	ErlangVmTotalRunQueueLengths                        int64 `json:"erlang_vm.statistics.total_run_queue_lengths"`
	ErlangVmTotalRunQueueLengthsAll                     int64 `json:"erlang_vm.statistics.total_run_queue_lengths_all"`
	ErlangVmWallClockTotalWallclockTime                 int64 `json:"erlang_vm.statistics.wall_clock.total_wallclock_time"`
	ErlangVmWallClockWallclockTimeSinceLastCall         int64 `json:"erlang_vm.statistics.wall_clock.wallclock_time_since_last_call"`
	// ErlangVmActiveTasks                                 []int64 `json:"erlang_vm.statistics.active_tasks"`
	// ErlangVmActiveTasksAll                              []int64 `json:"erlang_vm.statistics.active_tasks_all"`
	// ErlangVmRunQueueLengths                             []int64 `json:"erlang_vm.statistics.run_queue_lengths"`
	// ErlangVmRunQueueLengthsAll                          []int64 `json:"erlang_vm.statistics.run_queue_lengths_all"`
}
