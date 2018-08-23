package shipManager


import(
    "kalaxia-game-api/database"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/model"
)

func GetFleet(id uint16) *model.Fleet{
	/**
     * Get Fleet data ( may return incomplete information).
     *  If the player is the owner of the ship all the data are send
     */
    
    var fleet model.Fleet
    if err := database.
        Connection.
        Model(&fleet).
        Column("fleet.*", "Player", "Location", "Journey","Journey.CurrentStep","Location.System").
        Where("fleet.id = ?", id).
        Select(); err != nil {
            panic(exception.NewHttpException(404, "Fleet not found", err))
    }
    
    return &fleet;
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
		PlayerId : player.Id,
        Location : planet,
		LocationId : planet.Id,
		Journey : nil,
	};
	
	if err := database.Connection.Insert(&fleet); err != nil {
		panic(exception.NewHttpException(500, "Fleet could not be created", err))
    }
	
	
	
	return &fleet;
}



func GetAllFleets(player *model.Player) []model.Fleet{
	var fleets []model.Fleet
    if err := database.
        Connection.
        Model(&fleets).
        Column("fleet.*", "Player","Location","Journey").
        Where("fleet.player_id = ?", player.Id).
        Select(); err != nil {
            panic(exception.NewHttpException(404, "Fleets not found", err))
    }
    return fleets
}


func GetFleetsOnPlanet(player *model.Player, planet *model.Planet) []model.Fleet {
	var fleets []model.Fleet
    if err := database.
        Connection.
        Model(&fleets).
        Column("fleet.*", "Player","Location","Journey").
        Where("fleet.player_id = ?", player.Id).
		Where("fleet.location_id = ?", planet.Id).
        Select(); err != nil {
            panic(exception.NewHttpException(404, "Fleets not found", err))
    }
	
    return fleets
}

func AssignShipsToFleet ( ships []*model.Ship, fleet *model.Fleet){
    
    for _,ship := range ships {
        isShipInTheCorrectLocation := IsShipInSamePositionAsFleet(*ship, *fleet);
    	
    	if !isShipInTheCorrectLocation{
    		panic(exception.NewHttpException(400, "wrong location", nil));
    	} 
        //ELSE
		//ship.IsShipInFleet = true;
        
		ship.Fleet = fleet;
		ship.FleetId =fleet.Id;
		ship.Hangar = nil;
        ship.HangarId = 0;
    	
    }
    
    UpdateShips(ships);
}

func RemoveShipsFromFleet (ships []*model.Ship){
    
    for _,ship := range ships {
        if ship.Fleet != nil {
    		if (ship.Fleet.Location != nil) {
                //ship.IsShipInFleet = false;
                
        		ship.Hangar = ship.Fleet.Location;
        		ship.HangarId = ship.Fleet.Location.Id;
        		ship.Fleet = nil;
                ship.FleetId = 0;
                
            } else {
                panic(exception.NewHttpException(400, "Fleet not stationed", nil));
            }
    	} else{
    		panic(exception.NewHttpException(400, "Ship is not in a fleet", nil));
    	}
    }
    
    UpdateShips(ships);
    
}

func GetFleetShip (fleet model.Fleet) []model.Ship{
    /*
     * get all ships in a fleet
     */
    var ships []model.Ship
    
    if err := database.
        Connection.
        Model(&ships).
        Column("ship.*", "Hangar", "Fleet", "Model").
        Where("construction_state_id IS NULL").
        Where("ship.fleet_id = ?", fleet.Id).
        Select(); err != nil {
            panic(exception.NewHttpException(404, "ship not found", err));
    }
    
    return ships;
}

func DeleteFleet(fleet *model.Fleet){
    
    if err := database.Connection.Delete(fleet); err != nil {
        panic(exception.NewHttpException(500, "Fleet could not be deleted", err));
    }
}

func UpdateFleet (fleet *model.Fleet){
    if err := database.Connection.Update(fleet); err != nil {
        panic(exception.NewHttpException(500, "Fleet could not be updated", err))
    }
}

func UpdateFleetInternal (fleet *model.Fleet){
    if err := database.Connection.Update(fleet); err != nil {
        panic(exception.NewException("Fleet could not be updated", err))
    }
}
