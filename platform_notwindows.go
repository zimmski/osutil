//go:build !windows

package osutil

// BatchFileExtension returns the common file extension of a batch file.
func BatchFileExtension() (extension string) {
	return ""
}

// BinaryExtension returns the common file extension of a binary.
func BinaryExtension() (extension string) {
	return ""
}

// CommandFileExtension returns the common file extension of a command file.
func CommandFileExtension() (extension string) {
	return ""
}
