package main

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "message from broker",
	}
	_ = app.WriteJson(w, http.StatusOK, payload)

}
