package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(
	_ *http.Request,
	err error,
) {
	app.logger.Println(err)
}

func (app *application) errorResponse(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	message any,
) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}

// Note that the errors parameter here has the type map[string]string, which is exactly
// the same as the errors map contained in our Validator type.
func (app *application) failedValidationResponse(
	w http.ResponseWriter,
	r *http.Request,
	errors map[string]string,
) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) badRequestResponse(
	w http.ResponseWriter,
	r *http.Request, err error,
) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) serverErrorResponse(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {

	app.logError(r, err)

	msg := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, msg)

}

func (app *application) notFoundResponse(
	w http.ResponseWriter,
	r *http.Request,
) {

	msg := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, msg)
}

func (app *application) methodNotAllowedResponse(
	w http.ResponseWriter,
	r *http.Request,
) {

	msg := fmt.Sprintf("the %s mehtod is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, msg)
}
