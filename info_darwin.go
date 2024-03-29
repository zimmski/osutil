//go:build darwin

package osutil

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/pkg/errors"
)

const (
	// KernelVersionIdentifier holds the identifier for the kernel version OS information.
	KernelVersionIdentifier = "KernelVersion"
	// OperatingSystemBuildIdentifier holds the identifier for the OS build OS information.
	OperatingSystemBuildIdentifier = "BuildVersion"
	// OperatingSystemIdentifier holds the identifier for the OS name OS information.
	OperatingSystemIdentifier = "ProductName"
	// OperatingSystemVersionIdentifier holds the identifier for the OS version OS information.
	OperatingSystemVersionIdentifier = "ProductVersion"
)

// Info returns a list of OS relevant information.
func Info() (info map[string]string, err error) {
	info = map[string]string{}

	kernelVersion, err := syscall.Sysctl("kern.osrelease")
	if err != nil {
		return nil, errors.Wrap(err, "failed to query kernel version")
	}
	info[KernelVersionIdentifier] = strings.TrimSpace(kernelVersion)

	softwareVersionsCommand := exec.Command("sw_vers")
	var softwareVersions bytes.Buffer
	softwareVersionsCommand.Stdout = &softwareVersions
	if err := softwareVersionsCommand.Run(); err != nil {
		return nil, errors.Wrap(err, "failed to query software versions")
	}
	for _, line := range strings.Split(softwareVersions.String(), "\n") {
		ls := strings.SplitN(line, ":", 2)
		if len(ls) == 1 { // Ignore empty lines
			continue
		}
		info[strings.TrimSpace(ls[0])] = strings.TrimSpace(ls[1])
	}

	info[EnvironmentPathIdentifier] = os.Getenv(EnvironmentPathIdentifier)

	return info, err
}
