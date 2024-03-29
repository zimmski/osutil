//go:build linux

package osutil

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	// KernelVersionIdentifier holds the identifier for the kernel version OS information.
	KernelVersionIdentifier = "KernelVersion"
	// OperatingSystemIdentifier holds the identifier for the OS name OS information.
	OperatingSystemIdentifier = "ProductName"
	// OperatingSystemVersionIdentifier holds the identifier for the OS version OS information.
	OperatingSystemVersionIdentifier = "ProductVersion"
)

// Info returns a list of OS relevant information.
func Info() (info map[string]string, err error) {
	info = map[string]string{}

	kernelVersionCommand := exec.Command("uname", "-r")
	var kernelVersion bytes.Buffer
	kernelVersionCommand.Stdout = &kernelVersion
	if err := kernelVersionCommand.Run(); err != nil {
		return nil, errors.Wrap(err, "failed to query kernel versions")
	}
	info[KernelVersionIdentifier] = strings.TrimSpace(kernelVersion.String())

	osReleaseData, osReleaseError := os.ReadFile("/etc/os-release")
	if osReleaseError != nil {
		osReleaseData, err = os.ReadFile("/etc/lsb-release")
		if err != nil {
			return nil, errors.Wrap(err, "failed to query /etc/os-release")
		}
	}
	for _, line := range strings.Split(string(osReleaseData), "\n") {
		ls := strings.SplitN(line, "=", 2)
		if len(ls) == 1 { // Ignore empty lines
			continue
		}
		k := strings.TrimSpace(ls[0])
		v := strings.TrimSpace(ls[1])
		if v[0] == '"' {
			v, err = strconv.Unquote(v)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unquote data in /etc/os-release")
			}
		}
		switch k {
		case "NAME":
			info[OperatingSystemIdentifier] = v
		case "VERSION_ID":
			info[OperatingSystemVersionIdentifier] = v
		}
	}

	info[EnvironmentPathIdentifier] = os.Getenv(EnvironmentPathIdentifier)

	return info, err
}
