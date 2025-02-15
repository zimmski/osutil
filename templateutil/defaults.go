package templateutil

import (
	"encoding/base64"
	"text/template"
	"unicode"

	"github.com/symflower/pretty"
	"github.com/zimmski/osutil/bytesutil"
)

// DefaultFuncMap holds common template functions.
var DefaultFuncMap = template.FuncMap{
	"base64": func(in string) string {
		return base64.StdEncoding.EncodeToString([]byte(in))
	},
	"prefixContinuationLinesWith": bytesutil.PrefixContinuationLinesWith,
	"lowerFirst": func(s string) string {
		return string(unicode.ToLower(rune(s[0]))) + s[1:]
	},
	"pretty": func(data any) string {
		return pretty.Sprintf("%# v", data)
	},
	"prettyLazy": func(data any) string {
		return pretty.LazySprintf("%# v", data)
	},
	"quote": func(data any) string {
		return pretty.Sprintf("%q", data)
	},
}

// MergeFuncMaps returns all functions of "a" and all functions of "b" in a new function mapping.
// For entries that are defined in both maps the entry defined in b is chosen.
func MergeFuncMaps(a template.FuncMap, b template.FuncMap) template.FuncMap {
	c := template.FuncMap{}

	for n, f := range a {
		c[n] = f
	}
	for n, f := range b {
		c[n] = f
	}

	return c
}
