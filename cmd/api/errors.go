package main

import (
	"fmt"
	"net/http"
)

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		var (
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Error(err.Error(), "method", method, "uri", uri)
		w.WriteHeader(500)
	}
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request) {
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) dateTooEarlyResponse(w http.ResponseWriter, r *http.Request) {
	message := "date is before the adoption of the republican calendar"
	app.errorResponse(w, r, http.StatusBadRequest, message)
}

func (app *application) dateTooHighResponse(w http.ResponseWriter, r *http.Request) {
	message := "date must be 9999-12-31 or lower"
	app.errorResponse(w, r, http.StatusBadRequest, message)
}

func (app *application) invalidDateResponse(w http.ResponseWriter, r *http.Request) {
	message := "date must be in YYYY-MM-DD format"
	app.errorResponse(w, r, http.StatusBadRequest, message)
}
