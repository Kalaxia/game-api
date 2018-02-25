package main

import (
	"os"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"

	"kalaxia-game-api/base"
	"kalaxia-game-api/faction"
	"kalaxia-game-api/galaxy"
	"kalaxia-game-api/player"
	"kalaxia-game-api/server"
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
		player.AuthenticateAction,
		false,
	},
	Route{
		"Create Server",
		"POST",
		"/api/servers",
		server.CreateServerAction,
		false,
	},
	Route{
		"Get Current Player",
		"GET",
		"/api/me",
		player.GetCurrentPlayerAction,
		true,
	},
	Route{
		"Get Player Planets",
		"GET",
		"/api/players/{id}/planets",
		galaxy.GetPlayerPlanetsAction,
		true,
	},
	Route{
		"Launch Building Construction",
		"POST",
		"/api/planets/{id}/buildings",
		base.CreateBuildingAction,
		true,
	},
	Route{
		"Register Player",
		"POST",
		"/api/players",
		player.RegisterPlayerAction,
		true,
	},
	Route{
		"Get Map",
		"GET",
		"/api/map",
		galaxy.GetMapAction,
		true,
	},
	Route{
		"Get System",
		"GET",
		"/api/systems/{id}",
		galaxy.GetSystemAction,
		true,
	},
	Route{
		"Get Planet",
		"GET",
		"/api/planets/{id}",
		galaxy.GetPlanetAction,
		true,
	},
	Route{
		"Get Factions",
		"GET",
		"/api/factions",
		faction.GetFactionsAction,
		true,
	},
	Route{
		"Get Faction Planet Choices",
		"GET",
		"/api/factions/{id}/planet-choices",
		faction.GetFactionPlanetChoicesAction,
		true,
	},
}
