package route

import (
	"os"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"

	"kalaxia-game-api/controller/ship"
	"kalaxia-game-api/controller/trade"
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
		"Create Server",
		"POST",
		"/api/servers/{id}",
		controller.RemoveServer,
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
		"Cancel Building Construction",
		"DELETE",
		"/api/planets/{planet-id}/buildings/{building-id}",
		controller.CancelBuilding,
		true,
	},
	Route{
		"Launch Ship Construction",
		"POST",
		"/api/planets/{id}/ships",
		shipController.CreateShip,
		true,
	},
	/*******************************/
	// Fleets
	Route{
		"Get Fleet Ships",
		"GET",
		"/api/fleets/{id}/ships",
		shipController.GetFleetShip,
		true,
	},
	Route{
		"Get Fleet Ship Groups",
		"GET",
		"/api/fleets/{id}/ships/groups",
		shipController.GetFleetShipGroups,
		true,
	},
	Route{
		"Transfer Ships Between Fleet and Hangar",
		"PATCH",
		"/api/fleets/{fleetId}/ships",
		shipController.TransferShips,
		true,
	},
	Route{ // data of the planet send in json
		"Create Fleet",
		"POST",
		"/api/fleets",
		shipController.CreateFleet,
		true,
	},
	Route{
		"Get Current Player Fleets",
		"GET",
		"/api/fleets",
		shipController.GetAllFleets,
		true,
	},
	Route{
		"Get Current Player Fleets on Planet",
		"GET",
		"/api/planets/{id}/fleets",
		shipController.GetFleetsOnPlanet,
		true,
	},
	Route{
		"Delete Fleet",
		"DELETE",
		"/api/fleets/{id}",
		shipController.DeleteFleet,
		true,
	},
	Route{
		"Get Fleet",
		"GET",
		"/api/fleets/{id}",
		shipController.GetFleet,
		true,
	},
	/*******************************/
	// journeys
	Route{
		"Send Fleet On Journey",
		"POST",
		"/api/fleets/{id}/journey",
		shipController.SendFleetOnJourney,
		true,
	},
	Route{
		"Add Steps To Journey",
		"PATCH",
		"/api/fleets/{id}/journey",
		shipController.AddStepsToJourney,
		true,
	},
	Route{
		"Get Journey",
		"GET",
		"/api/fleets/{id}/journey",
		shipController.GetJourney,
		true,
	},
	Route{
		"Get Fleet Steps",
		"GET",
		"/api/fleets/{id}/steps",
		shipController.GetFleetSteps,
		true,
	},
	Route{
		"Get Range",
		"GET",
		"/api/fleets/{id}/range",
		shipController.GetRange,
		true,
	},
	Route{
		"Get Time laws",
		"GET",
		"/api/fleets/{id}/times",
		shipController.GetTimeLaws,
		true,
	},
	Route{
		"Get Range",
		"GET",
		"/api/fleets/range",
		shipController.GetRange,
		true,
	},
	Route{
		"Get Time laws",
		"GET",
		"/api/fleets/times",
		shipController.GetTimeLaws,
		true,
	},
	Route{
		"Remove Step And Following Form Journey Associated With Fleet",
		"DELETE",
		"/api/fleets/{id}/steps/{idStep}",
		shipController.RemoveStepAndFollowingFormJourneyAssociatedWithFleet,
		true,
	},
	/*******************************/
	Route{
		"Get Hangar Ships",
		"GET",
		"/api/planets/{id}/ships",
		shipController.GetHangarShips,
		true,
	},
	Route{
		"Get Hangar Ship Groups",
		"GET",
		"/api/planets/{id}/ships/groups",
		shipController.GetHangarShipGroups,
		true,
	},
	Route{
		"Get Currently Constructing Ships",
		"GET",
		"/api/planets/{id}/ships/currently-constructing",
		shipController.GetCurrentlyConstructingShips,
		true,
	},
	Route{
		"Get Constructing Ships",
		"GET",
		"/api/planets/{id}/ships/constructing",
		shipController.GetConstructingShips,
		true,
	},
	/************** TRADE **********/
	Route{
		"Create offer",
		"POST",
		"/api/planets/{id}/offers",
		tradeController.CreateOffer,
		true,
	},
	Route{
		"Cancel offer",
		"DELETE",
		"/api/offers/{id}",
		tradeController.CancelOffer,
		true,
	},
	Route{
		"Get offer",
		"GET",
		"/api/offers/{id}",
		tradeController.GetOffer,
		true,
	},
	Route{
		"Search offers",
		"POST",
		"/api/offers",
		tradeController.SearchOffers,
		true,
	},
	Route{
		"Accept offer",
		"POST",
		"/api/offers/{id}/accept",
		tradeController.AcceptOffer,
		true,
	},
	/************** PLAYER *********/
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
		"Get Sector Systems",
		"GET",
		"/api/systems",
		controller.GetSectorSystems,
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
