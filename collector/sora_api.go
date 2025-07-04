package collector

type soraGetStatsReport struct {
	SoraVersion string `json:"version"`
	soraConnectionReport
	soraWebhookReport
	soraSrtpReport
	soraSctpReport
	SoraClientReport          soraClientReport          `json:"sora_client,omitempty"`
	SoraConnectionErrorReport soraConnectionErrorReport `json:"error,omitempty"`
	ErlangVMReport            erlangVMReport            `json:"erlang_vm,omitempty"`
	ClusterReport             soraClusterReport         `json:"cluster,omitempty"`
	ClusterRelay              []soraClusterRelay        `json:"cluster_relay,omitempty"`
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
	TotalIgnoredSessionWebhook    int64 `json:"total_ignored_session_webhook"`
	TotalSuccessfulEventWebhook   int64 `json:"total_successful_event_webhook"`
	TotalFailedEventWebhook       int64 `json:"total_failed_event_webhook"`
	TotalIgnoredEventWebhook      int64 `json:"total_ignored_event_webhook"`
	TotalSuccessfulStatsWebhook   int64 `json:"total_successful_stats_webhook"`
	TotalFailedStatsWebhook       int64 `json:"total_failed_stats_webhook"`
	TotalIgnoredStatsWebhook      int64 `json:"total_ignored_stats_webhook"`
}

type soraSrtpReport struct {
	TotalReceivedSrtp          int64 `json:"total_received_srtp"`
	TotalReceivedSrtpByteSize  int64 `json:"total_received_srtp_byte_size"`
	TotalSentSrtp              int64 `json:"total_sent_srtp"`
	TotalSentSrtpByteSize      int64 `json:"total_sent_srtp_byte_size"`
	TotalSentSrtpSfuDelayUs    int64 `json:"total_sent_srtp_sfu_delay_us"`
	TotalDecryptedSrtp         int64 `json:"total_decrypted_srtp"`
	TotalDecryptedSrtpByteSize int64 `json:"total_decrypted_srtp_byte_size"`
}

type soraSctpReport struct {
	TotalReceivedSctp         int64 `json:"total_received_sctp"`
	TotalReceivedSctpByteSize int64 `json:"total_received_sctp_byte_size"`
	TotalSentSctp             int64 `json:"total_sent_sctp"`
	TotalSentSctpByteSize     int64 `json:"total_sent_sctp_byte_size"`
}

type soraClientStatistics struct {
	SoraAndroidSdk              int64 `json:"sora_android_sdk"`
	SoraCSdk                    int64 `json:"sora_c_sdk"`
	SoraCppSdk                  int64 `json:"sora_cpp_sdk"`
	SoraFlutterSdk              int64 `json:"sora_flutter_sdk"`
	SoraIosSdk                  int64 `json:"sora_ios_sdk"`
	SoraJsSdk                   int64 `json:"sora_js_sdk"`
	SoraUnitySdk                int64 `json:"sora_unity_sdk"`
	ObsStudioWhip               int64 `json:"obs_studio_whip"`
	ObsStudioWhep               int64 `json:"obs_studio_whep"`
	SoraPythonSdk               int64 `json:"sora_python_sdk"`
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

type soraClusterReport struct {
	RaftState       string `json:"raft_state"`
	RaftTerm        int64  `json:"raft_term"`
	RaftCommitIndex int64  `json:"raft_commit_index"`
}

type soraClusterNode struct {
	ClusterNodeName string `json:"cluster_node_name"`
	NodeName        string `json:"node_name"`
	Mode            string `json:"mode"`
	Connected       bool   `json:"connected"`
	TemporaryNode   bool   `json:"temporary_node"`
}

type soraClusterRelay struct {
	NodeName              string                   `json:"node_name"`
	TotalReceivedByteSize int64                    `json:"total_received_byte_size"`
	TotalSentByteSize     int64                    `json:"total_sent_byte_size"`
	TotalReceived         int64                    `json:"total_received"`
	TotalSent             int64                    `json:"total_sent"`
	Plumtree              soraClusterRelayPlumtree `json:"plumtree"`
}

type soraClusterRelayPlumtree struct {
        TotalSentGossip        int64 `json:"total_sent_gossip"`
        TotalReceivedGossip    int64 `json:"total_received_gossip"`
        TotalReceivedGossipHop int64 `json:"total_received_gossip_hop"`
        TotalSentIhave         int64 `json:"total_sent_ihave"`
        TotalReceivedIhave     int64 `json:"total_received_ihave"`
        TotalSentGraft         int64 `json:"total_sent_graft"`
        TotalReceivedGraft     int64 `json:"total_received_graft"`
        TotalSentPrune         int64 `json:"total_sent_prune"`
        TotalReceivedPrune     int64 `json:"total_received_prune"`
        TotalGraftMiss         int64 `json:"total_graft_miss"`
        TotalSkippedSend       int64 `json:"total_skipped_send"`
        TotalIhaveOverflow     int64 `json:"total_ihave_overflow"`
        TotalIgnored           int64 `json:"total_ignored"`
}

type soraLicenseInfo struct {
	ExpiredAt      string `json:"expired_at"`
	MaxConnections int64  `json:"max_connections"`
	MaxNodes       *int64 `json:"max_nodes"`
	ProductName    string `json:"product_name"`
	SerialCode     string `json:"serial_code"`
	Type           string `json:"type"`
}
