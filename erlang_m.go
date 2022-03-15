package main

import "github.com/prometheus/client_golang/prometheus"

var (
	erlangVmMetrics = ErlangVmMetrics{
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
)

type ErlangVmMetrics struct {
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
	// XXX(tnamao): Sora GetStatsReport API が array を返す値は未対応です。
	// erlangVmActiveTasks                                 *prometheus.Desc
	// erlangVmActiveTasksAll                              *prometheus.Desc
	// erlangVmRunQueueLengths                             *prometheus.Desc
	// erlangVmRunQueueLengthsAll                          *prometheus.Desc
}

func (m *ErlangVmMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.erlangVmMemoryTotal
	ch <- m.erlangVmMemoryProcesses
	ch <- m.erlangVmMemoryProcessesUsed
	ch <- m.erlangVmMemorySystem
	ch <- m.erlangVmMemoryAtom
	ch <- m.erlangVmMemoryAtomUsed
	ch <- m.erlangVmMemoryBinary
	ch <- m.erlangVmMemoryCode
	ch <- m.erlangVmMemoryEts
	ch <- m.erlangVmContextSwitches
	ch <- m.erlangVmExactReductionsExactReductionsSinceLastCall
	ch <- m.erlangVmExactReductionsTotalExactReductions
	ch <- m.erlangVmGarbageCollectionNumberOfGcs
	ch <- m.erlangVmGarbageCollectionWordsReclaimed
	ch <- m.erlangVmIoInput
	ch <- m.erlangVmIoOutput
	ch <- m.erlangVmReductionsReductionsSinceLastCall
	ch <- m.erlangVmReductionsTotalReductions
	ch <- m.erlangVmRunQueue
	ch <- m.erlangVmRuntimeTimeSinceLastCall
	ch <- m.erlangVmRuntimeTotalRunTime
	ch <- m.erlangVmTotalActiveTasks
	ch <- m.erlangVmTotalActiveTasksAll
	ch <- m.erlangVmTotalRunQueueLengths
	ch <- m.erlangVmTotalRunQueueLengthsAll
	ch <- m.erlangVmWallClockTotalWallclockTime
	ch <- m.erlangVmWallClockWallclockTimeSinceLastCall
	// ch <- m.erlangVmActiveTasks
	// ch <- m.erlangVmActiveTasksAll
	// ch <- m.erlangVmRunQueueLengths
	// ch <- m.erlangVmRunQueueLengthsAll
}

func (m *ErlangVmMetrics) Collect(ch chan<- prometheus.Metric, report erlangVmReport) {
	ch <- newGauge(m.erlangVmMemoryTotal, float64(report.ErlangVmMemoryTotal))
	ch <- newGauge(m.erlangVmMemoryProcesses, float64(report.ErlangVmMemoryProcesses))
	ch <- newGauge(m.erlangVmMemoryProcessesUsed, float64(report.ErlangVmMemoryProcessesUsed))
	ch <- newGauge(m.erlangVmMemorySystem, float64(report.ErlangVmMemorySystem))
	ch <- newGauge(m.erlangVmMemoryAtom, float64(report.ErlangVmMemoryAtom))
	ch <- newGauge(m.erlangVmMemoryAtomUsed, float64(report.ErlangVmMemoryAtomUsed))
	ch <- newGauge(m.erlangVmMemoryBinary, float64(report.ErlangVmMemoryBinary))
	ch <- newGauge(m.erlangVmMemoryCode, float64(report.ErlangVmMemoryCode))
	ch <- newGauge(m.erlangVmMemoryEts, float64(report.ErlangVmMemoryEts))
	ch <- newGauge(m.erlangVmContextSwitches, float64(report.ErlangVmContextSwitches))
	ch <- newGauge(m.erlangVmExactReductionsExactReductionsSinceLastCall, float64(report.ErlangVmExactReductionsExactReductionsSinceLastCall))
	ch <- newGauge(m.erlangVmExactReductionsTotalExactReductions, float64(report.ErlangVmExactReductionsTotalExactReductions))
	ch <- newGauge(m.erlangVmGarbageCollectionNumberOfGcs, float64(report.ErlangVmGarbageCollectionNumberOfGcs))
	ch <- newGauge(m.erlangVmGarbageCollectionWordsReclaimed, float64(report.ErlangVmGarbageCollectionWordsReclaimed))
	ch <- newGauge(m.erlangVmIoInput, float64(report.ErlangVmIoInput))
	ch <- newGauge(m.erlangVmIoOutput, float64(report.ErlangVmIoOutput))
	ch <- newGauge(m.erlangVmReductionsReductionsSinceLastCall, float64(report.ErlangVmReductionsReductionsSinceLastCall))
	ch <- newGauge(m.erlangVmReductionsTotalReductions, float64(report.ErlangVmReductionsTotalReductions))
	ch <- newGauge(m.erlangVmRunQueue, float64(report.ErlangVmRunQueue))
	ch <- newGauge(m.erlangVmRuntimeTimeSinceLastCall, float64(report.ErlangVmRuntimeTimeSinceLastCall))
	ch <- newGauge(m.erlangVmRuntimeTotalRunTime, float64(report.ErlangVmRuntimeTotalRunTime))
	ch <- newGauge(m.erlangVmTotalActiveTasks, float64(report.ErlangVmTotalActiveTasks))
	ch <- newGauge(m.erlangVmTotalActiveTasksAll, float64(report.ErlangVmTotalActiveTasksAll))
	ch <- newGauge(m.erlangVmTotalRunQueueLengths, float64(report.ErlangVmTotalRunQueueLengths))
	ch <- newGauge(m.erlangVmTotalRunQueueLengthsAll, float64(report.ErlangVmTotalRunQueueLengthsAll))
	ch <- newGauge(m.erlangVmWallClockTotalWallclockTime, float64(report.ErlangVmWallClockTotalWallclockTime))
	ch <- newGauge(m.erlangVmWallClockWallclockTimeSinceLastCall, float64(report.ErlangVmWallClockWallclockTimeSinceLastCall))
	// ch <- newGauge(m.erlangVmActiveTasks, float64(report.ErlangVmActiveTasks))
	// ch <- newGauge(m.erlangVmActiveTasksAll, float64(report.ErlangVmActiveTasksAll))
	// ch <- newGauge(m.erlangVmRunQueueLengths, float64(report.ErlangVmRunQueueLengths))
	// ch <- newGauge(m.erlangVmRunQueueLengthsAll, float64(report.ErlangVmRunQueueLengthsAll))
}
