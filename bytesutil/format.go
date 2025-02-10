package bytesutil

import (
	"fmt"

	"github.com/symflower/pretty"
)

// FormatToGoObject formats the given object to equivalent Go code.
func FormatToGoObject(object any) string {
	return fmt.Sprintf("%# v", pretty.LazyFormatter(object))
}
