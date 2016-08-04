package osutil

/*
#include <stdio.h>
*/
import "C"

import (
	"bytes"
	"io"
	"os"
)

// Capture captures stderr and stdout of a given function call.
func Capture(call func()) ([]byte, error) {
	originalStdErr, originalStdOut := os.Stderr, os.Stdout
	defer func() {
		os.Stderr, os.Stdout = originalStdErr, originalStdOut
	}()

	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}

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
func CaptureWithCGo(call func()) ([]byte, error) {
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

	f := C.fdopen((C.int)(w.Fd()), C.CString("w"))

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
