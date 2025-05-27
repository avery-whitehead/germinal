package main

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rickb777/date"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/", app.todayHandler)
	router.HandlerFunc(http.MethodGet, "/:date", app.dateHandler)

	return router
}

func (app *application) todayHandler(w http.ResponseWriter, r *http.Request) {
	date, err := app.toRepublican(date.TodayUTC())
	if err != nil {
		app.serverErrorResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"date": date}, nil)
	if err != nil {
		app.serverErrorResponse(w, r)
	}
}

func (app *application) dateHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	dateParam, err := date.ParseISO(params.ByName("date"))
	if err != nil {
		app.invalidDateResponse(w, r)
		return
	}

	date, err := app.toRepublican(dateParam)
	if err != nil {
		switch {
		case errors.Is(err, ErrBeforeCalendar):
			app.dateTooEarlyResponse(w, r)
		case errors.Is(err, ErrDateTooHigh):
			app.dateTooHighResponse(w, r)
		default:
			app.serverErrorResponse(w, r)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"date": date}, nil)
}
