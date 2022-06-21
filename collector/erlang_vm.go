package collector

import "github.com/prometheus/client_golang/prometheus"

// この統計情報はアンドキュメントです
var (
	erlangVMMetrics = ErlangVMMetrics{
		erlangVMMemoryTotal:                                 newDesc("erlang_vm_memory_total", "The total amount of memory currently allocated. This is the same as the sum of the memory size for processes and system."),
		erlangVMMemoryProcesses:                             newDesc("erlang_vm_memory_processes", "The total amount of memory currently allocated for the Erlang processes."),
		erlangVMMemoryProcessesUsed:                         newDesc("erlang_vm_memory_processes_used", "The total amount of memory currently used by the Erlang processes. This is part of the memory presented as processes memory."),
		erlangVMMemorySystem:                                newDesc("erlang_vm_memory_system", "The total amount of memory currently allocated for the emulator that is not directly related to any Erlang process. Memory presented as processes is not included in this memory. instrument(3) can be used to get a more detailed breakdown of what memory is part of this type."),
		erlangVMMemoryAtom:                                  newDesc("erlang_vm_memory_atom", "The total amount of memory currently allocated for atoms. This memory is part of the memory presented as system memory."),
		erlangVMMemoryAtomUsed:                              newDesc("erlang_vm_memory_atom_used", "The total amount of memory currently used for atoms. This memory is part of the memory presented as atom memory."),
		erlangVMMemoryBinary:                                newDesc("erlang_vm_memory_binary", "The total amount of memory currently allocated for binaries. This memory is part of the memory presented as system memory."),
		erlangVMMemoryCode:                                  newDesc("erlang_vm_memory_code", "The total amount of memory currently allocated for Erlang code. This memory is part of the memory presented as system memory."),
		erlangVMMemoryEts:                                   newDesc("erlang_vm_memory_ets", "The total amount of memory currently allocated for ETS tables. This memory is part of the memory presented as system memory."),
		erlangVMContextSwitches:                             newDesc("erlang_vm_context_switches", "The total number of context switches since the system started."),
		erlangVMExactReductionsExactReductionsSinceLastCall: newDesc("erlang_vm_exact_reductions_exact_reductions_since_last_call", "The number of exact reductions since last call."),
		erlangVMExactReductionsTotalExactReductions:         newDesc("erlang_vm_exact_reductions_total_exact_reductions", "The total number of exact reductions."),
		erlangVMGarbageCollectionNumberOfGcs:                newDesc("erlang_vm_garbage_collection_number_of_gcs", "The number of information about garbage collection."),
		erlangVMGarbageCollectionWordsReclaimed:             newDesc("erlang_vm_garbage_collection_words_reclaimed", "The number of information about garbage collection word reclaimed."),
		erlangVMIoInput:                                     newDesc("erlang_vm_io_input", "The total number of bytes received through ports."),
		erlangVMIoOutput:                                    newDesc("erlang_vm_io_output", "The total number of bytes output through ports."),
		erlangVMReductionsReductionsSinceLastCall:           newDesc("erlang_vm_reductions_reductions_since_last_call", "The number of information about reductions."),
		erlangVMReductionsTotalReductions:                   newDesc("erlang_vm_reductions_total_reductions", "The total number of information about reductions."),
		erlangVMRunQueue:                                    newDesc("erlang_vm_run_queue", "The total length of all normal and dirty CPU run queues."),
		erlangVMRuntimeTimeSinceLastCall:                    newDesc("erlang_vm_runtime_time_since_last_call", "The number of information about runtime since last call, in milliseconds."),
		erlangVMRuntimeTotalRunTime:                         newDesc("erlang_vm_runtime_total_run_time", "The number of information about runtime, in milliseconds."),
		erlangVMTotalActiveTasks:                            newDesc("erlang_vm_total_active_tasks", "The number of active processes and ports on each run queue and its associated schedulers. (only tasks that are expected to be CPU bound are part of the result.)"),
		erlangVMTotalActiveTasksAll:                         newDesc("erlang_vm_total_active_tasks_all", "The number of active processes and ports on each run queue and its associated schedulers."),
		erlangVMTotalRunQueueLengths:                        newDesc("erlang_vm_total_run_queue_lengths", "The number of processes and ports ready to run for each run queue. (only run queues with work that is expected to be CPU bound is part of the result.)"),
		erlangVMTotalRunQueueLengthsAll:                     newDesc("erlang_vm_total_run_queue_lengths_all", "The number of processes and ports ready to run for each run queue."),
		erlangVMWallClockTotalWallclockTime:                 newDesc("erlang_vm_wall_clock_total_wallclock_time", "The number of information about wall clock."),
		erlangVMWallClockWallclockTimeSinceLastCall:         newDesc("erlang_vm_wall_clock_wallclock_time_since_last_call", "The number of information about wall clock since last call."),
		// erlangVMActiveTasks:                                 newDesc("erlang_vm_active_tasks", ""),
		// erlangVMActiveTasksAll:                              newDesc("erlang_vm_active_tasks_all", ""),
		// erlangVMRunQueueLengths:                             newDesc("erlang_vm_run_queue_lengths", ""),
		// erlangVMRunQueueLengthsAll:                          newDesc("erlang_vm_run_queue_lengths_all", ""),
	}
)

type ErlangVMMetrics struct {
	erlangVMMemoryTotal                                 *prometheus.Desc
	erlangVMMemoryProcesses                             *prometheus.Desc
	erlangVMMemoryProcessesUsed                         *prometheus.Desc
	erlangVMMemorySystem                                *prometheus.Desc
	erlangVMMemoryAtom                                  *prometheus.Desc
	erlangVMMemoryAtomUsed                              *prometheus.Desc
	erlangVMMemoryBinary                                *prometheus.Desc
	erlangVMMemoryCode                                  *prometheus.Desc
	erlangVMMemoryEts                                   *prometheus.Desc
	erlangVMContextSwitches                             *prometheus.Desc
	erlangVMExactReductionsExactReductionsSinceLastCall *prometheus.Desc
	erlangVMExactReductionsTotalExactReductions         *prometheus.Desc
	erlangVMGarbageCollectionNumberOfGcs                *prometheus.Desc
	erlangVMGarbageCollectionWordsReclaimed             *prometheus.Desc
	erlangVMIoInput                                     *prometheus.Desc
	erlangVMIoOutput                                    *prometheus.Desc
	erlangVMReductionsReductionsSinceLastCall           *prometheus.Desc
	erlangVMReductionsTotalReductions                   *prometheus.Desc
	erlangVMRunQueue                                    *prometheus.Desc
	erlangVMRuntimeTimeSinceLastCall                    *prometheus.Desc
	erlangVMRuntimeTotalRunTime                         *prometheus.Desc
	erlangVMTotalActiveTasks                            *prometheus.Desc
	erlangVMTotalActiveTasksAll                         *prometheus.Desc
	erlangVMTotalRunQueueLengths                        *prometheus.Desc
	erlangVMTotalRunQueueLengthsAll                     *prometheus.Desc
	erlangVMWallClockTotalWallclockTime                 *prometheus.Desc
	erlangVMWallClockWallclockTimeSinceLastCall         *prometheus.Desc
	// XXX(tnamao): Sora GetStatsReport API が array を返す値は未対応です。
	// erlangVMActiveTasks                                 *prometheus.Desc
	// erlangVMActiveTasksAll                              *prometheus.Desc
	// erlangVMRunQueueLengths                             *prometheus.Desc
	// erlangVMRunQueueLengthsAll                          *prometheus.Desc
}

func (m *ErlangVMMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.erlangVMMemoryTotal
	ch <- m.erlangVMMemoryProcesses
	ch <- m.erlangVMMemoryProcessesUsed
	ch <- m.erlangVMMemorySystem
	ch <- m.erlangVMMemoryAtom
	ch <- m.erlangVMMemoryAtomUsed
	ch <- m.erlangVMMemoryBinary
	ch <- m.erlangVMMemoryCode
	ch <- m.erlangVMMemoryEts
	ch <- m.erlangVMContextSwitches
	ch <- m.erlangVMExactReductionsExactReductionsSinceLastCall
	ch <- m.erlangVMExactReductionsTotalExactReductions
	ch <- m.erlangVMGarbageCollectionNumberOfGcs
	ch <- m.erlangVMGarbageCollectionWordsReclaimed
	ch <- m.erlangVMIoInput
	ch <- m.erlangVMIoOutput
	ch <- m.erlangVMReductionsReductionsSinceLastCall
	ch <- m.erlangVMReductionsTotalReductions
	ch <- m.erlangVMRunQueue
	ch <- m.erlangVMRuntimeTimeSinceLastCall
	ch <- m.erlangVMRuntimeTotalRunTime
	ch <- m.erlangVMTotalActiveTasks
	ch <- m.erlangVMTotalActiveTasksAll
	ch <- m.erlangVMTotalRunQueueLengths
	ch <- m.erlangVMTotalRunQueueLengthsAll
	ch <- m.erlangVMWallClockTotalWallclockTime
	ch <- m.erlangVMWallClockWallclockTimeSinceLastCall
	// ch <- m.erlangVMActiveTasks
	// ch <- m.erlangVMActiveTasksAll
	// ch <- m.erlangVMRunQueueLengths
	// ch <- m.erlangVMRunQueueLengthsAll
}

func (m *ErlangVMMetrics) Collect(ch chan<- prometheus.Metric, report erlangVMReport) {
	ch <- newGauge(m.erlangVMMemoryTotal, float64(report.ErlangVMMemory.Total))
	ch <- newGauge(m.erlangVMMemoryProcesses, float64(report.ErlangVMMemory.Processes))
	ch <- newGauge(m.erlangVMMemoryProcessesUsed, float64(report.ErlangVMMemory.ProcessesUsed))
	ch <- newGauge(m.erlangVMMemorySystem, float64(report.ErlangVMMemory.System))
	ch <- newGauge(m.erlangVMMemoryAtom, float64(report.ErlangVMMemory.Atom))
	ch <- newGauge(m.erlangVMMemoryAtomUsed, float64(report.ErlangVMMemory.AtomUsed))
	ch <- newGauge(m.erlangVMMemoryBinary, float64(report.ErlangVMMemory.Binary))
	ch <- newGauge(m.erlangVMMemoryCode, float64(report.ErlangVMMemory.Code))
	ch <- newGauge(m.erlangVMMemoryEts, float64(report.ErlangVMMemory.Ets))
	ch <- newCounter(m.erlangVMContextSwitches, float64(report.ErlangVMStatistics.ContextSwitches))
	ch <- newGauge(m.erlangVMExactReductionsExactReductionsSinceLastCall, float64(report.ErlangVMStatistics.ExactReductions.ExactReductionsSinceLastCall))
	ch <- newCounter(m.erlangVMExactReductionsTotalExactReductions, float64(report.ErlangVMStatistics.ExactReductions.TotalExactReductions))
	ch <- newGauge(m.erlangVMGarbageCollectionNumberOfGcs, float64(report.ErlangVMStatistics.GarbageCollection.NumberOfGcs))
	ch <- newGauge(m.erlangVMGarbageCollectionWordsReclaimed, float64(report.ErlangVMStatistics.GarbageCollection.WordsReclaimed))
	ch <- newCounter(m.erlangVMIoInput, float64(report.ErlangVMStatistics.Io.Input))
	ch <- newCounter(m.erlangVMIoOutput, float64(report.ErlangVMStatistics.Io.Output))
	ch <- newGauge(m.erlangVMReductionsReductionsSinceLastCall, float64(report.ErlangVMStatistics.Reductions.ReductionsSinceLastCall))
	ch <- newCounter(m.erlangVMReductionsTotalReductions, float64(report.ErlangVMStatistics.Reductions.TotalReductions))
	ch <- newGauge(m.erlangVMRunQueue, float64(report.ErlangVMStatistics.RunQueue))
	ch <- newGauge(m.erlangVMRuntimeTimeSinceLastCall, float64(report.ErlangVMStatistics.Runtime.TimeSinceLastCall))
	ch <- newGauge(m.erlangVMRuntimeTotalRunTime, float64(report.ErlangVMStatistics.Runtime.TotalRunTime))
	ch <- newGauge(m.erlangVMTotalActiveTasks, float64(report.ErlangVMStatistics.TotalActiveTasks))
	ch <- newGauge(m.erlangVMTotalActiveTasksAll, float64(report.ErlangVMStatistics.TotalActiveTasksAll))
	ch <- newGauge(m.erlangVMTotalRunQueueLengths, float64(report.ErlangVMStatistics.TotalRunQueueLengths))
	ch <- newGauge(m.erlangVMTotalRunQueueLengthsAll, float64(report.ErlangVMStatistics.TotalRunQueueLengthsAll))
	ch <- newGauge(m.erlangVMWallClockTotalWallclockTime, float64(report.ErlangVMStatistics.WallClock.TotalWallclockTime))
	ch <- newGauge(m.erlangVMWallClockWallclockTimeSinceLastCall, float64(report.ErlangVMStatistics.WallClock.WallclockTimeSinceLastCall))
	// ch <- newGauge(m.erlangVMActiveTasks, float64(report.ErlangVMActiveTasks))
	// ch <- newGauge(m.erlangVMActiveTasksAll, float64(report.ErlangVMActiveTasksAll))
	// ch <- newGauge(m.erlangVMRunQueueLengths, float64(report.ErlangVMRunQueueLengths))
	// ch <- newGauge(m.erlangVMRunQueueLengthsAll, float64(report.ErlangVMRunQueueLengthsAll))
}
