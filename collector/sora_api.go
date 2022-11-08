package collector

type soraGetStatsReport struct {
	SoraVersion string `json:"version"`
	soraConnectionReport
	soraWebhookReport
	SoraClientReport          soraClientReport          `json:"sora_client,omitempty"`
	SoraConnectionErrorReport soraConnectionErrorReport `json:"error,omitempty"`
	ErlangVMReport            erlangVMReport            `json:"erlang_vm,omitempty"`
}

type soraConnectionReport struct {
	TotalConnectionCreated            int64 `json:"total_connection_created"`
	TotalConnectionUpdated            int64 `json:"total_connection_updated"`
	TotalConnectionDestroyed          int64 `json:"total_connection_destroyed"`
	TotalSuccessfulConnections        int64 `json:"total_successful_connections"`
	TotalOngoingConnections           int64 `json:"total_ongoing_connections"`
	TotalFailedConnections            int64 `json:"total_failed_connections"`
	TotalDurationSec                  int64 `json:"total_duration_sec"`
	TotalTurnUDPConnections           int64 `json:"total_turn_udp_connections"`
	TotalTurnTCPConnections           int64 `json:"total_turn_tcp_connections"`
	AverageDurationSec                int64 `json:"average_duration_sec"`
	AverageSetupTimeMsec              int64 `json:"average_setup_time_msec"`
	TotalSessionCreated               int64 `json:"total_session_created"`
	TotalSessionDestroyed             int64 `json:"total_session_destroyed"`
	TotalReceivedInvalidTurnTCPPacket int64 `json:"total_received_invalid_turn_tcp_packet"`
}

type soraWebhookReport struct {
	TotalAuthWebhookAllowed       int64 `json:"total_auth_webhook_allowed"`
	TotalAuthWebhookDenied        int64 `json:"total_auth_webhook_denied"`
	TotalSuccessfulAuthWebhook    int64 `json:"total_successful_auth_webhook"`
	TotalFailedAuthWebhook        int64 `json:"total_failed_auth_webhook"`
	TotalSuccessfulSessionWebhook int64 `json:"total_successful_session_webhook"`
	TotalFailedSessionWebhook     int64 `json:"total_failed_session_webhook"`
	TotalSuccessfulEventWebhook   int64 `json:"total_successful_event_webhook"`
	TotalFailedEventWebhook       int64 `json:"total_failed_event_webhook"`
}

type soraClientStatistics struct {
	SoraAndroidSdk              int64 `json:"sora_android_sdk"`
	SoraCppSdk                  int64 `json:"sora_cpp_sdk"`
	SoraFlutterSdk              int64 `json:"sora_flutter_sdk"`
	SoraIosSdk                  int64 `json:"sora_ios_sdk"`
	SoraJsSdk                   int64 `json:"sora_js_sdk"`
	SoraUnitySdk                int64 `json:"sora_unity_sdk"`
	Unknown                     int64 `json:"unknown"`
	WebrtcLoadTestingToolZakuro int64 `json:"webrtc_load_testing_tool_zakuro"`
	WebrtcNativeClientMomo      int64 `json:"webrtc_native_client_momo"`
}

type soraClientReport struct {
	TotalFailedSoraClientType     soraClientStatistics `json:"total_failed_sora_client_type"`
	TotalSuccessfulSoraClientType soraClientStatistics `json:"total_successful_sora_client_type"`
}

type soraConnectionErrorReport struct {
	SdpGenerationError int64 `json:"sdp_generation_error"`
	SignalingError     int64 `json:"signaling_error"`
}

type erlangVMMemory struct {
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

type erlangVMStatistics struct {
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
	// ErlangVMActiveTasks                                 []int64 `json:"active_tasks"`
	// ErlangVMActiveTasksAll                              []int64 `json:"active_tasks_all"`
	// ErlangVMRunQueueLengths                             []int64 `json:"run_queue_lengths"`
	// ErlangVMRunQueueLengthsAll                          []int64 `json:"run_queue_lengths_all"`
}

type erlangVMReport struct {
	ErlangVMMemory     erlangVMMemory     `json:"memory"`
	ErlangVMStatistics erlangVMStatistics `json:"statistics"`
}

type soraClusterNode struct {
	ClusterNodeName *string `json:"cluster_node_name"`
	NodeName        *string `json:"node_name"`
	Mode            *string `json:"mode"`
}
