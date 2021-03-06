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
		"Delete notification",
		"DELETE",
		"/api/me/notifications/{id}",
		api.DeleteNotification,
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
		"Update planet tax rate",
		"PATCH",
		"/api/planets/{id}/tax-rate",
		api.UpdatePlanetTaxRate,
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
		"Launch Building Compartment Construction",
		"POST",
		"/api/planets/{planet-id}/buildings/{building-id}/compartments",
		api.CreateBuildingCompartment,
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
	// Combats
	Route{
		"Get Combat Report Round",
		"GET",
		"/api/combats/{combatId}/rounds/{roundId}",
		api.GetCombatRound,
		true,
	},
	/*******************************/
	// Fleets
	Route{
		"Get Fleet Squadrons",
		"GET",
		"/api/fleets/{id}/squadrons",
		api.GetFleetSquadrons,
		true,
	},
	Route{
		"Create Fleet Ships",
		"POST",
		"/api/fleets/{id}/squadrons",
		api.CreateFleetSquadron,
		true,
	},
	Route{
		"Transfer Ships Between Fleet Squadron and Hangar",
		"PATCH",
		"/api/fleets/{fleetId}/squadrons/{squadronId}",
		api.AssignFleetSquadronShips,
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
		"Calculate Fleet Travel Duration",
		"POST",
		"/api/fleets/{id}/travel-duration",
		api.CalculateFleetTravelDuration,
		true,
	},
	Route{
		"Calculate Fleet Range",
		"POST",
		"/api/fleets/{id}/range",
		api.CalculateFleetRange,
		true,
	},
	Route{
		"Send Fleet On Journey",
		"POST",
		"/api/fleets/{id}/journey",
		api.SendFleetOnJourney,
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
		"Load fleet Cargo",
		"PATCH",
		"/api/fleets/{id}/load-cargo",
		api.LoadFleetCargo,
		true,
	},
	Route{
		"Unload fleet cargo",
		"PATCH",
		"/api/fleets/{id}/unload-cargo",
		api.UnloadFleetCargo,
		true,
	},
	/*******************************/
	Route{
		"Get Hangar Groups",
		"GET",
		"/api/planets/{id}/ships",
		api.GetHangarGroups,
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
		"/api/offers/search",
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
	/******** MAP ********/
	Route{
		"Get Map",
		"GET",
		"/api/map",
		api.GetMap,
		true,
	},
	Route{
		"Get Map Territories",
		"GET",
		"/api/map/territories",
		api.GetMapTerritories,
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
	Route{
		"Get Faction Motions",
		"GET",
		"/api/factions/{id}/motions",
		api.GetFactionCurrentMotions,
		true,
	},
	Route{
		"Get Faction previous Motions",
		"GET",
		"/api/factions/{id}/motions/previous",
		api.GetFactionPreviousMotions,
		true,
	},
	Route{
		"Create Faction Motion",
		"POST",
		"/api/factions/{id}/motions",
		api.CreateFactionMotion,
		true,
	},
	Route{
		"Get Faction Motion",
		"GET",
		"/api/factions/{id}/motions/{motion_id}",
		api.GetFactionMotion,
		true,
	},
	Route{
		"Vote Faction Motion",
		"POST",
		"/api/factions/{faction_id}/motions/{motion_id}/votes",
		api.VoteFactionMotion,
		true,
	},
	Route{
		"Get Faction Motion Vote",
		"GET",
		"/api/factions/{faction_id}/motions/{motion_id}/votes/me",
		api.GetFactionVote,
		true,
	},
	Route{
		"Get Faction Motion Votes",
		"GET",
		"/api/factions/{faction_id}/motions/{motion_id}/votes",
		api.GetFactionVotes,
		true,
	},
	Route{
		"Get Faction Wars",
		"GET",
		"/api/factions/{id}/wars",
		api.GetFactionWars,
		true,
	},
	Route{
		"Get Faction War",
		"GET",
		"/api/factions/{faction_id}/wars/{war_id}",
		api.GetFactionWar,
		true,
	},
	Route{
		"Get Faction Casus Belli",
		"GET",
		"/api/factions/{faction_id}/casus_belli/unanswered",
		api.GetFactionUnansweredCasusBelli,
		true,
	},
	Route{
		"Get Faction Casus Belli",
		"GET",
		"/api/factions/{faction_id}/casus_belli/{casus_belli_id}",
		api.GetFactionCasusBelli,
		true,
	},
	/********** RANKINGS *********/
	Route{
		"Get Territorial Ranking",
		"GET",
		"/api/rankings/territorial",
		api.GetTerritorialRanking,
		true,
	},
	Route{
		"Get Financial Ranking",
		"GET",
		"/api/rankings/financial",
		api.GetFinancialRanking,
		true,
	},
}
