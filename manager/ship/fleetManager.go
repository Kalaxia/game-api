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
        Column("Player", "Journey.CurrentStep","Location.System").
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
	
	//fleetJourney = nil
	
	fleet := model.Fleet{
        Player : player,
		PlayerId : player.Id,
        Location : planet,
		LocationId : planet.Id,
		Journey : nil,
	}
	
	if err := database.Connection.Insert(&fleet); err != nil {
		panic(exception.NewHttpException(500, "Fleet could not be created", err))
    }
	
	
	
	return &fleet
}



func GetAllFleets(player *model.Player) []model.Fleet{
	var fleets []model.Fleet
    if err := database.
        Connection.
        Model(&fleets).
        Column("Player","Location","Journey").
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
        Column( "Player","Location","Journey").
        Where("fleet.player_id = ?", player.Id).
		Where("fleet.location_id = ?", planet.Id).
        Select(); err != nil {
            panic(exception.NewHttpException(404, "Fleets not found", err))
    }
	
    return fleets
}

func AssignShipsToFleet (fleet *model.Fleet, modelId int, quantity int) {
    for _, ship := range GetHangarShipsByModel(fleet.Location, modelId, quantity) {
		ship.Fleet = fleet
		ship.FleetId = fleet.Id
		ship.Hangar = nil
        ship.HangarId = 0
        UpdateShip(&ship)
    }
}

func RemoveShipsFromFleet (fleet *model.Fleet, modelId int, quantity int) {
    if (fleet.Location == nil) {
        panic(exception.NewHttpException(400, "Fleet not stationed", nil))
    }
    for _, ship := range GetFleetShipsByModel(fleet, modelId, quantity) {
        ship.Hangar = fleet.Location
        ship.HangarId = fleet.Location.Id
        ship.Fleet = nil
        ship.FleetId = 0
        UpdateShip(&ship)
    }
}

func GetFleetShip (fleet model.Fleet) []model.Ship {
    var ships []model.Ship
    
    if err := database.
        Connection.
        Model(&ships).
        Column("Model").
        Where("construction_state_id IS NULL").
        Where("ship.fleet_id = ?", fleet.Id).
        Select(); err != nil {
            panic(exception.NewHttpException(404, "fleet not found", err))
    }
    
    return ships
}


func GetFleetShipsByModel(fleet *model.Fleet, modelId int, quantity int) []model.Ship {
    ships := make([]model.Ship, 0)
    if err := database.
        Connection.
        Model(&ships).
        Column("Hangar", "Fleet").
        Where("construction_state_id IS NULL").
        Where("fleet_id = ?", fleet.Id).
        Where("model_id = ?", modelId).
        Limit(quantity).
        Select(); err != nil {
        panic(exception.NewHttpException(404, "Planet not found", err))
    }
    return ships
}

func GetFleetShipGroups(fleet model.Fleet) []model.ShipGroup {
    ships := make([]model.ShipGroup, 0)

    if err := database.
        Connection.
        Model((*model.Ship)(nil)).
        ColumnExpr("model.id, model.name, model.type, model.frame_slug, count(*) AS quantity").
        Join("INNER JOIN ship__models as model ON model.id = ship.model_id").
        Group("model.id").
        Where("ship.construction_state_id IS NULL").
        Where("ship.fleet_id = ?", fleet.Id).
        Select(&ships); err != nil {
            panic(exception.NewHttpException(404, "fleet not found", err))
    }

    return ships
}

func DeleteFleet(fleet *model.Fleet){
    
    if err := database.Connection.Delete(fleet); err != nil {
        panic(exception.NewHttpException(500, "Fleet could not be deleted", err))
    }
}


func UpdateFleet (fleet *model.Fleet){
    if err := database.Connection.Update(fleet); err != nil {
        panic(exception.NewException("Fleet could not be updated on UpdateFleet", err))
    }
}
