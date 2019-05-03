package main

import (
	"os"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"

	"kalaxia-game-api/api"
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
			os.Stdout, api.ErrorHandler(
				api.JwtHandler(
					api.AuthorizationHandler(
						http.HandlerFunc(route.HandlerFunc),
					), route.IsProtected),
				),
			),
		).Methods(route.Method)
	}
		
    router.
        PathPrefix("/api/resources/").
        Handler(http.StripPrefix("/api/resources/", http.FileServer(http.Dir("./resources/"))))
    return router
}

var routes = Routes{
	Route{
		"Authenticate",
		"POST",
		"/api/auth",
		api.Authenticate,
		false,
	},
	Route{
		"Create Server",
		"POST",
		"/api/servers",
		api.CreateServer,
		false,
	},
	Route{
		"Create Server",
		"POST",
		"/api/servers/{id}",
		api.RemoveServer,
		false,
	},
	Route{
		"Get Current Player",
		"GET",
		"/api/me",
		api.GetCurrentPlayer,
		true,
	},
	Route{
		"Update Current Planet",
		"PATCH",
		"/api/me/current-planet",
		api.UpdateCurrentPlanet,
		true,
	},
	Route{
		"Get Current Player Ship Models",
		"GET",
		"/api/me/ship-models",
		api.GetPlayerShipModels,
		true,
	},
	Route{
		"Get Current Player Ship Model",
		"GET",
		"/api/me/ship-models/{id}",
		api.GetShipModel,
		true,
	},
	Route{
		"Create Current Player Ship Model",
		"POST",
		"/api/me/ship-models",
		api.CreateShipModel,
		true,
	},
	Route{
		"Get Player Planets",
		"GET",
		"/api/players/{id}/planets",
		api.GetPlayerPlanets,
		true,
	},
	Route{
		"Launch Building Construction",
		"POST",
		"/api/planets/{id}/buildings",
		api.CreateBuilding,
		true,
	},
	Route{
		"Cancel Building Construction",
		"DELETE",
		"/api/planets/{planet-id}/buildings/{building-id}",
		api.CancelBuilding,
		true,
	},
	Route{
		"Launch Ships Construction",
		"POST",
		"/api/planets/{id}/ships",
		api.CreateShips,
		true,
	},
	/*******************************/
	// Combats
	Route{
		"Get Combat Report",
		"GET",
		"/api/combats/{id}",
		api.GetCombatReport,
		true,
	},
	Route{
		"Get Combat Reports",
		"GET",
		"/api/combats",
		api.GetCombatReports,
		true,
	},
	/*******************************/
	// Fleets
	Route{
		"Get Fleet Ships",
		"GET",
		"/api/fleets/{id}/ships",
		api.GetFleetShips,
		true,
	},
	Route{
		"Get Fleet Ship Groups",
		"GET",
		"/api/fleets/{id}/ships/groups",
		api.GetFleetShipGroups,
		true,
	},
	Route{
		"Transfer Ships Between Fleet and Hangar",
		"PATCH",
		"/api/fleets/{fleetId}/ships",
		api.TransferShips,
		true,
	},
	Route{ // data of the planet send in json
		"Create Fleet",
		"POST",
		"/api/fleets",
		api.CreateFleet,
		true,
	},
	Route{
		"Get Current Player Fleets",
		"GET",
		"/api/fleets",
		api.GetAllFleets,
		true,
	},
	Route{
		"Get travelling fleets",
		"GET",
		"/api/fleets/travelling",
		api.GetTravellingFleets,
		true,
	},
	Route{
		"Get Current Player Fleets on Planet",
		"GET",
		"/api/planets/{id}/fleets",
		api.GetPlanetFleets,
		true,
	},
	Route{
		"Delete Fleet",
		"DELETE",
		"/api/fleets/{id}",
		api.DeleteFleet,
		true,
	},
	Route{
		"Get Fleet",
		"GET",
		"/api/fleets/{id}",
		api.GetFleet,
		true,
	},
	/*******************************/
	// journeys
	Route{
		"Send Fleet On Journey",
		"POST",
		"/api/fleets/{id}/journey",
		api.SendFleetOnJourney,
		true,
	},
	Route{
		"Add Steps To Journey",
		"PATCH",
		"/api/fleets/{id}/journey",
		api.AddStepsToJourney,
		true,
	},
	Route{
		"Get Journey",
		"GET",
		"/api/fleets/{id}/journey",
		api.GetJourney,
		true,
	},
	Route{
		"Get Fleet Steps",
		"GET",
		"/api/fleets/{id}/steps",
		api.GetFleetSteps,
		true,
	},
	Route{
		"Remove Step And Following Form Journey Associated With Fleet",
		"DELETE",
		"/api/fleets/{id}/steps/{stepId}",
		api.RemoveFleetJourneyStep,
		true,
	},
	/*******************************/
	Route{
		"Get Hangar Ships",
		"GET",
		"/api/planets/{id}/ships",
		api.GetHangarShips,
		true,
	},
	Route{
		"Get Hangar Ship Groups",
		"GET",
		"/api/planets/{id}/ships/groups",
		api.GetHangarShipGroups,
		true,
	},
	Route{
		"Get Currently Constructing Ships",
		"GET",
		"/api/planets/{id}/ships/currently-constructing",
		api.GetCurrentlyConstructingShips,
		true,
	},
	Route{
		"Get Constructing Ships",
		"GET",
		"/api/planets/{id}/ships/constructing",
		api.GetConstructingShips,
		true,
	},
	Route{
		"Get Coming Fleets",
		"GET",
		"/api/planets/{id}/fleets/coming",
		api.GetComingFleets,
		true,
	},
	Route{
		"Get Leaving Fleets",
		"GET",
		"/api/planets/{id}/fleets/leaving",
		api.GetLeavingFleets,
		true,
	},
	/************** TRADE **********/
	Route{
		"Create offer",
		"POST",
		"/api/planets/{id}/offers",
		api.CreateOffer,
		true,
	},
	Route{
		"Cancel offer",
		"DELETE",
		"/api/offers/{id}",
		api.CancelOffer,
		true,
	},
	Route{
		"Get offer",
		"GET",
		"/api/offers/{id}",
		api.GetOffer,
		true,
	},
	Route{
		"Search offers",
		"POST",
		"/api/offers",
		api.SearchOffers,
		true,
	},
	Route{
		"Accept offer",
		"POST",
		"/api/offers/{id}/accept",
		api.AcceptOffer,
		true,
	},
	/************** PLAYER *********/
	Route{
		"Register Player",
		"POST",
		"/api/players",
		api.RegisterPlayer,
		true,
	},
	Route{
		"Get Player",
		"GET",
		"/api/players/{id}",
		api.GetPlayer,
		true,
	},
	Route{
		"Get Map",
		"GET",
		"/api/map",
		api.GetMap,
		true,
	},
	Route{
		"Get Sector Systems",
		"GET",
		"/api/systems",
		api.GetSectorSystems,
		true,
	},
	Route{
		"Get System",
		"GET",
		"/api/systems/{id}",
		api.GetSystem,
		true,
	},
	Route{
		"Get Planet",
		"GET",
		"/api/planets/{id}",
		api.GetPlanet,
		true,
	},
	Route{
		"Update Planet Settings",
		"PUT",
		"/api/planets/{id}/settings",
		api.AffectPopulationPoints,
		true,
	},
	Route{
		"Get Factions",
		"GET",
		"/api/factions",
		api.GetFactions,
		true,
	},
	Route{
		"Get Faction",
		"GET",
		"/api/factions/{id}",
		api.GetFaction,
		true,
	},
	Route{
		"Get Faction Planet Choices",
		"GET",
		"/api/factions/{id}/planet-choices",
		api.GetFactionPlanetChoices,
		true,
	},
	Route{
		"Get Faction Members",
		"GET",
		"/api/factions/{id}/members",
		api.GetFactionMembers,
		true,
	},
}
