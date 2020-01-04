package api

import(
    "github.com/gorilla/mux"
    "github.com/gorilla/context"
    "net/http"
    "strconv"
    "time"
)

type(
    Fleet struct {
        tableName struct{} `pg:"fleet__fleets"`

        Id uint16 `json:"id"`
        Player *Player `json:"player"`
        PlayerId uint16 `json:"-"`
        Place *Place `json:"place"`
        PlaceId uint32 `json:"-"`
        Journey *FleetJourney `json:"journey"`
        JourneyId uint16 `json:"-"`
        Squadrons []*FleetSquadron `json:"squadrons" pg:",use_zero"`
        ShipSummary []FleetShipSummary `json:"ship_summary,omitempty" pg:"-"`
	    Cargo map[string]uint16 `json:"cargo" pg:",notnull,use_zero"`
        CreatedAt time.Time `json:"created_at"`
        DeletedAt time.Time `json:"deleted_at"`
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
	SendJsonResponse(w, 200, planet.getComingFleets(player))
}

func GetLeavingFleets(w http.ResponseWriter, r *http.Request) {
	player := context.Get(r, "player").(*Player)
	planetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := player.getPlanet(uint16(planetId))
    
	if (player.Id != planet.Player.Id) {
		panic(NewHttpException(http.StatusForbidden, "", nil))
	}
	SendJsonResponse(w, 200, planet.getLeavingFleets(player))
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

func LoadFleetCargo(w http.ResponseWriter, r *http.Request) {
	player := context.Get(r, "player").(*Player)
    fleetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    
    data := DecodeJsonRequest(r)
	planetId := uint16(data["planet_id"].(float64))
	planet := player.getPlanet(uint16(planetId))
    fleet := getFleet(uint16(fleetId))

    if fleet.Player.Id != player.Id || planet.Player.Id != player.Id {
        panic(NewHttpException(403, "access_denied", nil))
    }
    if !fleet.isOnPlanet(planet) {
        panic(NewHttpException(400, "fleet.delivery.not_on_planet", nil))
    }
    planet.transferResourcesToFleet(fleet, data)

    w.WriteHeader(204)
    w.Write([]byte(""))
}

func UnloadFleetCargo(w http.ResponseWriter, r *http.Request) {
	player := context.Get(r, "player").(*Player)
    fleetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    
    data := DecodeJsonRequest(r)
	planetId := uint16(data["planet_id"].(float64))
	planet := player.getPlanet(uint16(planetId))
    fleet := getFleet(uint16(fleetId))

    if fleet.Player.Id != player.Id {
        panic(NewHttpException(403, "access_denied", nil))
    }
    if !fleet.isOnPlanet(planet) {
        panic(NewHttpException(400, "fleet.delivery.not_on_planet", nil))
    }
    fleet.transferResourcesToPlanet(planet, data)

    w.WriteHeader(204)
    w.Write([]byte(""))
}

func DeleteFleet(w http.ResponseWriter, r *http.Request){
	player := context.Get(r, "player").(*Player)
	fleetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	fleet := getFleet(uint16(fleetId))
	
	if fleet.Player.Id != player.Id {
        panic(NewHttpException(http.StatusForbidden, "", nil))
    }
    if (fleet.Journey != nil){
        panic(NewHttpException(400, "Cannot delete moving fleet", nil))
    }
    if (len(fleet.getSquadrons()) != 0){
        panic(NewHttpException(400, "Cannot delete fleet with remaining ships", nil))
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
        Relation("Place.Planet.System").
        Relation("Place.Planet.Player.Faction").
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
        Relation("Journey.CurrentStep.StartPlace.Planet.System").
        Relation("Journey.CurrentStep.EndPlace.Planet.System").
        Where("fleet.journey_id IS NOT NULL").
        Select(); err != nil {
            panic(NewHttpException(404, "Could not retrieve travelling fleets", err))
    }
    return fleets
}

func (p *Player) createFleet(planet *Planet) *Fleet {
    place := NewPlace(planet, float64(planet.System.X), float64(planet.System.Y))
	fleet := &Fleet{
        Player : p,
        PlayerId : p.Id,
        Place: place,
        PlaceId: place.Id,
        Journey : nil,
        CreatedAt: time.Now(),
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
        Relation("Place.Planet").
        Relation("Journey.CurrentStep.StartPlace.Planet").
        Relation("Journey.CurrentStep.EndPlace.Planet").
        Where("fleet.deleted_at IS NULL").
        Where("fleet.player_id = ?", p.Id).
        Select(); err != nil {
            panic(NewHttpException(404, "Fleets not found", err))
    }
    for _, f := range fleets {
        f.injectData()
    }
    return fleets
}

func (p *Planet) getComingFleets(player *Player) []*Fleet {
    fleets := make([]*Fleet, 0)
    steps := make([]FleetJourneyStep, 0)
    if err := Database.
        Model(&steps).
        Relation("Journey").
        Relation("EndPlace.Planet").
        Where("end_place__planet.id = ?", p.Id).
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
    for i, f := range fleets {
        if f.PlayerId != player.Id {
            continue
        }
        fleets[i] = getFleet(f.Id)
        fleets[i].injectData()
    }
    return fleets
}

func (p *Planet) getLeavingFleets(player *Player) []*Fleet {
    fleets := make([]*Fleet, 0)
    steps := make([]FleetJourneyStep, 0)
    if err := Database.
        Model(&steps).
        Relation("Journey").
        Relation("StartPlace.Planet").
        Where("start_place__planet.id = ?", p.Id).
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
    for i, f := range fleets {
        if f.PlayerId != player.Id {
            continue
        }
        fleets[i] = getFleet(f.Id)
        fleets[i].injectData()
    }
    return fleets
}

func (p *Planet) getOrbitingFleets() []*Fleet {
    fleets := make([]*Fleet, 0)
    if err := Database.
        Model(&fleets).
        Relation("Place.Planet").
        Relation("Player.Faction").
        Relation("Journey").
        Where("place__planet.id = ?", p.Id).
        Where("fleet.deleted_at IS NULL").
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
        Relation("Place.Planet").
        Relation("Journey").
        Where("fleet.deleted_at IS NULL").
        Where("fleet.player_id = ?", player.Id).
		Where("place__planet.id = ?", p.Id).
        Select(); err != nil {
            return fleets
    }
    for _, f := range fleets {
        f.injectData()
    }
    return fleets
}

func (f *Fleet) injectData() {
    f.ShipSummary = f.getShipSummary()
}

func (f *Fleet) getShipSummary() []FleetShipSummary {
    summary := make([]FleetShipSummary, 0)
    if err := Database.Model((*FleetSquadron)(nil)).Column("model.type").ColumnExpr("SUM(quantity) as nb_ships").Join("INNER JOIN ship__models as model ON model.id = fleet_squadron.ship_model_id").Group("model.type").Where("fleet_id = ?", f.Id).Select(&summary); err != nil {
        panic(NewException("Could not retrieve fleet ship summary", err))
    }
    return summary
}

func (f *Fleet) transferResourcesToPlanet(p *Planet, data map[string]interface{}) {
    resource := data["resource"].(string)
    quantity := uint16(data["quantity"].(float64))

    f.unloadCargo(resource, quantity)
    p.Storage.storeResource(resource, int16(quantity))
    p.Storage.update()
    f.update()
    f.notifyDelivery(p, quantity, true)
}

func (p *Planet) transferResourcesToFleet(f *Fleet, data map[string]interface{}) {
    resource := data["resource"].(string)
    quantity := uint16(data["quantity"].(float64))

    if !p.Storage.hasResource(resource, quantity) {
        panic(NewHttpException(400, "planet.storage.not_enough_resources", nil))
    }
    f.loadCargo(resource, quantity)
    p.Storage.storeResource(resource, -int16(quantity))
    p.Storage.update()
    f.update()
}

func (f *Fleet) deliver(p *Planet, deliveryData map[string]interface{}) {
    totalQuantity := uint16(0)
    for _, d := range deliveryData["resources"].([]interface{}) {
        data := d.(map[string]interface{})
        resource := data["resource"].(string)
        quantity := uint16(data["quantity"].(float64))

        f.unloadCargo(resource, quantity)
        p.Storage.storeResource(resource, int16(quantity))
        totalQuantity += quantity
    }
    f.notifyDelivery(p, totalQuantity, false)
    f.update()
    p.Storage.update()
}

func (f *Fleet) notifyDelivery(p *Planet, quantity uint16, manualDelivery bool) {
    data := map[string]interface{}{
        "fleet_name": f.Id,
        "owner_name": f.Player.Pseudo,
        "planet_name": p.Name,
        "player_name": p.Player.Pseudo,
        "quantity": quantity,
    }
    if p.Player != nil && f.PlayerId != p.PlayerId {
        p.Player.notify(NotificationTypeTrade, "fleet.delivery.received", data)
    }
    if !manualDelivery {
        f.Player.notify(NotificationTypeTrade, "fleet.delivery.success", data)
    }
}

func (f *Fleet) loadCargo(resource string, quantity uint16) {
    if quantity > f.getAvailableCargoCapacity() {
        panic(NewHttpException(400, "fleet.delivery.cargo_overload", nil))
    }
    if f.Cargo == nil {
        f.Cargo = make(map[string]uint16, 0)
    }
    if _, ok := f.Cargo[resource]; !ok {
        f.Cargo[resource] = 0
    }
    f.Cargo[resource] += quantity
}

func (f *Fleet) unloadCargo(resource string, quantity uint16) {
    if _, ok := f.Cargo[resource]; !ok {
        panic(NewHttpException(400, "fleet.delivery.missing_resource", nil))
    }
    if f.Cargo[resource] < quantity {
        panic(NewHttpException(400, "fleet.delivery.not_enough_resources", nil))
    } else if f.Cargo[resource] == quantity {
        delete(f.Cargo, resource)
    } else {
        f.Cargo[resource] -= quantity
    }
}

func (f *Fleet) getAvailableCargoCapacity() (capacity uint16) {
    capacity = f.getCargoCapacity()
    for _, quantity := range f.Cargo {
        capacity -= quantity
    }
    return
}

func (f *Fleet) hasResource(resource string, quantity uint16) bool {
    q, ok := f.Cargo[resource]
    return ok && q >= quantity
}

func (f *Fleet) getCargoCapacity() (capacity uint16) {
    for _, s := range f.getSquadrons() {
        capacity += s.ShipModel.Stats[ShipStatCargo] * uint16(s.Quantity)
    }
    return
}

func (f *Fleet) delete() {
    f.DeletedAt = time.Now()
    f.update()
}

func (f *Fleet) update(){
    if err := Database.Update(f); err != nil {
        panic(NewException("Fleet could not be updated on UpdateFleet", err))
    }
}
