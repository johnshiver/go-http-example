package users

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"

	"github.com/johnshiver/asapp_challenge/db"
)

type env struct {
	userService UserService
	db          *sqlx.DB
}

func InitRoutes(mux *httprouter.Router) {
	dbCon, err := db.Get()
	if err != nil {
		log.Fatal(err)
	}
	e := env{&UserManager{}, dbCon}
	e.addRoutes(mux)
}

func (e *env) addRoutes(mux *httprouter.Router) {
	mux.Handler("POST", "/users", http.HandlerFunc(e.CreateUser))
	mux.Handler("POST", "/login", http.HandlerFunc(e.Login))
}
