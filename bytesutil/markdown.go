package bytesutil

import (
	"bytes"
	"os"

	"github.com/yuin/goldmark"
)

// RenderMarkdownFileToHTMLFile reads in the Markdown file, renders it as HTML and writes that output in the HTML file.
func RenderMarkdownFileToHTMLFile(markdownFilePath string, htmlFilePath string) (err error) {
	data, err := os.ReadFile(markdownFilePath)
	if err != nil {
		return err
	}

	var html bytes.Buffer
	if err := goldmark.Convert(data, &html); err != nil {
		return err
	}

	if err := os.WriteFile(htmlFilePath, html.Bytes(), 0640); err != nil {
		return err
	}

	return nil
}
