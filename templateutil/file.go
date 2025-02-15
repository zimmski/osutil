package templateutil

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"text/template"

	pkgerrors "github.com/pkg/errors"
	"github.com/symflower/pretty"
)

// WriteTemplateToFile executes a template with the given data and saves the result into a file.
func WriteTemplateToFile(filePath string, tmpl *template.Template, data any) error {
	var driver bytes.Buffer

	err := tmpl.Execute(&driver, data)
	if err != nil {
		return pkgerrors.Wrap(err, pretty.LazySprintf("%# v", data))
	}

	err = os.WriteFile(filePath, driver.Bytes(), 0640)
	if err != nil {
		return pkgerrors.Wrap(err, filePath)
	}

	return nil
}

// RewriteFileAsTemplate read in a file, execute it as a template with the given data and save the result into the same file.
func RewriteFileAsTemplate(filePath string, funcMap template.FuncMap, data any) (err error) {
	tmpl, err := template.New(filepath.Base(filePath)).Funcs(funcMap).ParseFiles(filePath) // REMARK Use the file name as template identifier because otherwise `template.ParseFiles` fails.
	if err != nil {
		return pkgerrors.Wrap(err, filePath)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return pkgerrors.Wrap(err, filePath)
	}
	defer func() {
		if e := f.Close(); e != nil {
			err = errors.Join(err, e)
		}
	}()

	if err := tmpl.Execute(f, data); err != nil {
		return pkgerrors.Wrap(err, filePath)
	}

	return nil
}
