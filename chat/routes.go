package chat

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"

	"github.com/johnshiver/asapp_challenge/db"
	"github.com/johnshiver/asapp_challenge/middleware"
)

type env struct {
	chatService Service
	db          *sqlx.DB
}

func InitRoutes(mux *httprouter.Router) {
	dbCon, err := db.Get()
	if err != nil {
		log.Fatal(err)
	}
	e := env{&Manager{}, dbCon}
	e.addRoutes(mux)
}

func (e *env) addRoutes(mux *httprouter.Router) {
	authMW := middleware.AuthenticateUser
	mux.Handler("GET", "/messages", authMW(http.HandlerFunc(e.GetMessages)))
	mux.Handler("POST", "/messages", authMW(http.HandlerFunc(e.CreateMessage)))
}
