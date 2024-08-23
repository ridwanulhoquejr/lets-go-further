package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

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
		Runtime int32    `json:"runtime"`
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
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movie.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			app.logger.Printf("Error in movies handler: %s", err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	// extract id for r.body
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// get the movie from Get methtod using id that extracted above
	movie, err := app.models.Movie.Get(id)
	if err != nil {
		switch {
		case errors.Is(sql.ErrNoRows, err):
			app.notFoundResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	// create payload input
	var input struct {
		Title   string   `json:"title"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
		Year    int32    `json:"year,omitempty"`
	}

	// decode the payload using our readJSON helper
	err = app.readJSON(w, r, &input)

	// Copy the values of decoded payload field in movie fetch by Get(id)
	movie.Title = input.Title
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres
	movie.Year = input.Year

	// Initialize a new Validator instance.
	v := validator.New()

	// Call the ValidateMovie() function and return a response containing the errors if
	// any of the checks fail.
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// now we can proceed with actual update task
	err = app.models.Movie.Update(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// header for product path url
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	// convert our response to json by using writeJSON helper method
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {

	// extract the id
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Movie.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"messege": "movie succesfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
