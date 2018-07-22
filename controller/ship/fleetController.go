package shipController


import (
    "net/http"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/manager"
	"kalaxia-game-api/manager/ship"
    "kalaxia-game-api/model"
	"kalaxia-game-api/utils"
    "strconv"
)

func AssignShipToFleet(w http.ResponseWriter, r *http.Request) {
	/**
	 * treat the http request to assign a ship in  a fleet
	 *
	 */
	
	
	player := context.Get(r, "player").(*model.Player)
	
	idFleet, _ := strconv.ParseUint(mux.Vars(r)["fleetId"], 10, 16);
	idShip, _ := strconv.ParseUint(mux.Vars(r)["shipId"], 10, 16)
	
	ship := shipManager.GetShip(uint16(idShip))
	fleet := shipManager.GetFleet(uint16(idFleet))
	
	
	// TODO check on the Journey ?
	// there is no verification if the this is in another fleet or not we can move ship inbetween fleet
	
	isShipInTheCorrectLocation := shipManager.IsShipInSamePositionAsFleet(*ship, *fleet);
	
	if ( player.Id != fleet.Player.Id || // this is the owner of the fleet
	  player.Id != ship.Hangar.Player.Id){// this is the owner of the ship
		panic(exception.NewHttpException(http.StatusForbidden, "", nil))
    }
	if !isShipInTheCorrectLocation{ // the fleet is on the right plante
		panic(exception.NewHttpException(400, "Wrong location", nil));
	}
	
	shipManager.AssignShipToFleet(ship,fleet)
	
    utils.SendJsonResponse(w, 200,fleet)
}


func RemoveShipFromFleet(w http.ResponseWriter, r *http.Request){
	/**
	 * treat the http request to remove a ship form a fleet into an hangar
	 *
	 *
	 */
	
	
	player := context.Get(r, "player").(*model.Player)
	
	idShip, _ := strconv.ParseUint(mux.Vars(r)["shipId"], 10, 16)
	
	ship := shipManager.GetShip(uint16(idShip))
	
	// TODO check on the Journey ?
	if ship.Fleet == nil{ // the ship is not in a fleet
		panic(exception.NewHttpException(400, "Ship already is not in a fleet", nil));
	}
	if player.Id != ship.Hangar.Player.Id || // this is the owner of the fleet
	  ship.Fleet.Location.Player.Id !=   player.Id { // if the hangard is on a planet the player owns
		panic(exception.NewHttpException(http.StatusForbidden, "", nil));
    }
	/* Depreciated
	if ! ship.IsShipInFleet { 
		panic(exception.NewHttpException(400, "Ship already in hangard", nil));
	}
	*/
	shipManager.AssignShipToHangard(ship)
	
    utils.SendJsonResponse(w, 200, nil /*TODO*/) // What do I return ?
}

func CreateFleet(w http.ResponseWriter, r *http.Request){
	/*
	 * Create fleet on a given planet
	 */
	player := context.Get(r, "player").(*model.Player)
	
	data := utils.DecodeJsonRequest(r)
	
	idPlanet, _ := strconv.ParseUint(data["planet_id"].(string), 10, 16)
	
	planet := manager.GetPlanet(uint16(idPlanet), player.Id)
	
	if (player.Id != planet.Player.Id) { // the player does not own the planet
		panic(exception.NewHttpException(http.StatusForbidden, "", nil));
	} else{
		utils.SendJsonResponse(w, 200,shipManager.CreateFleet(player,planet));
	}
	
}

func GetAllFleets (w http.ResponseWriter, r *http.Request){
	/*
	 * return all the fleets a player controll
	 */
	player := context.Get(r, "player").(*model.Player)
	
	utils.SendJsonResponse(w, 200,shipManager.GetAllFleets(player));
}

func GetFleet (w http.ResponseWriter, r *http.Request){
	/*
	 * return a specifique fleet by id
	 */
	player := context.Get(r, "player").(*model.Player)
	
	idFleet, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	
	fleet := shipManager.GetFleet(uint16(idFleet));
	
	if fleet.Player.Id != player.Id {
		panic(exception.NewHttpException(http.StatusForbidden, "", nil));
	}
	
	utils.SendJsonResponse(w, 200,fleet);
}

func GetFleetsOnPlanet (w http.ResponseWriter, r *http.Request){
	/*
	 * return all the fleet on a fiven planet
	 */
	
	player := context.Get(r, "player").(*model.Player)
	
	idPlanet, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	planet := manager.GetPlanet(uint16(idPlanet), player.Id)
	
	if (player.Id != planet.Player.Id) { // the player does not own the planet
		panic(exception.NewHttpException(http.StatusForbidden, "", nil));
	}
	
	utils.SendJsonResponse(w, 200,shipManager.GetFleetsOnPlanet(player,planet));
}
