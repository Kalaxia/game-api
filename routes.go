package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"kalaxia-game-api/controller"
)

type(
	Route struct {
		Name        string
		Method      string
		Pattern     string
		HandlerFunc http.HandlerFunc
	}
	Routes []Route
)

func NewRouter() *mux.Router {
    router := mux.NewRouter().StrictSlash(true)
    for _, route := range routes {
			router.HandleFunc(route.Pattern, route.HandlerFunc).Methods(route.Method)
    }
    return router
}

var routes = Routes{
		Route{
				"Authenticate player",
				"POST",
				"/auth",
				controller.AuthenticatePlayer,
		},
}
