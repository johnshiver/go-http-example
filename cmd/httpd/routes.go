package main

import (
	"github.com/julienschmidt/httprouter"

	"github.com/johnshiver/asapp_challenge/chat"
	"github.com/johnshiver/asapp_challenge/users"
)

type routerInit func(mux *httprouter.Router)

var routerInits = []routerInit{
	users.InitRoutes,
	chat.InitRoutes,
	applyHealthRoutes,
}

func initRoutes(mux *httprouter.Router) {
	for _, ri := range routerInits {
		ri(mux)
	}
}
