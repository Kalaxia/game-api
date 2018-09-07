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


func CreateFleet(w http.ResponseWriter, r *http.Request){
	/*
	 * Create fleet on a given planet
	 */
	player := context.Get(r, "player").(*model.Player)
	
	data := utils.DecodeJsonRequest(r)
	
	//idPlanet, _ := strconv.ParseUint(data["planet_id"].(string), 10, 16);
	idPlanet := data["planet_id"].(float64);
	planet := manager.GetPlanet(uint16(idPlanet), player.Id)
	
	if (player.Id != planet.Player.Id) { // the player does not own the planet
		panic(exception.NewHttpException(http.StatusForbidden, "", nil));
	}
	// else
	utils.SendJsonResponse(w, 201,shipManager.CreateFleet(player,planet));
	
	
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
	
	fleetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	
	fleet := shipManager.GetFleet(uint16(fleetId));
	
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

func AssignShipsToFleet (w http.ResponseWriter, r *http.Request){
    /*
     * Assign mutliple ship to a fleet by theire id given in the body ( json) {"data-ships" : [id,id,id]}
     */
    
    player := context.Get(r, "player").(*model.Player)
    
    fleetId, _ := strconv.ParseUint(mux.Vars(r)["fleetId"], 10, 16)
    fleet := shipManager.GetFleet(uint16(fleetId));
    data := utils.DecodeJsonRequest(r)["data-ships"].([]interface{});
    
    var dataConverted []uint16
    for i := range data { // This is the solution according to stack overflow to convert into an array
        dataConverted = append(dataConverted, uint16(data[i].(float64)))
    }
    
    ships :=shipManager.GetShipsByIds(dataConverted)
    
    if (player.Id != fleet.Player.Id) { // the player does not own the planet
		panic(exception.NewHttpException(http.StatusForbidden, "", nil));
	}
    
    for i := range ships {
        /*idShip := data[i].(float64)
        ship := shipManager.GetShip(uint16(idShip));*/
        
        // TODO check on the Journey ?
    	// there is no verification if the this is in another fleet or not we can move ship inbetween fleet
    	
    	isShipInTheCorrectLocation := shipManager.IsShipInSamePositionAsFleet(*ships[i], *fleet);
    	
    	if (player.Id != ships[i].Hangar.Player.Id) { // this is the owner of the ship
    		panic(exception.NewHttpException(http.StatusForbidden, "", nil))
        }
        
        
    	if !isShipInTheCorrectLocation { // the fleet is on the right plante
    		panic(exception.NewHttpException(400, "Wrong location", nil));
    	}
        
    }
    
    shipManager.AssignShipsToFleet(ships,fleet);
    utils.SendJsonResponse(w, 202,"");
    
    
    
}

func RemoveShipsFromFleet (w http.ResponseWriter, r *http.Request){
    /*
     * Remove mutliple ship by theire id given in the body (json) {"data-ships" : [id,id,id]}
     */
    player := context.Get(r, "player").(*model.Player)
    data := utils.DecodeJsonRequest(r)["data-ships"].([]interface{});
    
    var dataConverted []uint16
    
    for i := range data {
        dataConverted = append(dataConverted, uint16(data[i].(float64)))
    }
    
    ships :=shipManager.GetShipsByIds(dataConverted)
    
    for i := range ships {
        
        /*idShip := data[i].(float64);
        ship := shipManager.GetShip(uint16(idShip));*/
        
        // TODO check on the Journey ?
        
    	if ships[i].Fleet == nil { // the ship is not in a fleet
    		panic(exception.NewHttpException(400, "Ship is not in a fleet", nil));
    	}
        if ships[i].Fleet.Location == nil { // the ship is not in a fleet
    		panic(exception.NewHttpException(400, "Fleet not on a planet", nil));
    	}
        
    	if ( player.Id != ships[i].Fleet.Player.Id || // this is the owner of the fleet
    	     ships[i].Fleet.Location.Player.Id != player.Id ) { // if the hangar is on a planet the player owns
    		panic(exception.NewHttpException(http.StatusForbidden, "", nil));
        }
        
        
    }
    
    shipManager.RemoveShipsFromFleet(ships);
    utils.SendJsonResponse(w, 202,"");
    
}

func GetFleetShip (w http.ResponseWriter, r *http.Request){
    /*
     * return all the ships in a fleet if the player controll the fleet
     */
    player := context.Get(r, "player").(*model.Player)
    
    fleetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    fleet := shipManager.GetFleet(uint16(fleetId));
    
    if (player.Id != fleet.Player.Id) { // the player does not own the fleet
		panic(exception.NewHttpException(http.StatusForbidden, "", nil));
	}
    
    utils.SendJsonResponse(w, 200,shipManager.GetFleetShip(*fleet));
}

func DeleteFleet (w http.ResponseWriter, r *http.Request){
	player := context.Get(r, "player").(*model.Player)
	fleetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	fleet := shipManager.GetFleet(uint16(fleetId));
	
	if fleet.Player.Id != player.Id {
        panic(exception.NewHttpException(http.StatusForbidden, "", nil));
	}
    if (fleet.Journey != nil){
        panic(exception.NewHttpException(400, "Cannot delete moving fleet", nil));
    }
    
    ships := shipManager.GetFleetShip(*fleet);
    if (len(ships) != 0){
        panic(exception.NewHttpException(400, "Cannot delete moving fleet", nil));
    }
    shipManager.DeleteFleet(fleet);
	utils.SendJsonResponse(w, 204,"Deleted");
}
