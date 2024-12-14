package main

import (
	"fmt"
	"log"
	"net/http"
)

const WEBPORT = "80"

type Config struct{}

func main() {
	app := Config{}
	fmt.Println("running server")
	serv := &http.Server{
		Addr:    fmt.Sprintf(":%s", WEBPORT),
		Handler: app.routes(),
	}

	err := serv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
