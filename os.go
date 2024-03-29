package osutil

import "runtime"

// IsDarwin returns whether the operating system is Darwin.
func IsDarwin() bool {
	return runtime.GOOS == Darwin
}

// IsLinux returns whether the operating system is Linux.
func IsLinux() bool {
	return runtime.GOOS == Linux
}

// IsWindows returns whether the operating system is Windows.
func IsWindows() bool {
	return runtime.GOOS == Windows
}

// IsArchitectureARMWith32Bit returns wheter the operating system runs on ARM with 32 bits.
func IsArchitectureARMWith32Bit() bool {
	return runtime.GOARCH == "arm"
}

// IsArchitectureARMWith64Bit returns wheter the operating system runs on ARM with 64 bits.
func IsArchitectureARMWith64Bit() bool {
	return runtime.GOARCH == "arm64"
}

// IsArchitectureX86With32Bit returns wheter the operating system runs on x86 with 32 bits.
func IsArchitectureX86With32Bit() bool {
	return runtime.GOARCH == "386"
}

// IsArchitectureX86With64Bit returns wheter the operating system runs on x86 with 64 bits.
func IsArchitectureX86With64Bit() bool {
	return runtime.GOARCH == "amd64"
}
