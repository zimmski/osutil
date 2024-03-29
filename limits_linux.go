//go:build linux

package osutil

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

// EnforceProcessTreeLimits constrains the current process and all descendant processes to the specified limits. The current process exits when the limits are exceeded.
func EnforceProcessTreeLimits(limits ProcessTreeLimits) {
	if limits.MaxMemoryInMiB <= 0 {
		return
	}

	var watchdogInterval time.Duration
	if limits.WatchdogInterval == 0 {
		watchdogInterval = 2 * time.Second
	} else {
		watchdogInterval = limits.WatchdogInterval
	}

	go func() {
		for {
			memoryUsageInKiB, err := getProcessTreeMemoryUsage()
			if err != nil {
				panic(fmt.Errorf("Failed to check memory usage: %w", err))
			}

			currentMemoryInMiB := memoryUsageInKiB / 1024
			if currentMemoryInMiB > limits.MaxMemoryInMiB {
				limits.OnMemoryLimitReached(currentMemoryInMiB, limits.MaxMemoryInMiB)
			}
			time.Sleep(watchdogInterval)
		}
	}()
}

// getProcessTreeMemoryUsage returns the total memory usage in KiB from the current process and all child processes.
//
// REMARK This is currently a rough approximation. Memory shared between descendant processes is counted multiple times.
func getProcessTreeMemoryUsage() (memoryUsageInKiB uint, err error) {
	psCmd := exec.Command("ps", "-H", "-o", "rss=", strconv.Itoa(os.Getpid()))
	psOutput, err := psCmd.CombinedOutput()
	if err != nil {
		return 0, err
	}
	lines := bytes.Split(bytes.TrimSpace(psOutput), []byte("\n"))
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			return 0, errors.New("Encountered empty line in \"ps\" output")
		}
		rss, err := strconv.Atoi(string(line))
		if err != nil {
			return 0, err
		}
		memoryUsageInKiB += uint(rss)
	}

	return memoryUsageInKiB, nil
}

// SetRLimitFiles temporarily changes the file descriptor resource limit while calling the given function.
func SetRLimitFiles(limit uint64, call func(limit uint64)) (err error) {
	var tmp syscall.Rlimit
	if err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &tmp); err != nil {
		return nil
	}
	defer func() {
		if err == nil {
			err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &tmp)
		}
	}()

	if err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &syscall.Rlimit{
		Cur: limit,
		Max: tmp.Max,
	}); err != nil {
		return err
	}

	call(limit)

	return nil
}
