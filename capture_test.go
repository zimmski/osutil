package osutil

import (
	"syscall"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCapture(t *testing.T) {
	var limit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit)
	assert.Nil(t, err)
	//fmt.Fprintf(os.Stderr, "%v file descriptors out of a maximum of %v available\n", limit.Cur, limit.Max) 
	// use limit.Cur or something like 1024
	for i :=0; i <= int(limit.Cur); i++ {
		testCapture(t)
	}
}

func TestCaptureWithCGo(t *testing.T) {
	var limit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit)
	assert.Nil(t, err)
	//fmt.Fprintf(os.Stderr, "%v file descriptors out of a maximum of %v available\n", limit.Cur, limit.Max) 
	// use limit.Cur or something like 512
	for i:=0; i <= int(limit.Cur); i++ {
		testCaptureWithCGo(t)
	}
}
