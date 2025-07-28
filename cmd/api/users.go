package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/tanvir-rifat007/gymBuddy/internal/data"
	"github.com/tanvir-rifat007/gymBuddy/token"
	"github.com/tanvir-rifat007/gymBuddy/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	// password hashing
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		if errors.Is(err, data.ErrDuplicateEmail) {
			v.AddError("email", "email already exists")
			app.failedValidationResponse(w, r, v.Errors)
			return
		}
		app.serverErrorResponse(w, r, err)
		return

	}

	// after creating user send a jwt token to the client:
	user.JWT = token.CreateJWT(*user, app.logger)

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// and then send a token to the user gmail to activate the account

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	type ActivationData struct {
		User  *data.User  `json:"user"`
		Token *data.Token `json:"token"`
	}

	activationData := ActivationData{
		User:  user,
		Token: token,
	}

	emails := []string{
		user.Email,
	}

	err = app.sendEmail(emails, "user activation", "./internal/mailer/templates/user_welcome.tmpl.html", activationData)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := &data.User{
		Email: input.Email,
		Password: data.Password{
			Plaintext: &input.Password,
		},
	}
	user, err = app.models.Users.Login(user.Email, *user.Password.Plaintext)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrInvalidCredentials):
			v := validator.New()
			v.AddError("email", "invalid credentials")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// after login user send a jwt token to the client:
	user.JWT = token.CreateJWT(*user, app.logger)
	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	tokenPlainText := r.URL.Query().Get("token")

	// Validate the plaintext token provided by the client.
	v := validator.New()

	if data.ValidateTokenPlaintext(v, tokenPlainText); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetUserFromToken(data.ScopeActivation, tokenPlainText)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Update the user's activation status.
	user.Activated = true

	// Save the updated user record in our database, checking for any edit conflicts in
	// the same way that we did for our movie records.
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// If everything went successfully, then we delete all activation tokens for the
	// user.
	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// 	// Send a success message to the client.
	// 	// Simple HTML response.
	// w.Header().Set("Content-Type", "text/html")
	// w.WriteHeader(http.StatusOK)
	// fmt.Fprintf(w, `<h1>Your account has been successfully activated!</h1><p>You can now log in.</p>`)

	// then send a jwt token to the client with activated status

	user.JWT = token.CreateJWT(*user, app.logger)
	// Redirect user to frontend with JWT in query string
	redirectURL := fmt.Sprintf("https://gym-trackr-production.up.railway.app/activated?token=%s", user.JWT)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)

}
