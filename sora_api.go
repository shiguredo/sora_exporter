package main

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
