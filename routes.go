package main

import (
	"os"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"

	"kalaxia-game-api/controller"
	"kalaxia-game-api/handler"
)

type(
	Route struct {
		Name        string
		Method      string
		Pattern     string
		HandlerFunc http.HandlerFunc
		IsProtected bool
	}
	Routes []Route
)

func NewRouter() *mux.Router {
    router := mux.NewRouter().StrictSlash(true)
    for _, route := range routes {
			router.Handle(route.Pattern, handlers.LoggingHandler(os.Stdout, handler.JwtHandler(http.HandlerFunc(route.HandlerFunc), route.IsProtected))).Methods(route.Method)
    }
    return router
}

var routes = Routes{
		Route{
				"Authenticate",
				"POST",
				"/api/auth",
				controller.Authenticate,
				false,
		},
		Route{
				"Get Current Player",
				"GET",
				"/api/me",
				controller.GetCurrentPlayer,
				true,
		},
}
