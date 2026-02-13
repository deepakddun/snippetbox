package main

import (
	"net/http"
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
