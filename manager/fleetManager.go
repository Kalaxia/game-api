package manager

import(
    "kalaxia-game-api/database"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/model"
)

func GetFleet(id uint16, playerId uint16) *model.Fleet{
	/**
     * Get Fleet data ( may return incomplete information).
     *  If the player is the owner of the ship all the data are send
     * 
     */
    // TODO Someone pls check if
    
    var fleet model.Fleet
    if err := database.
        Connection.
        Model(&fleet).
        Column("fleet.*", "Player", "Location", "Journey").
        Where("planet.id = ?", id).
        Select(); err != nil {
            panic(exception.NewHttpException(404, "Fleet not found", err))
    }
    if fleet.Player != nil && playerId == fleet.Player.Id {
        getShipOwnerData(&fleet)
    }
    return &fleet
}


func getShipOwnerData(fleet *model.Fleet) {
   // TODO 
}


func AssignShipToFleet (ship *model.Ship,fleet *model.Fleet) {
	
	isShipInTheCorrectLocation := ( ! ship.IsShipInFleet && fleet.Location.Id !=  ship.Hangar.Id ) || // ship in Hangard and hangard same pos as the fleet
	  (ship.IsShipInFleet && ship.Fleet.Location.Id !=  fleet.Location.Id); // ship in fleet  and both fleet are a the same place
	
	if !isShipInTheCorrectLocation{
		panic(exception.NewHttpException(400, "wrong location", nil));
	}
	else{
		ship.IsShipInFleet = true;
		ship.Fleet = fleet;
		ship.FleetId=fleet.Id;
		
	}
	
}

func AssignShipToHangard (ship *model.Ship){
	ship.IsShipInFleet = false;
	ship.Hangar = ship.Fleet.Location;
	ship.HangarId = ship.Hangar.Id;
}
