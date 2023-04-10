package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/onahvictor/Snippet/internal/models"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	// CSRFToken       string
}

func humanDate(t time.Time) string {
	if t.IsZero(){
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	//this returns a slice of all filepaths that ends with .html
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		//extracts the file name from the full path
		name := filepath.Base(page)

		temp := template.New(name).Funcs(functions)
		temp, err := temp.ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		//parse the files into a templare set.
		temp, err = temp.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		temp, err = temp.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = temp
	}

	return cache, nil

}
