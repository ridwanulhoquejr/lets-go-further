package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ridwanulhoquejr/lets-go-further/internal/data"
	"github.com/ridwanulhoquejr/lets-go-further/internal/validator"
)

func (app *application) createMovieHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	fmt.Println("Create new movie handler")

	var input struct {
		Title   string   `json:"title"`
		Runtime int      `json:"runtime"`
		Genres  []string `json:"genres"`
		Year    int32    `json:"year,omitempty"`
	}

	// use our readJSON helper method
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the values from the input struct to a new Movie struct.
	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	// Initialize a new Validator instance.
	v := validator.New()

	// Call the ValidateMovie() function and return a response containing the errors if
	// any of the checks fail.
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Call the Insert() method on our movies model, passing in a pointer to the
	// validated movie struct. This will create a record in the database and update the
	// movie struct with the system-generated information.
	err = app.models.Movie.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// When sending a HTTP response, we want to include a Location header to let the
	// client know which URL they can find the newly-created resource at. We make an
	// empty http.Header map and then use the Set() method to add a new Location header,
	// interpolating the system-generated ID for our new movie in the URL.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	// fmt.Fprintf(w, "%+v\n", input)
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, headers)
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
