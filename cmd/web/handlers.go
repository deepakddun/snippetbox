package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Server", "Go")

	files := []string{
		"./../../ui/html/pages/base.tmpl",
		"./../../ui/html/pages/home.tmpl",
		"./../../ui/html/pages/nav.tmpl",
	}

	tp, err := template.ParseFiles(files...)
	// if tp.Lookup("base") == nil {
	// 	log.Println(`Lookup("base") = nil (NOT defined)`)
	// } else {
	// 	log.Println(`Lookup("base") OK`)
	// }

	// for _, t := range tp.Templates() {
	// 	log.Println("template:", t.Name())
	// }

	// for _, f := range files {
	// 	abs, _ := filepath.Abs(f)
	// 	log.Println("using file:", abs)
	// }

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = tp.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	//w.WriteHeader()
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		app.logger.Error(err.Error())
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Viewing SnippetBox ...%d", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Creating SnippetBox ...")
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Server", "Go")
	w.WriteHeader(http.StatusCreated)
	//w.WriteHeader(200)

	fmt.Fprintf(w, "Creating SnippetBox .Post..")
}
