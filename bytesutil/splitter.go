package bytesutil

import (
	"bytes"
)

// Split iteratively splits argument s and returns the split items over the returned channel.
func Split(s []byte, sep byte) <-chan []byte {
	ch := make(chan []byte)

	go func() {
		for {
			end := bytes.IndexByte(s, sep)
			if end == -1 {
				end = len(s)
			}

			ch <- s[:end]

			if end == len(s) {
				break
			}

			s = s[end+1:]
		}

		close(ch)
	}()

	return ch
}
