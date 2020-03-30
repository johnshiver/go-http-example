package main

import (
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/johnshiver/asapp_challenge/config"
	"github.com/johnshiver/asapp_challenge/middleware"
)

func main() {
	mux := httprouter.New()
	initRoutes(mux)

	srv := http.Server{
		Addr:         config.Get().ServerPort,
		Handler:      middleware.Standard.Then(mux),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("starting server on port %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
