package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) readIDParam(r *http.Request) (int64, error) {

	// get the params from the request!
	params := httprouter.ParamsFromContext(r.Context())

	// parse the id param into int64 from the string
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)

	if err != nil {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}
