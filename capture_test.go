package osutil

import (
	"testing"
)

func TestCapture(t *testing.T) {
	testCapture(t)
}

func TestCaptureWithCGo(t *testing.T) {
	testCaptureWithCGo(t)
}
