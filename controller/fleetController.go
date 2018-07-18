package controller


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
	// string of the request route wich give the fleet id
	const ID_FLEET string = "idFleet"  //< you might want to change that
	// string of the request route wich give the ship id
	const ID_SHIP string = "idShip"  //< you might want to change that
	
	player := context.Get(r, "player").(*model.Player)
	
	idFleet, _ := strconv.ParseUint(mux.Vars(r)[ID_FLEET], 10, 16);
	idShip, _ := strconv.ParseUint(mux.Vars(r)[ID_SHIP], 10, 16)
	
	ship := shipManager.GetShip(uint16(idShip), player.Id)
	fleet := manager.GetFleet(uint16(idFleet), player.Id)
	
	
	// TODO check on the Journey ?
	// there is no verification if the this is in another fleet or not we can move ship inbetween fleet
	
	isShipInTheCorrectLocation := ( ! ship.IsShipInFleet && fleet.Location.Id !=  ship.Hangar.Id ) || // ship in Hangard and hangard same pos as the fleet
	  (ship.IsShipInFleet && ship.Fleet.Location.Id !=  fleet.Location.Id) // ship in fleet  and both fleet are a the same place
	
	if ( player.Id != fleet.Player.Id || // this is the owner of the fleet
	  player.Id != ship.Hangar.Player.Id){// this is the owner of the ship
		panic(exception.NewHttpException(http.StatusForbidden, "", nil))
    }
	if !isShipInTheCorrectLocation{ // the fleet is on the right plante
		panic(exception.NewHttpException(400, "Wrong location", nil));
	}
	
	manager.AssignShipToFleet(ship,fleet)
	
    utils.SendJsonResponse(w, 200,fleet)
}


func RemoveShipFormFleet(w http.ResponseWriter, r *http.Request){
	/**
	 * treat the http request to remove a ship form a fleet into an hangar
	 *
	 *
	 */
	
	// string of the request route wich give the ship id
	const ID_SHIP string = "idShip"  //< you migth want to change that
	
	player := context.Get(r, "player").(*model.Player)
	
	idShip, _ := strconv.ParseUint(mux.Vars(r)[ID_SHIP], 10, 16)
	
	ship := shipManager.GetShip(uint16(idShip), player.Id)
	
	// TODO check on the Journey ?
	
	if player.Id != ship.Hangar.Player.Id || // this is the owner of the fleet
	  ship.Fleet.Location.Player.Id !=   player.Id { // if the hangard is on a planet the player owns
		panic(exception.NewHttpException(http.StatusForbidden, "", nil));
    }
	if ! ship.IsShipInFleet { // the ship is not in a fleet
		panic(exception.NewHttpException(400, "Ship already in hangard", nil));
	}
	manager.AssignShipToHangard(ship)
	
    utils.SendJsonResponse(w, 200, nil /*TODO*/) // What do I return ?
}
