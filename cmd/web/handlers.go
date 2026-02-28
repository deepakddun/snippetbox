package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/deepakddun/snippetbox/internal/models"
)

type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest(r.Context())

	if err != nil {

		//app.logger.Error(err.Error())
		app.serverError(w, r, err)
		return
	}

	data := templateData{
		Snippets:    snippets,
		CurrentYear: app.getCurrentYear(),
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

	data := templateData{

		Snippet:     snippet,
		CurrentYear: app.getCurrentYear(),
	}

	app.render(w, r, "view.tmpl", http.StatusOK, data)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	data := templateData{

		Snippet:     models.Snippet{},
		CurrentYear: app.getCurrentYear(),
	}
	data.Form = snippetCreateForm{
		Expires: 7,
	}
	app.render(w, r, "create.tmpl", http.StatusOK, data)

}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	// 	w.Header().Add("Server", "Go")
	// 	//w.WriteHeader(http.StatusCreated)
	// 	//w.WriteHeader(200)

	// 	title := "O snail"
	// 	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	// 	expires := 7
	// 	id, err := app.snippets.Insert(r.Context(), title, content, expires)

	// 	app.logger.Info("Created Record With ID", "Id ", id)

	// 	if err != nil {
	// 		app.serverError(w, r, err)
	// 		return
	// 	}
	// 	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	// 	//fmt.Fprintf(w, "Creating SnippetBox .Post..")
	// }

	err := r.ParseForm()

	if err != nil {

		app.ClientError(w, r, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))

	if err != nil {

		app.ClientError(w, r, http.StatusBadRequest)
		return
	}

	// Check for errors

	errorMap := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		errorMap["title"] = "title cannot be empty"
	}

	if utf8.RuneCountInString(title) > 100 {
		errorMap["title"] = "title length cannot be greater than 100"
	}

	if strings.TrimSpace(content) == "" {
		errorMap["content"] = "Content cannot be empty"
	}

	if expires != 1 && expires != 7 && expires != 365 {
		errorMap["expires"] = "Expires can only be 1 , 7 or 365 Days"
	}
	data := templateData{

		CurrentYear: app.getCurrentYear(),
		Form: snippetCreateForm{
			Title:       title,
			Content:     content,
			FieldErrors: errorMap,
			Expires:     expires,
		},
	}
	if len(errorMap) > 0 {
		app.render(w, r, "create.tmpl", http.StatusUnprocessableEntity, data)
		return
	}

	id, err := app.snippets.Insert(r.Context(), title, content, expires)

	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}
