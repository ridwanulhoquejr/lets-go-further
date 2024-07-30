package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// we separating the routes from the main.go file
// instead creating the mux := http.NewServeMux from the default net/http
// we use httprouter's -> *httprouter.Router

// then we can replace the Handler value from the &Server in maing.go to app.routes() method
func (app *application) routes() *httprouter.Router {

	r := httprouter.New()

	r.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)
	r.HandlerFunc(http.MethodGet, "/v1/movies/get/:id", app.showMovieHandler)
	r.HandlerFunc(http.MethodPost, "/v1/movies/create", app.createMovieHandler)

	return r
}
