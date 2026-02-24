package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/deepakddun/snippetbox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Server", "Go")

	snippets, err := app.snippets.Latest(r.Context())

	if err != nil {

		//app.logger.Error(err.Error())
		app.serverError(w, r, err)
		return
	}

	// files := []string{
	// 	"./../../ui/html/pages/base.tmpl",
	// 	"./../../ui/html/pages/home.tmpl",
	// 	"./../../ui/html/pages/nav.tmpl",
	// }

	// ts, err := template.ParseFiles(files...)

	// if err != nil {

	// 	//app.logger.Error(err.Error())
	// 	app.serverError(w, r, err)
	// 	return
	// }

	data := templateData{
		Snippets: snippets,
	}

	app.render(w, r, "home.tmpl", http.StatusOK, data)
	// err = ts.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	//json.NewEncoder(w).Encode(snippets)

	// files := []string{
	// 	"./../../ui/html/pages/base.tmpl",
	// 	"./../../ui/html/pages/home.tmpl",
	// 	"./../../ui/html/pages/nav.tmpl",
	// }

	// tp, err := template.ParseFiles(files...)
	// // if tp.Lookup("base") == nil {
	// // 	log.Println(`Lookup("base") = nil (NOT defined)`)
	// // } else {
	// // 	log.Println(`Lookup("base") OK`)
	// // }

	// // for _, t := range tp.Templates() {
	// // 	log.Println("template:", t.Name())
	// // }

	// // for _, f := range files {
	// // 	abs, _ := filepath.Abs(f)
	// // 	log.Println("using file:", abs)
	// // }

	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	// err = tp.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }
	//w.WriteHeader()
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Starting view ...")
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil || id < 1 {
		app.serverError(w, r, err)
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(r.Context(), id)

	if err != nil {
		if errors.Is(err, models.ErrorNoRecord) {

			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)

		}
		return
	}

	fmt.Printf("%q\n", snippet.Content)

	files := []string{
		"./../../ui/html/pages/base.tmpl", "./../../ui/html/pages/nav.tmpl", "./../../ui/html/pages/view.tmpl",
	}

	ts, err := template.ParseFiles(files...)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := templateData{
		Snippet: snippet,
	}
	err = ts.ExecuteTemplate(w, "base", data)

	if err != nil {
		app.serverError(w, r, err)
	}

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Creating SnippetBox ...")
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Server", "Go")
	//w.WriteHeader(http.StatusCreated)
	//w.WriteHeader(200)

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7
	id, err := app.snippets.Insert(r.Context(), title, content, expires)

	app.logger.Info("Created Record With ID", "Id ", id)

	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	//fmt.Fprintf(w, "Creating SnippetBox .Post..")
}
