package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ridwanulhoquejr/lets-go-further/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create new movie handler")

	// body := httprouter.ParamsFromContext(r.Body())
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil || id <= 0 {
		http.NotFound(w, r)
		return
	}

	// Create a new instance of the Movie struct, containing the ID we extracted from
	// the URL and some dummy data. Also notice that we deliberately haven't set a
	// value for the Year field.

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Breaking Bad",
		Runtime:   100,
		Genres:    []string{"War", "Drug", "Fight"},
		// Year:      2012,
		Version: 10,
	}

	err = app.writeJSON(w, http.StatusOK, movie, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Some error", http.StatusInternalServerError)
	}

}
