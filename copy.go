package osutil

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/termie/go-shutil"
)

// CopyFile copies a file from src to dst.
func CopyFile(src string, dst string) (err error) {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		e := s.Close()
		if err == nil {
			err = e
		}
	}()

	d, err := os.Create(dst)
	if err != nil {
		// In case the file is a symlink, we need to remove the file before we can write to it.
		if _, e := os.Lstat(dst); e == nil {
			if e := os.Remove(dst); e != nil {
				return e
			}
			d, err = os.Create(dst)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer func() {
		e := d.Close()
		if err == nil {
			err = e
		}
	}()

	_, err = io.Copy(d, s)
	if err != nil {
		return err
	}

	i, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, i.Mode())
}

// CopyFileCompressed reads the file src and writes a compressed version to dst.
// The compression level can be gzip.DefaultCompression, gzip.NoCompression, gzip.HuffmanOnly or any integer value between gzip.BestSpeed and gzip.BestCompression inclusive.
func CopyFileCompressed(src string, dst string, compressionLevel int) (err error) {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		e := s.Close()
		if err == nil {
			err = e
		}
	}()

	d, err := os.Create(dst)
	if err != nil {
		// In case the file is a symlink, we need to remove the file before we can write to it.
		if _, e := os.Lstat(dst); e == nil {
			if e := os.Remove(dst); e != nil {
				return e
			}
			d, err = os.Create(dst)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer func() {
		e := d.Close()
		if err == nil {
			err = e
		}
	}()

	gzipWriter, err := gzip.NewWriterLevel(d, compressionLevel)
	if err != nil {
		return err
	}
	defer func() {
		e := gzipWriter.Close()
		if err == nil {
			err = e
		}
	}()

	if _, err := io.Copy(gzipWriter, s); err != nil {
		return err
	}

	i, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, i.Mode())
}

// CopyTree copies a whole file system tree from the source path to destination path.
func CopyTree(sourcePath string, destinationPath string) (err error) {
	return shutil.CopyTree(sourcePath, destinationPath, nil)
}

// CompressDirectory reads the directory srcDirectory and writes a compressed version to archive.
func CompressDirectory(srcDirectory string, archive string) (err error) {
	archiveFile, err := os.Create(archive)
	if err != nil {
		return err
	}
	defer func() {
		e := archiveFile.Close()
		if err == nil {
			err = e
		}
	}()

	zipWriter := zip.NewWriter(archiveFile)
	defer func() {
		e := zipWriter.Close()
		if err == nil {
			err = e
		}
	}()

	return filepath.WalkDir(srcDirectory, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			e := file.Close()
			if err == nil {
				err = e
			}
		}()

		relativePath, err := filepath.Rel(srcDirectory, path)
		if err != nil {
			return err
		}

		zipFileWriter, err := zipWriter.Create(relativePath)
		if err != nil {
			return err
		}

		if _, err := io.Copy(zipFileWriter, file); err != nil {
			return err
		}

		return nil
	})
}

// Uncompress extracts the given archive into the given destination.
func Uncompress(archive io.Reader, dstDirectory string) (err error) {
	data, err := io.ReadAll(archive)
	if err != nil {
		return err
	}
	byteReader := bytes.NewReader(data)

	zipReader, err := zip.NewReader(byteReader, byteReader.Size())
	if err != nil {
		return err
	}

	for _, zipFile := range zipReader.File {
		zipReaderFile, err := zipFile.Open()
		if err != nil {
			return err
		}

		destinationPath := filepath.Join(dstDirectory, zipFile.Name)
		if err := os.MkdirAll(filepath.Dir(destinationPath), 0700); err != nil {
			return err
		}
		destinationFile, err := os.Create(destinationPath)
		if err != nil {
			return err
		}

		if _, err := io.Copy(destinationFile, zipReaderFile); err != nil {
			_ = destinationFile.Close()
			_ = zipReaderFile.Close()

			return err
		}

		if err := destinationFile.Close(); err != nil {
			_ = zipReaderFile.Close()

			return err
		}
		if err := zipReaderFile.Close(); err != nil {
			return err
		}

		if err := os.Chmod(destinationPath, zipFile.Mode()); err != nil {
			return err
		}
	}

	return nil
}
