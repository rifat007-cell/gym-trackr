package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/tanvir-rifat007/gymBuddy/internal/data"
)


func (app *application) workoutLogHandler(w http.ResponseWriter, r *http.Request){
	var input struct{
		Exercise string `json:"exercise"`
		Sets     string    `json:"sets"`
		Reps     string    `json:"reps"`
		Duration string    `json:"duration"`
		Weight   string    `json:"weight"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Extract the email from the context
	// that i added in the middleware
	 email,ok := r.Context().Value("email").(string)

	if !ok {
		app.serverErrorResponse(w, r, errors.New("missing email in context"))
		return
	}

	// Get the user ID from the database using the email
	user, err := app.models.Users.GetUserByEmail(email)

	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	 sets,err:=strconv.Atoi(input.Sets)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	reps,err:=strconv.Atoi(input.Reps)


	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	duration,err:=strconv.Atoi(input.Duration)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	  weight,err:= strconv.Atoi(input.Weight)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	workoutLog := &data.WorkoutLog{
		UserID:   user.ID,
		Exercise: input.Exercise,
		Sets:     sets,
		Reps:     reps,
		Duration: duration,
		Weight:  weight,
	}

	if err := app.models.WorkoutLogs.Insert(workoutLog); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"workout_log": workoutLog}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}



}


func (app *application) getWorkoutVolumeHandler(w http.ResponseWriter,r *http.Request){
	// Extract the email from the context
	// that i added in the middleware
	email, ok := r.Context().Value("email").(string)

	if !ok {
		app.serverErrorResponse(w, r, errors.New("missing email in context"))
		return
	}

	// Get the user ID from the database using the email
	user, err := app.models.Users.GetUserByEmail(email)

	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	volume, err := app.models.WorkoutLogs.GetVolumeOverTime(user.ID)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"volumes": volume,
	}

	err = app.writeJSON(w, http.StatusOK, data, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}