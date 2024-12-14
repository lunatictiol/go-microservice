package main

import (
	"encoding/json"
	"net/http"
)


func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "message from broker",
	}

	out, _ := json.MarshalIndent(payload, "", "\t")
	w.Header().Set("COntent-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)

}
