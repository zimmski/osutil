package osutil

// BatchFileExtension returns the common file extension of a batch file.
func BatchFileExtension() (extension string) {
	return ".bat"
}

// BinaryExtension returns the common file extension of a binary.
func BinaryExtension() (extension string) {
	return ".exe"
}

// CommandFileExtension returns the common file extension of a command file.
func CommandFileExtension() (extension string) {
	return ".cmd"
}
