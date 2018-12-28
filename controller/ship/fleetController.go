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
	idPlanet := data["planet_id"].(float64)
	planet := manager.GetPlanet(uint16(idPlanet), player.Id)
	
	if (player.Id != planet.Player.Id) { // the player does not own the planet
		panic(exception.NewHttpException(http.StatusForbidden, "", nil))
	}
	// else
	utils.SendJsonResponse(w, 201, shipManager.CreateFleet(player,planet))
	
	
}

func GetAllFleets (w http.ResponseWriter, r *http.Request){
	/*
	 * return all the fleets a player controll
	 */
	player := context.Get(r, "player").(*model.Player)
	
	utils.SendJsonResponse(w, 200,shipManager.GetAllFleets(player))
}

func GetFleet (w http.ResponseWriter, r *http.Request){
	/*
	 * return a specifique fleet by id
	 */
	player := context.Get(r, "player").(*model.Player)
	
	fleetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	
	fleet := shipManager.GetFleet(uint16(fleetId))
	
	if fleet.Player.Id != player.Id {
		panic(exception.NewHttpException(http.StatusForbidden, "", nil))
	}
	
	utils.SendJsonResponse(w, 200,fleet)
}

func GetFleetsOnPlanet (w http.ResponseWriter, r *http.Request){
	/*
	 * return all the fleet on a fiven planet
	 */
	player := context.Get(r, "player").(*model.Player)
	
	idPlanet, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	planet := manager.GetPlanet(uint16(idPlanet), player.Id)
	
	if (player.Id != planet.Player.Id) { // the player does not own the planet
		panic(exception.NewHttpException(http.StatusForbidden, "", nil))
	}
	
	utils.SendJsonResponse(w, 200, shipManager.GetFleetsOnPlanet(player,planet))
}

func TransferShips (w http.ResponseWriter, r *http.Request){
    player := context.Get(r, "player").(*model.Player)
    
    fleetId, _ := strconv.ParseUint(mux.Vars(r)["fleetId"], 10, 16)
	fleet := shipManager.GetFleet(uint16(fleetId))
	
	data := utils.DecodeJsonRequest(r)
	modelId := int(data["model-id"].(float64))
	quantity := int(data["quantity"].(float64))
    
    if player.Id != fleet.Player.Id { // the player does not own the planet
		panic(exception.NewHttpException(http.StatusForbidden, "", nil))
	}
	
	if quantity > 0 {
		shipManager.AssignShipsToFleet(fleet, modelId, quantity)
	} else {
		shipManager.RemoveShipsFromFleet(fleet, modelId, -quantity)
	}
    w.WriteHeader(204)
    w.Write([]byte(""))  
}

func GetFleetShip (w http.ResponseWriter, r *http.Request){
    /*
     * return all the ships in a fleet if the player controll the fleet
     */
    player := context.Get(r, "player").(*model.Player)
    
    fleetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    fleet := shipManager.GetFleet(uint16(fleetId))
    
    if (player.Id != fleet.Player.Id) { // the player does not own the fleet
		panic(exception.NewHttpException(http.StatusForbidden, "", nil))
	}
    
    utils.SendJsonResponse(w, 200, shipManager.GetFleetShip(*fleet))
}

func GetFleetShipGroups (w http.ResponseWriter, r *http.Request) {
	player := context.Get(r, "player").(*model.Player)
    
    fleetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    fleet := shipManager.GetFleet(uint16(fleetId))
    
    if (player.Id != fleet.Player.Id) { // the player does not own the fleet
		panic(exception.NewHttpException(http.StatusForbidden, "", nil))
	}

	utils.SendJsonResponse(w, 200, shipManager.GetFleetShipGroups(*fleet))
}

func DeleteFleet (w http.ResponseWriter, r *http.Request){
	player := context.Get(r, "player").(*model.Player)
	fleetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	fleet := shipManager.GetFleet(uint16(fleetId))
	
	if fleet.Player.Id != player.Id {
        panic(exception.NewHttpException(http.StatusForbidden, "", nil))
	}
    if (fleet.Journey != nil){
        panic(exception.NewHttpException(400, "Cannot delete moving fleet", nil))
    }
    
    ships := shipManager.GetFleetShip(*fleet)
    if (len(ships) != 0){
        panic(exception.NewHttpException(400, "Cannot delete moving fleet", nil))
    }
    shipManager.DeleteFleet(fleet)
	utils.SendJsonResponse(w, 204,"Deleted")
}
