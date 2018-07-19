package route

import (
	"os"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"

	"kalaxia-game-api/controller/ship"
	"kalaxia-game-api/controller"
	"kalaxia-game-api/handler"
	"kalaxia-game-api/utils"
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
			os.Stdout, utils.ErrorHandler(
				handler.JwtHandler(
					handler.AuthorizationHandler(
						http.HandlerFunc(route.HandlerFunc),
					), route.IsProtected),
				),
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
		"Get Current Player Ship Models",
		"GET",
		"/api/me/ship-models",
		shipController.GetPlayerShipModels,
		true,
	},
	Route{
		"Get Current Player Ship Model",
		"GET",
		"/api/me/ship-models/{id}",
		shipController.GetShipModel,
		true,
	},
	Route{
		"Create Current Player Ship Model",
		"POST",
		"/api/me/ship-models",
		shipController.CreateShipModel,
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
		"Launch Ship Construction",
		"POST",
		"/api/planets/{id}/ships",
		shipController.CreateShip,
		true,
	},
	Route{
		"Assign Ship To Fleet",
		"POST",
		"/api/fleetId/{idFleet}/shipId/{idFleet}/assign",
		controller.AssignShipToFleet,
		true,
	},
	Route{
		"Remove Ship Form Fleet",
		"POST",
		"/api/fleetId/{idFleet}/remove",
		controller.RemoveShipFormFleet,
		true,
	},
	Route{
		"Creat Fleet",
		"POST",
		"/api/planets/{id}/creatfleet",
		controller.CreatFleet,
		true,
	},
	Route{
		"Get Hangar Ships",
		"GET",
		"/api/planets/{id}/ships",
		shipController.GetHangarShips,
		true,
	},
	Route{
		"Get Constructing Ships",
		"GET",
		"/api/planets/{id}/ships/constructing",
		shipController.GetConstructingShips,
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
		"Get Player",
		"GET",
		"/api/players/{id}",
		controller.GetPlayer,
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
		"Update Planet Settings",
		"PUT",
		"/api/planets/{id}/settings",
		controller.UpdatePlanetSettings,
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
		"Get Faction",
		"GET",
		"/api/factions/{id}",
		controller.GetFaction,
		true,
	},
	Route{
		"Get Faction Planet Choices",
		"GET",
		"/api/factions/{id}/planet-choices",
		controller.GetFactionPlanetChoices,
		true,
	},
	Route{
		"Get Faction Members",
		"GET",
		"/api/factions/{id}/members",
		controller.GetFactionMembers,
		true,
	},
}
