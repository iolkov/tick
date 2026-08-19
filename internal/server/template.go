package server

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, err
	}

	var layouts, partials []string

	collectHelpers := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".html" {
			base := filepath.Base(path)
			if strings.HasSuffix(base, ".layout.html") {
				layouts = append(layouts, path)
			} else if strings.HasSuffix(base, ".partial.html") {
				partials = append(partials, path)
			}
		}
		return nil
	}

	if err := filepath.Walk(dir, collectHelpers); err != nil {
		return nil, err
	}

	processPage := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".html" && strings.HasSuffix(filepath.Base(path), ".page.html") {
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}

			name := filepath.ToSlash(relPath)
			ts, err := template.ParseFiles(path)
			if err != nil {
				return err
			}

			if len(layouts) > 0 {
				ts, err = ts.ParseFiles(layouts...)
				if err != nil {
					return err
				}
			}

			if len(partials) > 0 {
				ts, err = ts.ParseFiles(partials...)
				if err != nil {
					return err
				}
			}

			cache[name] = ts
		}
		return nil
	}

	if err := filepath.Walk(dir, processPage); err != nil {
		return nil, err
	}

	return cache, nil
}
