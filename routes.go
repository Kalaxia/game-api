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
		router.Handle(route.Pattern, handlers.LoggingHandler(
			os.Stdout, handler.JwtHandler(
				handler.AuthorizationHandler(
					http.HandlerFunc(route.HandlerFunc),
				), route.IsProtected),
			),
		).Methods(route.Method)
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
		"Create Server",
		"POST",
		"/api/servers",
		controller.CreateServer,
		false,
	},
	Route{
		"Get Current Player",
		"GET",
		"/api/me",
		controller.GetCurrentPlayer,
		true,
	},
	Route{
		"Get Player Planets",
		"GET",
		"/api/players/{id}/planets",
		controller.GetPlayerPlanets,
		true,
	},
	Route{
		"Launch Building Construction",
		"POST",
		"/api/planets/{id}/buildings",
		controller.CreateBuilding,
		true,
	},
	Route{
		"Register Player",
		"POST",
		"/api/players",
		controller.RegisterPlayer,
		true,
	},
	Route{
		"Get Map",
		"GET",
		"/api/map",
		controller.GetMap,
		true,
	},
	Route{
		"Get System",
		"GET",
		"/api/systems/{id}",
		controller.GetSystem,
		true,
	},
	Route{
		"Get Planet",
		"GET",
		"/api/planets/{id}",
		controller.GetPlanet,
		true,
	},
	Route{
		"Get Factions",
		"GET",
		"/api/factions",
		controller.GetFactions,
		true,
	},
	Route{
		"Get Faction Planet Choices",
		"GET",
		"/api/factions/{id}/planet-choices",
		controller.GetFactionPlanetChoices,
		true,
	},
}
