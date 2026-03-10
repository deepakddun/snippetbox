package main

import (
	"html/template"
	"path/filepath"

	"github.com/deepakddun/snippetbox/internal/models"
)

type templateData struct {
	Snippet         models.Snippet
	Snippets        []models.Snippet
	CurrentYear     int
	Form            any
	Flash           string
	IsAuthenticated bool
}

func newTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./../../ui/html/pages/*.tmpl")

	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		name := filepath.Base(page)
		// content, err := os.ReadFile(page)
		// if err != nil {
		// 	return nil, err
		// }
		// fmt.Println("PAGE:", page)
		// fmt.Println(string(content))

		files := []string{
			"./../../ui/html/pages/base.tmpl",
			"./../../ui/html/pages/nav.tmpl",
			page,
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts

		// fmt.Println("cache key:", name)
		// fmt.Println("files parsed:", files)
		// fmt.Println(ts.DefinedTemplates())
		// fmt.Println("-----")

	}

	return cache, nil
}
