package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) ReadJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxByte := 1048567

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxByte))
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(data)
	if err != nil {
		return err
	}
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only have a single JSON Value ")
	}
	return nil
}

func (app *Config) WriteJson(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for key, v := range headers[0] {
			w.Header()[key] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) WriteErrorJson(w http.ResponseWriter, err error, status ...int) error {

	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}
	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.WriteJson(w, statusCode, payload)
}
