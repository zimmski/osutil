package osutil

import (
	"os"
	"regexp"
	"time"
)

// StaticFile holds a single file or directory in-memory.
type StaticFile struct {
	Data  string
	Mime  string
	Mtime time.Time
	// Size is the size before compression.
	// If 0, it means the data is uncompressed.
	Size int
	// Hash is a SHA-256 hash of the file contents, which is used for the Etag, and useful for caching.
	Hash string
	// Directory determines if this file is a directory.
	Directory bool
}

// RewriteStaticIndexFile rewrites a `github.com/bouk/staticfiles` index file to be extendable by replacing inlined code to common code.
func RewriteStaticIndexFile(filePath string) (err error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	data = regexp.MustCompile(`(?s)(\t"golang.org/x/text/language"\n)`).ReplaceAll(data, []byte("\t\"github.com/zimmski/osutil\"\n$1"))
	data = regexp.MustCompile(`(?s)type StaticFilesFile struct {.+?}\n\s+`).ReplaceAll(data, []byte(""))
	data = regexp.MustCompile(`StaticFilesFile`).ReplaceAll(data, []byte("osutil.StaticFile"))

	if err := os.WriteFile(filePath, data, 0); err != nil {
		return err
	}

	return nil
}
