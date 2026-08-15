package server

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// Эту структуру не трогаем, она тут из прошлого урока...
// type templateData struct {
// 	Username string
// 	Name     string
// 	// AvatarURL string
// }

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("path %s was not found", dir)
	}

	// Собираем все layout и partial файлы
	var layouts, partials []string

	// Функция для сбора вспомогательных файлов
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

	// Собираем все вспомогательные файлы
	if err := filepath.Walk(dir, collectHelpers); err != nil {
		return nil, err
	}

	// Функция для обработки page шаблонов
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

			// Добавляем все layout файлы
			if len(layouts) > 0 {
				ts, err = ts.ParseFiles(layouts...)
				if err != nil {
					return err
				}
			}

			// Добавляем все partial файлы
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

	// Обрабатываем page шаблоны
	if err := filepath.Walk(dir, processPage); err != nil {
		return nil, err
	}

	return cache, nil
}
