package osutil

/*
#include <stdio.h>
#include <stdlib.h>
*/
import "C"

import (
	"bytes"
	"io"
	"os"
	"sync"
	"syscall"
)

var lockStdFileDescriptorsSwapping sync.Mutex // FIXME our solution is not concurrent-safe. Find a better solution because this might be a bottleneck in the future.

// Capture captures stderr and stdout of a given function call.
func Capture(call func()) (output []byte, err error) {
	lockStdFileDescriptorsSwapping.Lock()
	defer lockStdFileDescriptorsSwapping.Unlock()

	originalStdErr, originalStdOut := os.Stderr, os.Stdout
	defer func() {
		os.Stderr, os.Stdout = originalStdErr, originalStdOut
	}()

	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	defer func() {
		e := r.Close()
		if e != nil {
			err = e
		}
	}()

	os.Stderr, os.Stdout = w, w

	out := make(chan []byte)
	go func() {
		var b bytes.Buffer

		_, err := io.Copy(&b, r)
		if err != nil {
			panic(err)
		}

		out <- b.Bytes()
	}()

	call()

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return <-out, err
}

// CaptureWithCGo captures stderr and stdout as well as stderr and stdout of C of a given function call.
func CaptureWithCGo(call func()) (output []byte, err error) {
	lockStdFileDescriptorsSwapping.Lock()
	defer lockStdFileDescriptorsSwapping.Unlock()

	originalStdout, e := syscall.Dup(syscall.Stdout)
	if e != nil {
		return nil, e
	}

	originalStderr, e := syscall.Dup(syscall.Stderr)
	if e != nil {
		return nil, e
	}

	defer func() {
		if e := syscall.Dup2(originalStdout, syscall.Stdout); e != nil {
			err = e
		}
		if e := syscall.Close(originalStdout); e != nil {
			err = e
		}
		if e := syscall.Dup2(originalStderr, syscall.Stderr); e != nil {
			err = e
		}
		if e := syscall.Close(originalStderr); e != nil {
			err = e
		}
	}()

	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	defer func() {
		e := r.Close()
		if e != nil {
			err = e
		}
	}()

	if e := syscall.Dup2(int(w.Fd()), syscall.Stdout); e != nil {
		return nil, e
	}
	if e := syscall.Dup2(int(w.Fd()), syscall.Stderr); e != nil {
		return nil, e
	}

	out := make(chan []byte)
	go func() {
		var b bytes.Buffer

		_, err := io.Copy(&b, r)
		if err != nil {
			panic(err)
		}

		out <- b.Bytes()
	}()

	call()

	C.fflush(C.stdout)

	err = w.Close()
	if err != nil {
		return nil, err
	}
	if e := syscall.Close(syscall.Stdout); e != nil {
		return nil, e
	}
	if e := syscall.Close(syscall.Stderr); e != nil {
		return nil, e
	}

	return <-out, err
}
