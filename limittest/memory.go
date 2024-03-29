package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/zimmski/osutil"
)

func main() {
	memoryLimitInMiB, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	memoryToAllocateInMiB, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	osutil.EnforceProcessTreeLimits(osutil.ProcessTreeLimits{
		MaxMemoryInMiB: uint(memoryLimitInMiB),
		OnMemoryLimitReached: func(currentMemoryInMiB uint, maxMemoryInMiB uint) {
			os.Exit(5)
		},
	})

	// Consume requested memory at roughly 20 MiB/s.
	mebibytes := []byte{}
	for i := 0; i < memoryToAllocateInMiB; i++ {
		time.Sleep(50 * time.Millisecond)
		mebibytes = append(mebibytes, make([]byte, 1024*1024)...)
		fmt.Printf("%d\n", i+1)
	}

	// Write a flag.
	if err := os.WriteFile("success.txt", []byte("process finished successfully"), 0640); err != nil {
		panic(err)
	}
}
