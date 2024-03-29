//go:build windows

package osutil

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows/registry"
)

const (
	// OperatingSystemBuildIdentifier holds the identifier for the OS build OS information.
	OperatingSystemBuildIdentifier = "BuildVersion"
	// OperatingSystemIdentifier holds the identifier for the OS name OS information.
	OperatingSystemIdentifier = "ProductName"
	// OperatingSystemVersionIdentifier holds the identifier for the OS version OS information.
	OperatingSystemVersionIdentifier = "ProductVersion"
	// OperatingSystemVersionMajorIdentifier holds the identifier for the OS major version OS information.
	OperatingSystemVersionMajorIdentifier = "ProductVersionMajor"
	// OperatingSystemVersionMinorIdentifier holds the identifier for the OS minor version OS information.
	OperatingSystemVersionMinorIdentifier = "ProductVersionMinor"
)

const (
	registryKeyBuildUpdateVersion  = "UBR"
	registryKeyBuildVersion        = "CurrentBuildNumber"
	registryKeyProductMajorVersion = "CurrentMajorVersionNumber"
	registryKeyProductVersionMinor = "CurrentMinorVersionNumber"
	registryKeyProductName         = "ProductName"
	registryKeyProductVersion      = "CurrentVersion"
)

// Info returns a list of OS relevant information.
func Info() (info map[string]string, err error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open registry")
	}
	defer k.Close()

	info = map[string]string{}

	info[OperatingSystemBuildIdentifier], _, err = k.GetStringValue(registryKeyBuildVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to query key %s: %w", registryKeyBuildVersion, err)
	}
	{
		data, _, err := k.GetIntegerValue(registryKeyBuildUpdateVersion)
		if err != nil && err != registry.ErrNotExist {
			return nil, fmt.Errorf("failed to query key %s: %w", registryKeyBuildUpdateVersion, err)
		} else if err == nil {
			info[OperatingSystemBuildIdentifier] += "." + strconv.FormatUint(data, 10)
		}
	}

	info[OperatingSystemIdentifier], _, err = k.GetStringValue(registryKeyProductName)
	if err != nil {
		return nil, fmt.Errorf("failed to query key %s: %w", registryKeyProductName, err)
	}

	info[OperatingSystemVersionIdentifier], _, err = k.GetStringValue(registryKeyProductVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to query key %s: %w", registryKeyProductVersion, err)
	}
	{
		data, _, err := k.GetIntegerValue(registryKeyProductMajorVersion)
		if err != nil {
			return nil, fmt.Errorf("failed to query key %s: %w", registryKeyProductMajorVersion, err)
		}
		info[OperatingSystemVersionMajorIdentifier] = strconv.FormatUint(data, 10)
	}
	{
		data, _, err := k.GetIntegerValue(registryKeyProductVersionMinor)
		if err != nil {
			return nil, fmt.Errorf("failed to query key %s: %w", registryKeyProductVersionMinor, err)
		}
		info[OperatingSystemVersionMinorIdentifier] = strconv.FormatUint(data, 10)
	}

	info[EnvironmentPathIdentifier] = os.Getenv(EnvironmentPathIdentifier)

	return info, err
}
