package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/tanvir-rifat007/gymBuddy/token"
)

// type WebAuthnHandler struct {
// 	storage  data.PasskeyStore
// 	logger   *logger.Logger
// 	webauthn *webauthn.WebAuthn
// }

// func NewWebAuthnHandler(storage data.PasskeyStore, logger *logger.Logger, webauthn *webauthn.WebAuthn) *WebAuthnHandler {
// 	return &WebAuthnHandler{
// 		storage:  storage,
// 		logger:   logger,
// 		webauthn: webauthn,
// 	}
// }

func (h *application) writeJSONResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode response", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return err
	}
	return nil
}

func (h *application) WebAuthnRegistrationBeginHandler(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value("email").(string)

	fmt.Println("Email from context:", email)
	if !ok {
		// h.logger.Error("Unable to retrieve email", nil)
		http.Error(w, "Unable to retrieve email", http.StatusInternalServerError)
		return
	}
	user, err := h.models.Passkey.GetUserByEmail(email)
	fmt.Println("user",user)
	if err != nil {
		h.logger.Error("Failed to find user", "err",err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	options, session, err := h.webauthn.BeginRegistration(user)
	if err != nil {
		h.logger.Error("Unable to retrieve email", "err", err)
		http.Error(w, "Can't begin WebAuthn Registration", http.StatusInternalServerError)

		return
	}

	// Make a session key and store the sessionData values
	t, err := h.models.Passkey.GenSessionID()
	if err != nil {
		h.logger.Error("Can't generate session id: %s", "err",err)
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	h.models.Passkey.SaveSession(t, *session)

	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    t,
		Path:     "api/passkey/registerStart",
		MaxAge:   3600,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	h.writeJSONResponse(w, options)
}

func (h *application) WebAuthnRegistrationEndHandler(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value("email").(string)
	if !ok {
		h.logger.Error("Unable to retrieve email","err", nil)
		http.Error(w, "Unable to retrieve email", http.StatusInternalServerError)
		return
	}

	// Get the session key from cookie
	sid, err := r.Cookie("sid")
	fmt.Println("Session ID from cookie:", sid)
	if err != nil {
		h.logger.Error("Couldn't get the cookie for the session", "err", err)
	}

	// Get the session data
	session, _ := h.models.Passkey.GetSession(sid.Value)

	user, err := h.models.Passkey.GetUserByEmail(email)
	if err != nil {
		h.logger.Error("Failed to find user", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	credential, err := h.webauthn.FinishRegistration(user, session, r)
	if err != nil {
		h.logger.Error("Coudln't finish the WebAuthn Registration", err)
		// clean up sid cookie
		http.SetCookie(w, &http.Cookie{
			Name:  "sid",
			Value: "",
		})
		http.Error(w, "Couldn't finish registration", http.StatusBadRequest)
		return
	}

	// Store the credential object
	user.AddCredential(credential)
	h.models.Passkey.SaveUser(*user)
	// Delete the session data
	h.models.Passkey.DeleteSession(sid.Value)
	http.SetCookie(w, &http.Cookie{
		Name:  "sid",
		Value: "",
	})

	h.writeJSONResponse(w, "{'success': true}")

}

func (h *application) WebAuthnAuthenticationBeginHandler(w http.ResponseWriter, r *http.Request) {
	type CollectionRequest struct {
		Email string `json:"email"`
	}
	var req CollectionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode collection request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	email := req.Email

	h.logger.Info("Finding user " + email)

	user, err := h.models.Passkey.GetUserByEmail(email) // Find the user

	if err != nil {
		h.logger.Error("Failed to find user", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	options, session, err := h.webauthn.BeginLogin(user)
	if err != nil {
		h.logger.Error("Coudln't start a login", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Make a session key and store the sessionData values
	t, err := h.models.Passkey.GenSessionID()
	if err != nil {
		h.logger.Error("Coudln't create a session ID", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	h.models.Passkey.SaveSession(t, *session)

	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    t,
		Path:     "api/passkey/loginStart",
		MaxAge:   3600,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // TODO: SameSiteStrictMode maybe?
	})

	h.writeJSONResponse(w, options)
}

func (h *application) WebAuthnAuthenticationEndHandler(w http.ResponseWriter, r *http.Request) {
	// Get the session key from cookie
	sid, err := r.Cookie("sid")
	if err != nil {
		h.logger.Error("Coudln't get a session", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	// Get the session data stored from the function above
	session, _ := h.models.Passkey.GetSession(sid.Value)

	userID, err := strconv.Atoi(string(session.UserID)) // Convert []byte to int
	if err != nil {
		h.logger.Error("Failed to convert UserID to int", err)
		http.Error(w, "Invalid session data", http.StatusBadRequest)
		return
	}
	user, err := h.models.Passkey.GetUserByID(userID) // Get the user

	if err != nil {
		h.logger.Error("Failed to find user", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	credential, err := h.webauthn.FinishLogin(user, session, r)
	if err != nil {
		h.logger.Error("Coudln't finish login", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Handle credential.Authenticator.CloneWarning
	if credential.Authenticator.CloneWarning {
		h.logger.Error("Couldn't finish login", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// If login was successful
	user.UpdateCredential(credential)
	h.models.Passkey.SaveUser(*user)

	// Delete the session data
	h.models.Passkey.DeleteSession(sid.Value)
	http.SetCookie(w, &http.Cookie{
		Name:  "sid",
		Value: "",
	})

	// Add the new session cookie
	t, err := h.models.Passkey.GenSessionID()
	if err != nil {
		h.logger.Error("Couldn't generate session", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	h.models.Passkey.SaveSession(t, webauthn.SessionData{
		Expires: time.Now().Add(time.Hour),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    t,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // TODO: SameSiteStrictMode maybe?
	})

	type PasskeyResponse struct {
		Success bool   `json:"success"`
		JWT     string `json:"jwt"`
	}
	h.logger.Info("Sending JWT for " + user.Name)



	// get the user data from db:

	getUser,err:=h.models.Users.GetUserByEmail(user.Name)

	fmt.Println("User data",getUser)

	if err!=nil{
		h.logger.Error("Failed to find user", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Return success response
	response := PasskeyResponse{
		Success: true,
		JWT:     token.CreateJWT(*getUser, h.logger),
	}

	h.writeJSONResponse(w, response)
}
