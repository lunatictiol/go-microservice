package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   authPayload `json:"auth,omitempty"`
}

type authPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "message from broker",
	}
	_ = app.WriteJson(w, http.StatusOK, payload)

}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := app.ReadJson(w, r, &requestPayload)
	if err != nil {
		log.Println("Hereeeee")
		app.WriteErrorJson(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)

	default:
		app.WriteErrorJson(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a authPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	req, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.WriteErrorJson(w, err)
		return
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		app.WriteErrorJson(w, err)
		return
	}
	defer res.Body.Close()
	log.Println(http.StatusAccepted)
	log.Println(res.StatusCode)
	if res.StatusCode == http.StatusUnauthorized {
		app.WriteErrorJson(w, errors.New("invalid credentials"))
		return

	} else if res.StatusCode != http.StatusAccepted {
		log.Println("I am the one")
		app.WriteErrorJson(w, errors.New("invalid credentials"))
		return

	}

	var jsonAuthResponse jsonResponse

	err = json.NewDecoder(res.Body).Decode(&jsonAuthResponse)
	if err != nil {
		app.WriteErrorJson(w, err)
		return
	}
	if jsonAuthResponse.Error {
		app.WriteErrorJson(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = jsonAuthResponse.Message
	payload.Data = jsonAuthResponse.Data

	app.WriteJson(w, http.StatusAccepted, payload)

}
