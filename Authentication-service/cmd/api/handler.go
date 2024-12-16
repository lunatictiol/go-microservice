package main

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var authPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.ReadJson(w, r, &authPayload)
	if err != nil {
		app.WriteErrorJson(w, err, http.StatusBadRequest)
		return
	}
	log.Println(authPayload.Email, authPayload.Password)
	user, err := app.Models.User.GetByEmail(authPayload.Email)
	log.Println(user.Email, user.Password)
	if err != nil {
		app.WriteErrorJson(w, errors.New("Invalid credentials"), http.StatusBadRequest)
		return
	}
	valid, err := user.PasswordMatches(authPayload.Password)
	log.Println(valid)

	if err != nil || !valid {
		log.Println(valid)
		app.WriteErrorJson(w, errors.New("Invalid credentials"), http.StatusBadRequest)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: "user authenticated",
		Data:    user,
	}
	app.WriteJson(w, http.StatusAccepted, payload)
}
