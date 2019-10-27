package api

import(
    "github.com/gorilla/mux"
    "github.com/gorilla/context"
    "net/http"
    "strconv"
)

type(
    Fleet struct {
        tableName struct{} `json:"-" pg:"fleet__fleets"`

        Id uint16 `json:"id"`
        Player *Player `json:"player"`
        PlayerId uint16 `json:"-"`
        Location *Planet `json:"location"`
        LocationId uint16 `json:"-"`
        Journey *FleetJourney `json:"journey"`
        JourneyId uint16 `json:"-"`
        MapPosX float64 `json:"map_pos_x" pg:"map_pos_x"`
        MapPosY float64 `json:"map_pos_y" pg:"map_pos_y"`
        Squadrons []*FleetSquadron `json:"squadrons" pg:",use_zero"`
        ShipSummary []FleetShipSummary `json:"ship_summary,omitempty" pg:"-"`
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

	//SendJsonResponse(w, 200, injectFleetsData(fleets))
	SendJsonResponse(w, 200, fleets)
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
        Relation("Journey.CurrentStep.StartPlace.Planet.System").
        Relation("Journey.CurrentStep.EndPlace.Planet.System").
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
        Relation("Journey.CurrentStep.StartPlace.Planet").
        Relation("Journey.CurrentStep.EndPlace.Planet").
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
        Relation("Journey.CurrentStep.StartPlace.Planet").
        Relation("Journey.CurrentStep.EndPlace.Planet").
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
        Relation("StartPlace.Planet").
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

// func injectFleetsData(fleets []*Fleet) []*Fleet {
//     for _, f := range fleets {
//         f.ShipSummary = f.getShipSummary()
//     }
//     return fleets
// }

// func (f *Fleet) getShipSummary() []FleetShipSummary {
//     summary := make([]FleetShipSummary, 0)
//     if err := Database.Model((*Ship)(nil)).Column("model.type").ColumnExpr("count(*) as nb_ships").Join("INNER JOIN ship__models as model ON model.id = ship.model_id").Group("model.type").Where("fleet_id = ?", f.Id).Select(&summary); err != nil {
//         panic(NewException("Could not retrieve fleet ship summary", err))
//     }
//     return summary
// }

func (f *Fleet) delete() {
    if (f.Journey != nil){
        panic(NewHttpException(400, "Cannot delete moving fleet", nil))
    }
    if (len(f.getSquadrons()) != 0){
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
