package main

import (
	"net/http"
	"time"
)

func (app *application) healthCheckHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	// Create a map which holds the information that we want to send in the response.
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	app.logger.Println("Starting long-running request")
	time.Sleep(3 * time.Second)
	app.logger.Println("Completed long-running request after 4 sec")

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w, r, err)
	}
}
