package bytesutil

import (
	"bytes"
	"io"
	"sync"
)

// SynchronizedBuffer holds a concurrency-safe buffer.
type SynchronizedBuffer struct {
	lock sync.Mutex

	b bytes.Buffer
}

// Bytes calls "Bytes" of "bytes.Buffer".
func (b *SynchronizedBuffer) Bytes() []byte {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.b.Bytes()
}

// Cap calls "Cap" of "bytes.Buffer".
func (b *SynchronizedBuffer) Cap() int {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.b.Cap()
}

// Grow calls "Grow" of "bytes.Buffer".
func (b *SynchronizedBuffer) Grow(n int) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.b.Grow(n)
}

// Len calls "Len" of "bytes.Buffer".
func (b *SynchronizedBuffer) Len() int {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.b.Len()
}

// Next calls "Next" of "bytes.Buffer".
func (b *SynchronizedBuffer) Next(n int) []byte {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.b.Next(n)
}

// Read calls "Read" of "bytes.Buffer".
func (b *SynchronizedBuffer) Read(p []byte) (n int, err error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.b.Read(p)
}

// ReadByte calls "ReadByte" of "bytes.Buffer".
func (b *SynchronizedBuffer) ReadByte() (byte, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.b.ReadByte()
}

// ReadBytes calls "ReadBytes" of "bytes.Buffer".
func (b *SynchronizedBuffer) ReadBytes(delim byte) (line []byte, err error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.b.ReadBytes(delim)
}

// ReadFrom calls "ReadFrom" of "bytes.Buffer".
func (b *SynchronizedBuffer) ReadFrom(r io.Reader) (n int64, err error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.b.ReadFrom(r)
}

// ReadRune calls "ReadRune" of "bytes.Buffer".
func (b *SynchronizedBuffer) ReadRune() (r rune, size int, err error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.b.ReadRune()
}

// ReadString calls "ReadString" of "bytes.Buffer".
func (b *SynchronizedBuffer) ReadString(delim byte) (line string, err error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.b.ReadString(delim)
}

// Reset calls "Reset" of "bytes.Buffer".
func (b *SynchronizedBuffer) Reset() {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.b.Reset()
}

// String calls "String" of "bytes.Buffer".
func (b *SynchronizedBuffer) String() string {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.b.String()
}

// Truncate calls "Truncate" of "bytes.Buffer".
func (b *SynchronizedBuffer) Truncate(n int) {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.b.Truncate(n)
}

// UnreadByte calls "UnreadByte" of "bytes.Buffer".
func (b *SynchronizedBuffer) UnreadByte() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.b.UnreadByte()
}

// UnreadRune calls "UnreadRune" of "bytes.Buffer".
func (b *SynchronizedBuffer) UnreadRune() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.b.UnreadRune()
}

// Write calls "Write" of "bytes.Buffer".
func (b *SynchronizedBuffer) Write(p []byte) (n int, err error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.b.Write(p)
}

// WriteByte calls "WriteByte" of "bytes.Buffer".
func (b *SynchronizedBuffer) WriteByte(c byte) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.b.WriteByte(c)
}

// WriteRune calls "WriteRune" of "bytes.Buffer".
func (b *SynchronizedBuffer) WriteRune(r rune) (n int, err error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.b.WriteRune(r)
}

// WriteString calls "WriteString" of "bytes.Buffer".
func (b *SynchronizedBuffer) WriteString(s string) (n int, err error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.b.WriteString(s)
}

// WriteTo calls "WriteTo" of "bytes.Buffer".
func (b *SynchronizedBuffer) WriteTo(w io.Writer) (n int64, err error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.b.WriteTo(w)
}
