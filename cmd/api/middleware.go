package main

import (
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
)

// 👇 a logging middleware
func (app *Config) Authenticate(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// var requestPayload struct {
		// 	// Email    string `json:"email"`
		// 	Username string `json:"username"`
		// 	Password string `json:"password"`
		// }

		username, password, ok := r.BasicAuth()

		// err := app.readJSON(w, r, &requestPayload)
		// if err != nil {
		// 	app.errorJSON(w, err, http.StatusBadRequest)
		// 	return
		// }

		log.Info().Msgf("Username:%s and Password:%s", username, password)

		if ok {
			// validate the user against the database
			user, err := app.Models.User.GetByEmail(username)
			if err != nil {
				log.Info().Msgf("error while retrieving user from db: %v", err)
				app.errorJSON(w, errors.New("invalid credentials: Username"), http.StatusBadRequest)
				return
			}

			// send errors for invalid users
			valid, err := user.PasswordMatches(password)
			if err != nil || !valid {
				app.errorJSON(w, errors.New("invalid credentials: Password"), http.StatusBadRequest)
				return
			}

		}

		handler.ServeHTTP(w, r)
	})
}
