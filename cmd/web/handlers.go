package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {

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

		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tp.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	//w.WriteHeader()
}

func snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Viewing SnippetBox ...%d", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Creating SnippetBox ...")
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Server", "Go")
	w.WriteHeader(http.StatusCreated)
	//w.WriteHeader(200)

	fmt.Fprintf(w, "Creating SnippetBox .Post..")
}
