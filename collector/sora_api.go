package collector

type soraGetStatsReport struct {
	SoraVersion string `json:"version"`
	soraConnectionReport
	SoraClientReport soraClientReport `json:"sora_client,omitempty"`
	SoraErrorReport  soraErrorReport  `json:"error,omitempty"`
	ErlangVmReport   erlangVmReport   `json:"erlang_vm,omitempty"`
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

type soraClientStatistics struct {
	SoraAndroidSdk         int64 `json:"sora_android_sdk"`
	SoraIosSdk             int64 `json:"sora_ios_sdk"`
	SoraJsSdk              int64 `json:"sora_js_sdk"`
	SoraUnitySdk           int64 `json:"sora_unity_sdk"`
	Unknown                int64 `json:"unknown"`
	WebrtcNativeClientMomo int64 `json:"webrtc_native_client_momo"`
}

type soraClientReport struct {
	TotalFailedSoraClientType     soraClientStatistics `json:"total_failed_sora_client_type"`
	TotalSuccessfulSoraClientType soraClientStatistics `json:"total_successful_sora_client_type"`
}

type soraErrorReport struct {
	SdpGenerationError int64 `json:"sdp_generation_error"`
	SignalingError     int64 `json:"signaling_error"`
}

type erlangVmMemory struct {
	Total         int64 `json:"total"`
	Processes     int64 `json:"processes"`
	ProcessesUsed int64 `json:"processes_used"`
	System        int64 `json:"system"`
	Atom          int64 `json:"atom"`
	AtomUsed      int64 `json:"atom_used"`
	Binary        int64 `json:"binary"`
	Code          int64 `json:"code"`
	Ets           int64 `json:"ets"`
}

type exatReducations struct {
	ExactReductionsSinceLastCall int64 `json:"exact_reductions_since_last_call"`
	TotalExactReductions         int64 `json:"total_exact_reductions"`
}

type garbageCollection struct {
	NumberOfGcs    int64 `json:"number_of_gcs"`
	WordsReclaimed int64 `json:"words_reclaimed"`
}

type erlangIO struct {
	Input  int64 `json:"input"`
	Output int64 `json:"output"`
}

type reductions struct {
	ReductionsSinceLastCall int64 `json:"reductions_since_last_call"`
	TotalReductions         int64 `json:"total_reductions"`
}

type runtime struct {
	TimeSinceLastCall int64 `json:"time_since_last_call"`
	TotalRunTime      int64 `json:"total_run_time"`
}

type wallClock struct {
	TotalWallclockTime         int64 `json:"total_wallclock_time"`
	WallclockTimeSinceLastCall int64 `json:"wallclock_time_since_last_call"`
}

type erlangVmStatistics struct {
	ContextSwitches         int64             `json:"context_switches"`
	ExactReductions         exatReducations   `json:"exact_reductions"`
	GarbageCollection       garbageCollection `json:"garbage_collection"`
	Io                      erlangIO          `json:"io"`
	Reductions              reductions        `json:"reductions"`
	RunQueue                int64             `json:"run_queue"`
	Runtime                 runtime           `json:"runtime"`
	TotalActiveTasks        int64             `json:"total_active_tasks"`
	TotalActiveTasksAll     int64             `json:"total_active_tasks_all"`
	TotalRunQueueLengths    int64             `json:"total_run_queue_lengths"`
	TotalRunQueueLengthsAll int64             `json:"total_run_queue_lengths_all"`
	WallClock               wallClock         `json:"wall_clock"`
	// ErlangVmActiveTasks                                 []int64 `json:"active_tasks"`
	// ErlangVmActiveTasksAll                              []int64 `json:"active_tasks_all"`
	// ErlangVmRunQueueLengths                             []int64 `json:"run_queue_lengths"`
	// ErlangVmRunQueueLengthsAll                          []int64 `json:"run_queue_lengths_all"`
}

type erlangVmReport struct {
	ErlangVmMemory     erlangVmMemory     `json:"memory"`
	ErlangVmStatistics erlangVmStatistics `json:"statistics"`
}
