package osutil

import (
	"time"
)

// ProcessTreeLimits holds limits to apply to the current process's resource usage, including the resource usage of its descendant processes.
type ProcessTreeLimits struct {
	// MaxMemoryInMiB holds the limit for the memory usage of the current process and all descendants in 1024-based mebibytes.
	// Zero means no limit.
	MaxMemoryInMiB uint
	// OnOutOfMemory may or may not run when the memory limit is reached, depending on the enforcement strategy. If it runs, the process will not be killed automatically and the function should end the process.
	OnMemoryLimitReached func(currentMemoryInMiB uint, maxMemoryInMiB uint)
	// WatchdogInterval holds the amount of time to sleep between checks if the limits have been exceeded, if a watchdog strategy is used.
	// The default is two seconds.
	WatchdogInterval time.Duration
}
