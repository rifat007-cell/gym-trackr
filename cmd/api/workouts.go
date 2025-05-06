package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/tanvir-rifat007/gymBuddy/internal/data"
)

// func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
// 	// Simulate an error (you can replace this with a real one later)
// 	err := errors.New("simulated health check failure")

// 	// Get the Sentry hub from the request context
// 	if hub := sentry.GetHubFromContext(r.Context()); hub != nil {
// 		hub.WithScope(func(scope *sentry.Scope) {
// 			scope.SetTag("handler", "healthCheck")
// 			scope.SetLevel(sentry.LevelError)
// 			hub.CaptureException(err)
// 		})
// 	}

// 	// Respond with error to client
// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// }



func (app *application) catchAllClientRequestHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the client application
	http.ServeFile(w,r,"./public/index.html")


}




func (app *application) getExercisesByWorkoutHandler(w http.ResponseWriter, r *http.Request) {
	var input struct{
		Goal string `json:"goal"`
		Label string `json:"label"`
	}


	err:= app.readJSON(w,r,&input)

	if err!=nil{
		app.badRequestResponse(w,r,err)
		return
	}

	fmt.Println("Decoded JSON:", input)

	// get from db

	workouts,err:=app.models.Workouts.GetAllExerciseBasedWorkoutName(input.Goal,input.Label)

	if err!=nil{
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}


	// get the user email from context
	email,ok := r.Context().Value("email").(string)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("missing email in context"))
		return
	}

	// get the user from db:

	user,err:= app.models.Users.GetUserByEmail(email)

	if err!=nil{
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// get the user id from the user
	userId:= user.ID

	// save the user daily plain
	_,err = app.models.UserDailyPlan.InsertDailyPlan(userId,workouts[0].ID)




	if err!=nil{
		app.serverErrorResponse(w,r,err)
		return
	}

// make a map of workouts
	data:= envelope{
		"workouts": workouts,
	}

	err = app.writeJSON(w, http.StatusOK, data, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}


	
	
}
