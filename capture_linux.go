//go:build cgo

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
	"os/signal"
	"sync"
	"syscall"
)

var lockStdFileDescriptorsSwapping sync.Mutex
var lockStdFileWithCGoDescriptorsSwapping sync.Mutex

// Capture captures stderr and stdout of a given function call.
func Capture(call func()) (output []byte, err error) {
	lockStdFileDescriptorsSwapping.Lock()

	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	defer func() {
		e := r.Close()
		if e != nil {
			err = e
		}
		if w != nil {
			e = w.Close()
			if err != nil {
				err = e
			}
		}
	}()

	originalStdErr, originalStdOut := os.Stderr, os.Stdout
	release := func() {
		os.Stderr, os.Stdout = originalStdErr, originalStdOut
	}
	defer func() {
		lockStdFileDescriptorsSwapping.Lock()

		if w != nil {
			release()
		}

		lockStdFileDescriptorsSwapping.Unlock()
	}()
	os.Stderr, os.Stdout = w, w

	out := make(chan []byte)
	go func() {
		defer func() {
			// If there is a panic in the function call, copying from "r" does not work anymore.
			_ = recover()
		}()

		var b bytes.Buffer

		_, err := io.Copy(&b, r)
		if err != nil {
			panic(err)
		}

		out <- b.Bytes()
	}()

	lockStdFileDescriptorsSwapping.Unlock()

	call()

	lockStdFileDescriptorsSwapping.Lock()

	err = w.Close()
	if err != nil {
		lockStdFileDescriptorsSwapping.Unlock()

		return nil, err
	}
	w = nil

	release()

	lockStdFileDescriptorsSwapping.Unlock()

	return <-out, err
}

// CaptureWithCGo captures stderr and stdout as well as stderr and stdout of C of a given function call.
// Currently this function cannot be nested.
func CaptureWithCGo(call func()) (output []byte, err error) {
	// FIXME At the moment this function does not work with nested calls (recursively). This might be because of the signal handler or the way we clone the file descriptors. I really do not know. Since we do not need recursive calls right now we can postpone this for later. https://$INTERNAL/symflower/symflower/-/issues/85
	lockStdFileWithCGoDescriptorsSwapping.Lock()
	defer lockStdFileWithCGoDescriptorsSwapping.Unlock()

	lockStdFileDescriptorsSwapping.Lock()

	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	defer func() {
		e := r.Close()
		if e != nil {
			err = e
		}
		if w != nil {
			e = w.Close()
			if err != nil {
				err = e
			}
		}
	}()

	originalStdout, err := syscall.Dup(syscall.Stdout)
	if err != nil {
		lockStdFileDescriptorsSwapping.Unlock()

		return nil, err
	}
	originalStderr, err := syscall.Dup(syscall.Stderr)
	if err != nil {
		lockStdFileDescriptorsSwapping.Unlock()

		return nil, err
	}
	release := func() {
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
	}
	defer func() {
		lockStdFileDescriptorsSwapping.Lock()

		if w != nil {
			release()
		}

		lockStdFileDescriptorsSwapping.Unlock()
	}()

	// WORKAROUND Since Go 1.10 `go test` can hang (randomly) if a subprocess has another subprocess that got the original STDOUT/STDERR if the defined STDOUT/STDERR are not a file from the Go side. This is a somewhat incomplete description of this bug. More details can be found here https://github.com/golang/go/issues/24050 and here https://github.com/golang/go/issues/23019. This bug occurs randomly not just with `go test` but anywhere the same APIs are used. However, the only time this happens with the current function is when a parent process kills the currently running process but not when we have, e.g. a panic in the call we are capturing. This is already handled by the defer calls. Since the exiting is not handled, we have to set up a signal handler to take care of the cleanup.

	exitSignalHandler := make(chan bool)
	sigs := make(chan os.Signal, 10)
	signal.Notify(sigs, syscall.SIGCHLD)
	defer func() {
		signal.Stop(sigs)
		exitSignalHandler <- true
	}()

	go func() {
		select {
		case <-sigs:
			_ = syscall.Close(originalStdout)
			_ = syscall.Close(originalStderr)
		case <-exitSignalHandler:
		}
	}()

	if e := syscall.Dup2(int(w.Fd()), syscall.Stdout); e != nil {
		lockStdFileDescriptorsSwapping.Unlock()

		return nil, e
	}
	if e := syscall.Dup2(int(w.Fd()), syscall.Stderr); e != nil {
		lockStdFileDescriptorsSwapping.Unlock()

		return nil, e
	}

	out := make(chan []byte)
	go func() {
		defer func() {
			// If there is a panic in the function call, copying from "r" does not work anymore.
			_ = recover()
		}()

		var b bytes.Buffer

		_, err := io.Copy(&b, r)
		if err != nil {
			panic(err)
		}

		out <- b.Bytes()
	}()

	lockStdFileDescriptorsSwapping.Unlock()

	call()

	lockStdFileDescriptorsSwapping.Lock()

	C.fflush(C.stderr)
	C.fflush(C.stdout)

	if err = w.Close(); err != nil {
		lockStdFileDescriptorsSwapping.Unlock()

		return nil, err
	}
	w = nil

	if err = syscall.Close(syscall.Stdout); err != nil {
		lockStdFileDescriptorsSwapping.Unlock()

		return nil, err
	}
	if err = syscall.Close(syscall.Stderr); err != nil {
		lockStdFileDescriptorsSwapping.Unlock()

		return nil, err
	}

	release()

	lockStdFileDescriptorsSwapping.Unlock()

	return <-out, err
}
