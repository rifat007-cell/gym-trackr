package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/golang-jwt/jwt"
	"github.com/tanvir-rifat007/gymBuddy/token"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				app.logger.Error("Panic recovered", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}


func (app *application) withSentry(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hub := app.sentry.Clone()
		hub.Scope().SetRequest(r)

		defer func() {
			if err := recover(); err != nil {
				hub.Recover(err)
				hub.Flush(2 * time.Second)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		r = r.WithContext(sentry.SetHubOnContext(r.Context(), hub))
		next.ServeHTTP(w, r)
	})
}


func (h *application) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix if present
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		// Parse and validate the token

		token, err := jwt.Parse(tokenStr,
			func(t *jwt.Token) (interface{}, error) {
				// Ensure the signing method is HMAC
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(token.GetJWTSecret(h.logger)), nil
			},
		)
		if err != nil || !token.Valid {
			h.writeJSON(w, http.StatusUnauthorized, envelope{"error": "Invalid token"}, nil)
			return
		}

		// Extract claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			h.writeJSON(w, http.StatusUnauthorized, envelope{"error": "Invalid token claims"}, nil)
			return
		}

		// Get the email from claims
		email, ok := claims["email"].(string)
		if !ok {
			h.writeJSON(w, http.StatusUnauthorized, envelope{"error": "Email not found in token"}, nil)
			return
		}

		// Inject email into the request context
		ctx := context.WithValue(r.Context(), "email", email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
