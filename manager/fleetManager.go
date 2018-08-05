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
        getFleetOwnerData(&fleet)
    }
    return &fleet
}


func getFleetOwnerData(fleet *model.Fleet) {
   // TODO 
}




func CreateFleet (player *model.Player, planet *model.Planet) *model.Fleet{
	
	/*
	fleetJourney := model.FleetJourney{ // TODO ?
		
	}
	
	if err := database.Connection.Insert(&fleetJourney); err != nil {
      panic(exception.NewHttpException(500, "Fleet Journey could not be created", err))
    }
    */
	
	//fleetJourney = nil;
	
	fleet := model.Fleet{
        Player : player,
        Location : planet,
        Journey : nil,
	};
	
	if err := database.Connection.Insert(&fleet); err != nil {
		panic(exception.NewHttpException(500, "Fleet could not be created", err))
    }
	
	
	
	return &fleet;
}

