package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

//! writeJSON helper method for our application
//
// Define a writeJSON() helper for sending responses. This takes the destination
// http.ResponseWriter, the HTTP status code to send, the data to encode to JSON, and a
// header map containing any additional HTTP headers we want to include in the response.

func (app *application) writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
	headers http.Header) error {

	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// At this point, we know that we won't encounter any more errors before writing the
	// response, so it's safe to add any headers that we want to include. We loop
	// through the header map and add each header to the http.ResponseWriter header map.
	// Note that it's OK if the provided header map is nil. Go doesn't throw an error
	// if you try to range over (or generally, read from) a nil map.
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "applicaiton/json")
	// for sending the status code via header
	w.WriteHeader(status)
	w.Write([]byte(js))

	return nil
}

func (app *application) readIDParam(r *http.Request) (int, error) {

	// get the params from the request!
	params := httprouter.ParamsFromContext(r.Context())

	// parse the id param into int64 from the string
	// id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	// ParseInt in same as Atoi but with just int
	id, err := strconv.Atoi(params.ByName("id"))

	if err != nil {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}
