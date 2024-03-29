package osutil

import (
	"io"
	"os"
)

// InMemoryStream allows to read and write to an in-memory stream.
type InMemoryStream struct {
	// reader holds the reading end of the stream.
	reader io.ReadCloser
	// writer holds the writing end of the stream.
	writer io.WriteCloser
}

var _ io.ReadWriteCloser = (*InMemoryStream)(nil)

func NewInMemoryStream(reader io.ReadCloser, writer io.WriteCloser) *InMemoryStream {
	return &InMemoryStream{
		reader: reader,
		writer: writer,
	}
}

// Read reads from the stream until the given buffer is full.
func (s *InMemoryStream) Read(buffer []byte) (n int, err error) {
	return s.reader.Read(buffer)
}

// Write writes the given data to the stream.
func (s *InMemoryStream) Write(data []byte) (n int, err error) {
	return s.writer.Write(data)
}

// Close closes the stream.
func (s *InMemoryStream) Close() error {
	if err := s.reader.Close(); err != nil {
		return err
	}
	if err := s.writer.Close(); err != nil {
		return err
	}

	return nil
}

// StandardStream allows to read from STDIN and write to STDOUT.
type StandardStream struct{}

var _ io.ReadWriteCloser = (*StandardStream)(nil)

// Read reads from the stream until the given buffer is full.
func (s *StandardStream) Read(buffer []byte) (n int, err error) {
	return os.Stdin.Read(buffer)
}

// Write writes the given data to the stream.
func (s *StandardStream) Write(data []byte) (n int, err error) {
	return os.Stdout.Write(data)
}

// Close closes the stream.
func (s *StandardStream) Close() error {
	if err := os.Stdin.Close(); err != nil {
		return err
	}
	if err := os.Stdout.Close(); err != nil {
		return err
	}

	return nil
}
