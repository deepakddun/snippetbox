package main

import (
	"html/template"
	"path/filepath"

	"github.com/deepakddun/snippetbox/internal/models"
)

type templateData struct {
	Snippet     models.Snippet
	Snippets    []models.Snippet
	CurrentYear int
	Form        any
}

func newTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./../../ui/html/pages/*.tmpl")

	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		name := filepath.Base(page)
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

	}

	return cache, nil
}
