package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/johnshiver/asapp_challenge/db"
)

func applyHealthRoutes(router *httprouter.Router) {
	router.Handler("POST", "/check", http.HandlerFunc(healthCheck))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	var res int
	dbCon, err := db.Get()
	if err != nil {
		http.Error(w, "DB connection error", http.StatusInternalServerError)
		return
	}
	err = dbCon.Get(&res, "SELECT 1")
	if err != nil {
		log.Printf("health: querying database: %v", err)
		http.Error(w, "DB connection error", http.StatusInternalServerError)
		return
	}

	if res != 1 {
		http.Error(w, "Unexpected query result", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(map[string]string{"health": "ok"}); err != nil {
		http.Error(w, "Write error", http.StatusInternalServerError)
	}
}
