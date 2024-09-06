package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// we separating the routes from the main.go file
// instead creating the mux := http.NewServeMux from the default net/http
// we use httprouter's -> *httprouter.Router

// then we can replace the Handler value from the &Server in maing.go to app.routes() method
func (app *application) routes() http.Handler {

	r := httprouter.New()

	// we use our error helper mthod to override the built-in `NotFound` & `methodNotAllowed` error responses.
	r.NotFound = http.HandlerFunc(app.notFoundResponse)
	r.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	r.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)
	r.HandlerFunc(http.MethodPost, "/v1/movie", app.createMovieHandler)
	r.HandlerFunc(http.MethodGet, "/v1/movie", app.listMovieHandler)
	r.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.showMovieHandler)
	r.HandlerFunc(http.MethodPatch, "/v1/movie/:id", app.updateMovieHandler)
	r.HandlerFunc(http.MethodDelete, "/v1/movie/:id", app.deleteMovieHandler)

	// User route handler
	r.HandlerFunc(http.MethodPost, "/v1/users", app.createUserHandler)

	return app.recoverPanic(r)
}
