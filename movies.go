package main

import (
	"fmt"
	"net/http"
)

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get the movie by the id")

	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
	}
	fmt.Fprintf(w, "show the details of the movie %d\n", id)
}
