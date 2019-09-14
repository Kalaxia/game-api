package api

import(
    "github.com/gorilla/mux"
    "github.com/gorilla/context"
    "net/http"
    "strconv"
)

type(
    Fleet struct {
        TableName struct{} `json:"-" sql:"fleet__fleets"`

        Id uint16 `json:"id"`
        Player *Player `json:"player"`
        PlayerId uint16 `json:"-"`
        Location *Planet `json:"location"`
        LocationId uint16 `json:"-"`
        Journey *FleetJourney `json:"journey"`
        JourneyId uint16 `json:"-"`
        MapPosX float64 `json:"map_pos_x" sql:"map_pos_x"`
        MapPosY float64 `json:"map_pos_y" sql:"map_pos_y"`
        ShipSummary []FleetShipSummary `json:"ship_summary,omitempty" sql:"-"`
    }

    FleetShipSummary struct {
        Type string `json:"type"`
        NbShips uint16 `json:"nb_ships"`
    }
)

func CreateFleet(w http.ResponseWriter, r *http.Request) {
	player := context.Get(r, "player").(*Player)
	
	data := DecodeJsonRequest(r)
	
	planetId := data["planet_id"].(float64)
	planet := player.getPlanet(uint16(planetId))
	
	if (player.Id != planet.Player.Id) {
		panic(NewHttpException(http.StatusForbidden, "", nil))
	}
	SendJsonResponse(w, 201, player.createFleet(planet))
}

func GetAllFleets(w http.ResponseWriter, r *http.Request) {
	player := context.Get(r, "player").(*Player)
    
    fleets := player.getFleets()

	SendJsonResponse(w, 200, injectFleetsData(fleets))
}

func GetTravellingFleets(w http.ResponseWriter, r *http.Request) {
    SendJsonResponse(w, 200, getTravellingFleets())
}

func GetFleet(w http.ResponseWriter, r *http.Request) {
	player := context.Get(r, "player").(*Player)
	fleetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    fleet := getFleet(uint16(fleetId))
	
	if fleet.Player.Id != player.Id {
		panic(NewHttpException(http.StatusForbidden, "", nil))
	}
    if fleet.Journey != nil {
        fleet.Journey.Steps = fleet.Journey.getSteps()
    }
	SendJsonResponse(w, 200, fleet)
}

func GetComingFleets(w http.ResponseWriter, r *http.Request) {
	player := context.Get(r, "player").(*Player)
	planetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := player.getPlanet(uint16(planetId))
    
	if (player.Id != planet.Player.Id) {
		panic(NewHttpException(http.StatusForbidden, "", nil))
	}
	SendJsonResponse(w, 200, planet.getComingFleets())
}

func GetLeavingFleets(w http.ResponseWriter, r *http.Request) {
	player := context.Get(r, "player").(*Player)
	planetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := player.getPlanet(uint16(planetId))
    
	if (player.Id != planet.Player.Id) {
		panic(NewHttpException(http.StatusForbidden, "", nil))
	}
	SendJsonResponse(w, 200, planet.getLeavingFleets())
}

func GetPlanetFleets(w http.ResponseWriter, r *http.Request) {
	player := context.Get(r, "player").(*Player)
	planetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	planet := player.getPlanet(uint16(planetId))
	
	if (player.Id != planet.Player.Id) {
		panic(NewHttpException(http.StatusForbidden, "", nil))
	}
	SendJsonResponse(w, 200, planet.getFleets(player))
}

func TransferShips(w http.ResponseWriter, r *http.Request){
    player := context.Get(r, "player").(*Player)
    
    fleetId, _ := strconv.ParseUint(mux.Vars(r)["fleetId"], 10, 16)
	fleet := getFleet(uint16(fleetId))
	
	data := DecodeJsonRequest(r)
	modelId := int(data["model-id"].(float64))
	quantity := int(data["quantity"].(float64))
    
    if player.Id != fleet.Player.Id { // the player does not own the planet
		panic(NewHttpException(http.StatusForbidden, "", nil))
    }
    if fleet.isOnJourney() {
        panic(NewHttpException(http.StatusBadRequest, "fleet.errors.ship_transfer_on_journey", nil))
    }
    if fleet.Location.Player.Id != player.Id {
        panic(NewHttpException(http.StatusBadRequest, "fleet.errors.ship_transfer_on_foreign_planet", nil))
    }
	
	nbShips := 0
	if quantity > 0 {
		nbShips = fleet.assignShips(modelId, quantity)
	} else {
		nbShips = fleet.removeShips(modelId, -quantity)
	}
	SendJsonResponse(w, 200, struct {
		Quantity int `json:"quantity"`
	}{
		nbShips,
	}) 
}

func GetFleetShips(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)
    fleetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    fleet := getFleet(uint16(fleetId))
    
    if (player.Id != fleet.Player.Id) {
		panic(NewHttpException(http.StatusForbidden, "", nil))
	}
    
    SendJsonResponse(w, 200, fleet.getShips())
}

func GetFleetShipGroups(w http.ResponseWriter, r *http.Request) {
	player := context.Get(r, "player").(*Player)
    fleetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    fleet := getFleet(uint16(fleetId))
    
    if (player.Id != fleet.Player.Id) {
		panic(NewHttpException(http.StatusForbidden, "", nil))
	}
	SendJsonResponse(w, 200, fleet.getShipGroups())
}

func DeleteFleet(w http.ResponseWriter, r *http.Request){
	player := context.Get(r, "player").(*Player)
	fleetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	fleet := getFleet(uint16(fleetId))
	
	if fleet.Player.Id != player.Id {
        panic(NewHttpException(http.StatusForbidden, "", nil))
	}
	fleet.delete()
	
	w.WriteHeader(204)
	w.Write([]byte(""))
}

func getFleet(id uint16) *Fleet {
    fleet := &Fleet{}
    if err := Database.
        Model(fleet).
        Relation("Player.Faction").
        Relation("Journey.CurrentStep.PlanetStart.System").
        Relation("Journey.CurrentStep.PlanetFinal.System").
        Relation("Location.System").
        Relation("Location.Player.Faction").
        Where("fleet.id = ?", id).
        Select(); err != nil {
            panic(NewHttpException(404, "Fleet not found", err))
    }
    return fleet
}

func getTravellingFleets() []*Fleet {
    fleets := make([]*Fleet, 0)
    if err := Database.
        Model(&fleets).
        Relation("Player.Faction").
        Relation("Journey.CurrentStep.PlanetStart").
        Relation("Journey.CurrentStep.PlanetFinal").
        Where("fleet.journey_id IS NOT NULL").
        Select(); err != nil {
            panic(NewHttpException(404, "Could not retrieve travelling fleets", err))
    }
    return fleets
}

func (p *Player) createFleet(planet *Planet) *Fleet {
	fleet := &Fleet{
        Player : p,
		PlayerId : p.Id,
        Location : planet,
		LocationId : planet.Id,
		Journey : nil,
	}
	if err := Database.Insert(fleet); err != nil {
		panic(NewHttpException(500, "Fleet could not be created", err))
    }
	return fleet
}

func (p *Player) getFleets() []*Fleet {
	fleets := make([]*Fleet, 0)
    if err := Database.
        Model(&fleets).
        Relation("Player").
        Relation("Location").
        Relation("Journey.CurrentStep.PlanetStart").
        Relation("Journey.CurrentStep.PlanetFinal").
        Where("fleet.player_id = ?", p.Id).
        Select(); err != nil {
            panic(NewHttpException(404, "Fleets not found", err))
    }
    return fleets
}

func (p *Planet) getComingFleets() []*Fleet {
    fleets := make([]*Fleet, 0)
    steps := make([]FleetJourneyStep, 0)
    if err := Database.
        Model(&steps).
        Relation("Journey").
        Relation("PlanetFinal").
        Where("planet_final.id = ?", p.Id).
        Select(); err != nil {
            panic(NewException("Coming journey steps could not be retrieved", err))
    }
    if len(steps) == 0 {
        return fleets
    }
    journeyIds := make([]uint16, len(steps))
    for i, step := range steps {
        journeyIds[i] = step.JourneyId
    }

    if err := Database.
        Model(&fleets).
        Relation("Player.Faction").
        WhereIn("fleet.journey_id IN (?)", journeyIds).
        Select(); err != nil {
            panic(NewException("Could not retrieve coming fleets", err))
    }
    return fleets
}

func (p *Planet) getLeavingFleets() []*Fleet {
    fleets := make([]*Fleet, 0)
    steps := make([]FleetJourneyStep, 0)
    if err := Database.
        Model(&steps).
        Relation("Journey").
        Relation("PlanetStart").
        Where("planet_start.id = ?", p.Id).
        Where("step_number = 1").
        Select(); err != nil {
            panic(NewException("Coming journey steps could not be retrieved", err))
    }
    if len(steps) == 0 {
        return fleets
    }
    journeyIds := make([]uint16, len(steps))
    for i, step := range steps {
        journeyIds[i] = step.JourneyId
    }

    if err := Database.
        Model(&fleets).
        Relation("Player.Faction").
        WhereIn("fleet.journey_id IN (?)", journeyIds).
        Where("fleet.player_id = ?", p.Player.Id).
        Select(); err != nil {
            panic(NewException("Could not retrieve leaving fleets", err))
    }
    return fleets
}

func (p *Planet) getOrbitingFleets() []*Fleet {
    fleets := make([]*Fleet, 0)
    if err := Database.
        Model(&fleets).
        Relation("Player.Faction").
        Relation("Journey").
        Where("fleet.location_id = ?", p.Id).
        Where("fleet.journey_id IS NULL").
        Select(); err != nil {
            return fleets
    }
    return fleets
}

func (p *Planet) getFleets(player *Player) []*Fleet {
	fleets := make([]*Fleet, 0)
    if err := Database.
        Model(&fleets).
        Relation("Player.Faction").
        Relation("Location").
        Relation("Journey").
        Where("fleet.player_id = ?", player.Id).
		Where("fleet.location_id = ?", p.Id).
        Select(); err != nil {
            return fleets
    }
    return fleets
}

func (f *Fleet) assignShips(modelId int, quantity int) int {
    ships := f.Location.getHangarShipsByModel(modelId, quantity)
    for _, ship := range ships {
		ship.Fleet = f
		ship.FleetId = f.Id
		ship.Hangar = nil
        ship.HangarId = 0
        ship.update()
    }
    return len(ships)
}

func (f *Fleet) removeShips(modelId int, quantity int) int {
    ships := f.getShipsByModel(modelId, quantity)
    for _, ship := range ships {
        ship.Hangar = f.Location
        ship.HangarId = f.Location.Id
        ship.Fleet = nil
        ship.FleetId = 0
        ship.update()
    }
    return -len(ships)
}

func (f *Fleet) getShips() []Ship {
    ships := make([]Ship, 0)
    
    if err := Database.
        Model(&ships).
        Relation("Model").
        Where("construction_state_id IS NULL").
        Where("ship.fleet_id = ?", f.Id).
        Select(); err != nil {
            panic(NewHttpException(404, "fleet not found", err))
    }
    return ships
}


func (f *Fleet) getShipsByModel(modelId int, quantity int) []Ship {
	ships := make([]Ship, 0)
	
    if err := Database.
        Model(&ships).
        Relation("Hangar").
        Relation("Fleet").
        Where("construction_state_id IS NULL").
        Where("fleet_id = ?", f.Id).
        Where("model_id = ?", modelId).
        Limit(quantity).
        Select(); err != nil {
        	panic(NewHttpException(404, "Planet not found", err))
    }
    return ships
}

func (f *Fleet) getShipGroups() []ShipGroup {
    ships := make([]ShipGroup, 0)

    if err := Database.
        Model((*Ship)(nil)).
        ColumnExpr("model.id, model.name, model.type, model.frame_slug, count(*) AS quantity").
        Join("INNER JOIN ship__models as model ON model.id = ship.model_id").
        Group("model.id").
        Where("ship.construction_state_id IS NULL").
        Where("ship.fleet_id = ?", f.Id).
        Select(&ships); err != nil {
            panic(NewHttpException(404, "fleet not found", err))
    }
    return ships
}

func injectFleetsData(fleets []*Fleet) []*Fleet {
    for _, f := range fleets {
        f.ShipSummary = f.getShipSummary()
    }
    return fleets
}

func (f *Fleet) getShipSummary() []FleetShipSummary {
    summary := make([]FleetShipSummary, 0)
    if err := Database.Model((*Ship)(nil)).Column("model.type").ColumnExpr("count(*) as nb_ships").Join("INNER JOIN ship__models as model ON model.id = ship.model_id").Group("model.type").Where("fleet_id = ?", f.Id).Select(&summary); err != nil {
        panic(NewException("Could not retrieve fleet ship summary", err))
    }
    return summary
}

func (f *Fleet) delete() {
    if (f.Journey != nil){
        panic(NewHttpException(400, "Cannot delete moving fleet", nil))
    }
    if (len(f.getShips()) != 0){
        panic(NewHttpException(400, "Cannot delete fleet with remaining ships", nil))
    }
    if err := Database.Delete(f); err != nil {
        panic(NewHttpException(500, "Fleet could not be deleted", err))
    }
}

func (f *Fleet) update(){
    if err := Database.Update(f); err != nil {
        panic(NewException("Fleet could not be updated on UpdateFleet", err))
    }
}
