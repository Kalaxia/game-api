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
        Column("fleet.*", "Player", "Location", "Journey").
        Where("fleet.id = ?", id).
        Select(); err != nil {
            panic(exception.NewHttpException(404, "Fleet not found", err))
    }
    
    return &fleet
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

func AssignShipToFleet (ship *model.Ship,fleet *model.Fleet) {
	
	isShipInTheCorrectLocation := IsShipInSamePositionAsFleet(*ship, *fleet);
	
	if !isShipInTheCorrectLocation{
		panic(exception.NewHttpException(400, "wrong location", nil));
	} else{
		//ship.IsShipInFleet = true;
        
		ship.Fleet = fleet;
		ship.FleetId =fleet.Id;
		ship.Hangar = nil;
        ship.HangarId = 0;
        
		/*shipToUpdate := model.Ship{
            Id : ship.Id,
            FleetId : fleet.Id,
            ModelId : ship.ModelId,
            ConstructionStateId : ship.ConstructionStateId,
        }*/
        
		UpdateShip(ship);
        
        /*if err := database.Connection.Delete(ship.HangarId); err != nil {
            panic(exception.NewException("Ship Hangar could not be removed", err))
        }*/
	}
	
}

func AssignShipToHangar (ship *model.Ship){
	if ship.Fleet != nil {
		if (ship.Fleet.Location != nil) {
            //ship.IsShipInFleet = false;
            
    		ship.Hangar = ship.Fleet.Location;
    		ship.HangarId = ship.Fleet.Location.Id;
    		ship.Fleet = nil;
            ship.FleetId = 0;
            
            /*
            shipToUpdate := model.Ship{
                Id : ship.Id,
                HangarId : ship.Fleet.Location.Id,
                ModelId : ship.ModelId,
                ConstructionStateId : ship.ConstructionStateId,
            }
    		*/
    		UpdateShip(ship);
            
            /*
            if err := database.Connection.Delete(ship.FleetId); err != nil {
                panic(exception.NewException("Ship Fleet could not be removed", err))
            }
            */
        } else {
            panic(exception.NewHttpException(400, "Fleet not stationed", nil));
        }
	} else{
		panic(exception.NewHttpException(400, "Ship is not in a fleet", nil));
	}
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

func AssignMultipleShipsToFleet ( ships []*model.Ship, fleet *model.Fleet){
    for i := range ships {
        AssignShipToFleet(ships[i],fleet);
    }
}

func RemoveMultipleShipsFromFleet (ships []*model.Ship){
    
    for i := range ships {
        AssignShipToHangar(ships[i]);
    }
    
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
        Where("ship.fleet_id = ?", fleet.Id).
        Select(); err != nil {
            panic(exception.NewHttpException(404, "ship not found", err));
    }
    
    return ships;
}
