package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {

	var method = r.Method
	var uri = r.URL.RequestURI()

	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) ClientError(w http.ResponseWriter, r *http.Request, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, page string, statusCode int, data templateData) {

	cache := app.templateCache
	fmt.Println(cache)
	fmt.Println(cache)
	buf := new(bytes.Buffer)
	ts, ok := cache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exists", page)
		app.serverError(w, r, err)
		return

	}

	fmt.Println("render page:", page)
	fmt.Println(ts.DefinedTemplates())
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	w.WriteHeader(statusCode)
	buf.WriteTo(w)

}

func (app *application) getCurrentYear() int {
	return time.Now().Year()
}

func (app *application) isAutheticated(r *http.Request) bool {

	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{

		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		CurrentYear:     app.getCurrentYear(),
		IsAuthenticated: app.isAutheticated(r),
	}
}
