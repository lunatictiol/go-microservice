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
	
	return nil
}
