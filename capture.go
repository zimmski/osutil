package osutil

/*
#include <stdio.h>
#include <stdlib.h>
*/
import "C"

import (
	"bytes"
	"errors"
	"io"
	"os"
	"sync"
	"unsafe"
)

var (
	ErrFDOpenFailed = errors.New("fdopen returned nil")
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

	originalStdErr, originalStdOut := os.Stderr, os.Stdout
	originalCStdErr, originalCStdOut := C.stderr, C.stdout
	defer func() {
		os.Stderr, os.Stdout = originalStdErr, originalStdOut
		C.stderr, C.stdout = originalCStdErr, originalCStdOut
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

	cw := C.CString("w")
	defer C.free(unsafe.Pointer(cw))

	f := C.fdopen((C.int)(w.Fd()), cw)
	if f == nil {
		return nil, ErrFDOpenFailed
	}
	defer C.fclose(f)

	os.Stderr, os.Stdout = w, w
	C.stderr, C.stdout = f, f

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

	C.fflush(f)

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return <-out, err
}
