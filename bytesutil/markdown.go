package bytesutil

import (
	"bytes"
	"os"

	pkgerrors "github.com/pkg/errors"
	"github.com/yuin/goldmark"
)

// RenderMarkdownFileToHTMLFile reads in the Markdown file, renders it as HTML and writes that output in the HTML file.
func RenderMarkdownFileToHTMLFile(markdownFilePath string, htmlFilePath string) (err error) {
	data, err := os.ReadFile(markdownFilePath)
	if err != nil {
		return pkgerrors.Wrap(err, markdownFilePath)
	}

	var html bytes.Buffer
	if err := goldmark.Convert(data, &html); err != nil {
		return pkgerrors.Wrap(err, markdownFilePath)
	}

	if err := os.WriteFile(htmlFilePath, html.Bytes(), 0640); err != nil {
		return pkgerrors.Wrap(err, html.String())
	}

	return nil
}
