package osutil

/*
#include <stdio.h>
#include <stdlib.h>
*/
import "C"

import (
	"bytes"
	"golang.org/x/sys/unix"
	"io"
	"os"
	"sync"
	"syscall"
)

var capturingFileDescLock = &sync.Mutex{}

// Capture captures stderr and stdout of a given function call.
//
// Note that because this requires modifying stdout and stderr,
// which are global to the process,
// this function (and the invocation of call) are wrapped in a mutex.
// (note also that this mutex is shared between Capture and CaptureWithCGo).
//
// Otherwise it would be possible for output to contain incorrect bytes,
// either missing what was written during call,
// and/or having the results of other invocations.
func Capture(call func()) (output []byte, err error) {
	capturingFileDescLock.Lock()
	defer func() {
		capturingFileDescLock.Unlock()
	}()

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
		if w != nil {
			e = w.Close()
			if err != nil {
				err = e
			}
		}
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

	call()

	err = w.Close()
	if err != nil {
		return nil, err
	}
	w = nil

	return <-out, err
}

// CaptureWithCGo captures stderr and stdout as well as stderr and stdout of C of a given function call.
//
// Note that because this requires modifying stdout and stderr,
// which are global to the process,
// this function (and the invocation of call) are wrapped in a mutex.
// (note also that this mutex is shared between Capture and CaptureWithCGo).
//
// Otherwise it would be possible for output to contain incorrect bytes,
// either missing what was written during call,
// and/or having the results of other invocations.

func CaptureWithCGo(call func()) (output []byte, err error) {
	capturingFileDescLock.Lock()
	defer func() {
		capturingFileDescLock.Unlock()
	}()

	originalStdout, e := syscall.Dup(syscall.Stdout)
	if e != nil {
		return nil, e
	}

	originalStderr, e := syscall.Dup(syscall.Stderr)
	if e != nil {
		return nil, e
	}

	defer func() {
		if e := unix.Dup2(originalStdout, syscall.Stdout); e != nil {
			err = e
		}
		if e := syscall.Close(originalStdout); e != nil {
			err = e
		}
		if e := unix.Dup2(originalStderr, syscall.Stderr); e != nil {
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
		if w != nil {
			e = w.Close()
			if err != nil {
				err = e
			}
		}
	}()

	if e := unix.Dup2(int(w.Fd()), syscall.Stdout); e != nil {
		return nil, e
	}
	if e := unix.Dup2(int(w.Fd()), syscall.Stderr); e != nil {
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

	call()

	C.fflush(C.stdout)

	err = w.Close()
	if err != nil {
		return nil, err
	}
	w = nil

	if e := syscall.Close(syscall.Stdout); e != nil {
		return nil, e
	}
	if e := syscall.Close(syscall.Stderr); e != nil {
		return nil, e
	}
	return <-out, err
}
