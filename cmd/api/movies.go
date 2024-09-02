package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/ridwanulhoquejr/lets-go-further/internal/data"
	"github.com/ridwanulhoquejr/lets-go-further/internal/validator"
)

// get muliple movie handler
func (app *application) listMovieHandler(w http.ResponseWriter, r *http.Request) {

	// To keep things consistent with our other handlers, we'll define an input struct
	// to hold the expected values from the request query string.
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}

	// instantiate the validator
	v := validator.New()

	// Call r.URL.Query() to get the url.Values map containing the query string data.
	qs := r.URL.Query()

	// Use our helpers to extract the title and genres query string values, falling back
	// to defaults of an empty string and an empty slice respectively if they are not
	// provided by the client.
	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})

	// Get the page and page_size query string values as integers. Notice that we set
	// the default page value to 1 and default page_size to 20, and that we pass the
	// validator instance as the final argument here.
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 10, v)

	// Extract the sort query string value, falling back to "id" if it is not provided
	// by the client (which will imply a ascending sort on movie ID).
	input.Filters.Sort = app.readString(qs, "sort", "id")

	// sort safe list
	input.Filters.SortSafelist = []string{"title", "id", "runtime", "year", "-title", "-id", "-runtime", "-year"}

	// Check the Validator instance for any errors and use the failedValidationResponse()
	// helper to send the client a response if necessary.
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Dump the contents of the input struct in a HTTP response.
	// fmt.Fprintf(w, "\n%+v\n", input)
	movies, metadata, err := app.models.Movie.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movies, "metadata": metadata}, nil)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

}

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
		Title   *string  `json:"title"`
		Runtime *int32   `json:"runtime"`
		Genres  []string `json:"genres"`
		Year    *int32   `json:"year,omitempty"`
	}
	// app.logger.Printf("They year before readJSON is %d", *input.Year)

	// decode the payload using our readJSON helper
	err = app.readJSON(w, r, &input)

	// Copy the values of decoded payload field in movie fetch by Get(id)

	// for partial updates, we will check wheather provided json is nil or not,
	// as this is now a pointer type this will be nil if not provider
	if input.Title != nil {
		movie.Title = *input.Title
	}

	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}

	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	if input.Year != nil {
		movie.Year = *input.Year
	}

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
