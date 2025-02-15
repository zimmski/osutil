package templateutil

import "io"

type Template interface {
	// Execute applies the template to the specified data object and writes the output to the writer.
	Execute(wr io.Writer, data any) (err error)
}
