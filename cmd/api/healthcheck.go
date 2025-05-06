package main

import "net/http"


func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data:= envelope{
		"status": "available",
		"sys_info": map[string]any{
			"environment": "development",
       "version": "1.0.0",

		},
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Error("Error writing JSON response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}