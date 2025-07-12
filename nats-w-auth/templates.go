package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"strings"
)

func LoadTemplates(fsys fs.FS) (*template.Template, error) {
	root := template.New("")

	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, ".gohtml") {
			content, err := fs.ReadFile(fsys, path)

			if err != nil {
				return fmt.Errorf("reading template %s: %w", path, err)
			}

			// Use file path as template name
			templateName, _ := strings.CutSuffix(path, ".gohtml")
			tmpl := root.New(templateName)
			_, err = tmpl.Parse(string(content))
			if err != nil {
				return fmt.Errorf("parsing template %s: %w", path, err)
			}
		}
		return nil
	})

	return root, err
}
