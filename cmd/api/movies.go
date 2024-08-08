package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ridwanulhoquejr/lets-go-further/internal/data"
)

func (app *application) createMovieHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	fmt.Println("Create new movie handler")

	var input struct {
		ID        int       `json:"id"`
		CreatedAt time.Time `json:"created_at"` // - tag will hide this field in respone object
		Title     string    `json:"title"`
		Runtime   int       `json:"runtime"`
		Genres    []string  `json:"genres"`
		Year      int       `json:"year,omitempty,string"`
	}

	// use our readJSON helper method
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// fmt.Fprintf(w, "%+v\n", input)
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": input}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) showMovieHandler(
	w http.ResponseWriter, r *http.Request,
) {

	id, err := app.readIDParam(r)
	if err != nil || id <= 0 {
		app.notFoundResponse(w, r)
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
		Year:      2012,
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w, r, err)
	}
}
