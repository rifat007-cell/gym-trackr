package main

import (
	"errors"
	"net/http"

	"github.com/tanvir-rifat007/gymBuddy/internal/data"
)


func (app *application) getMealByWorkoutHandler(w http.ResponseWriter, r *http.Request){
	var input struct{
		Goal string `json:"goal"`
		DietaryPreference string `json:"dietary_preference"`
	}

	err:= app.readJSON(w,r,&input)

	if err!=nil{
		app.badRequestResponse(w,r,err)
		return
	}

	meals,err:=app.models.Meals.GetAllMealByWorkoutName(input.Goal,input.DietaryPreference)

	if err!=nil{
		app.serverErrorResponse(w,r,err)
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
	_,err = app.models.UserDailyPlan.UpdateDailyPlan(meals[0].ID,userId)

	if err!=nil{
		app.serverErrorResponse(w,r,err)
		return
	}

	data:= envelope{"meals": meals}

	err = app.writeJSON(w, http.StatusOK, data, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}




}