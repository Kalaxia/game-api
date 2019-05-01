package api

import(
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "net/http"
    "strconv"
    "time"
)

type(
    ShipConstructionState struct {
        TableName struct{} `json:"-" sql:"ship__construction_states"`

        Id uint32 `json:"id"`
        CurrentPoints uint8 `json:"current_points" sql:",notnull"`
        Points uint8 `json:"points"`
    }
    Ship struct {
        TableName struct{} `json:"-" sql:"ship__ships"`

        Id uint32 `json:"id"`
        HangarId uint16 `json:"-"` //< if  the ship is not in a hangar the id wil be nil.
        Hangar *Planet `json:"hangar"` //< if  the ship is not in a hangar  the pointer will be nil
        FleetId uint16 `json:"-"` //< if  the ship is not in a fleet the id wil be nil
        Fleet *Fleet `json:"fleet"` //< if  the ship is not in a fleet the pointer will be nil
        //IsShipInFleet bool `json:"isShipInFleet"` //< Depreciated : this is used when a ship is try to be removed form the hangar in order to avoid teleporting ships
        ModelId uint `json:"-"` 
        Model *ShipModel `json:"model"`
        CreatedAt time.Time `json:"created_at"`
        ConstructionStateId uint32 `json:"-"`
        ConstructionState *ShipConstructionState `json:"construction_state"`
        // Combat fields
        Damage uint8 `json:"-" sql:"-"`
    }
    ShipConstructionGroup struct {
        TableName struct{} `json:"-" pg:",discard_unknown_columns"`

        Model *ShipModel `json:"model"`
        ConstructionState *ShipConstructionState `json:"construction_state"`
        Quantity uint `json:"quantity"`
    }
    ShipGroup struct {
        Id uint `json:"id"`
        Name string `json:"name"`
        Type string `json:"type"`
        FrameSlug string `json:"frame"`
        Quantity uint `json:"quantity"`
    }
)

func CreateShips(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)
    planetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := getPlayerPlanet(uint16(planetId), player.Id)

    if planet.Player.Id != player.Id {
        panic(NewHttpException(403, "You do not control this planet", nil))
    }
    SendJsonResponse(w, 201, player.createShips(planet, DecodeJsonRequest(r)))
}

func GetCurrentlyConstructingShips(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)
    planetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := getPlayerPlanet(uint16(planetId), player.Id)

    if planet.Player.Id != player.Id {
        panic(NewHttpException(403, "You do not control this planet", nil))
    }
    SendJsonResponse(w, 200, planet.getCurrentlyConstructingShips())
}

func GetConstructingShips(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)
    planetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := getPlayerPlanet(uint16(planetId), player.Id)

    if planet.Player.Id != player.Id {
        panic(NewHttpException(403, "You do not control this planet", nil))
    }
    SendJsonResponse(w, 200, planet.getConstructingShips())
}

func GetHangarShips(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)
    planetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := getPlayerPlanet(uint16(planetId), player.Id)

    if planet.Player.Id != player.Id {
        panic(NewHttpException(403, "You do not control this planet", nil))
    }
    SendJsonResponse(w, 200, planet.getHangarShips())
}

func GetHangarShipGroups(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)
    planetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := getPlayerPlanet(uint16(planetId), player.Id)

    if planet.Player.Id != player.Id {
        panic(NewHttpException(403, "You do not control this planet", nil))
    }
    SendJsonResponse(w, 200, planet.getHangarShipGroups())
}

func (p *Player) createShips(planet *Planet, data map[string]interface{}) []Ship {
    modelId := uint32(data["model"].(map[string]interface{})["id"].(float64))
    quantity := uint8(data["quantity"].(float64))
    shipModel := p.getShipModel(modelId)

    points := p.payShipCost(shipModel.Price, planet.Storage, quantity)

    constructionState := &ShipConstructionState{
        CurrentPoints: 0,
        Points: points,
    }
    if err := Database.Insert(constructionState); err != nil {
        panic(NewHttpException(500, "Ship construction state could not be created", err))
	}
	ships := make([]Ship, quantity)
    ship := Ship{
        ModelId: shipModel.Id,
        Model: shipModel,
        HangarId: planet.Id,
        Hangar: planet,
        CreatedAt: time.Now(),
        ConstructionState: constructionState,
        ConstructionStateId: constructionState.Id,
    }
    for i := uint8(0); i < quantity; i++ {
		ships[i] = ship
        if err := Database.Insert(&ships[i]); err != nil {
            panic(NewHttpException(500, "Ship could not be created", err))
        }
    }
    planet.Storage.update()
    p.update()
    return ships
}

func (p *Planet) getConstructingShips() []ShipConstructionGroup {
    scg := make([]ShipConstructionGroup, 0)
    if _, err := Database.Query(&scg, `SELECT m.id as model__id, m.name as model__name, m.type as model__type, m.frame_slug as model__frame_slug,
        cs.id as construction_state__id, cs.points as construction_state__points, cs.current_points as construction_state__current_points, COUNT(cs.id) as quantity
        FROM ship__ships s
        INNER JOIN ship__models m ON s.model_id = m.id
        INNER JOIN ship__construction_states cs ON s.construction_state_id = cs.id
        WHERE s.hangar_id = ?
        GROUP BY cs.id, m.id
        ORDER BY cs.id ASC`, p.Id); err != nil {
            panic(NewHttpException(404, "No constructing ship found", err))
    }
    return scg
}

func (p *Planet) getCurrentlyConstructingShips() *ShipConstructionGroup {
    scg := &ShipConstructionGroup{}
    if _, err := Database.
        Query(scg, `SELECT m.id as model__id, m.name as model__name, m.type as model__type, m.frame_slug as model__frame_slug,
        cs.id as construction_state__id, cs.points as construction_state__points, cs.current_points as construction_state__current_points, COUNT(cs.id) as quantity
        FROM ship__ships s
        INNER JOIN ship__models m ON s.model_id = m.id
        INNER JOIN ship__construction_states cs ON s.construction_state_id = cs.id
        WHERE s.hangar_id = ?
        GROUP BY cs.id, m.id
        ORDER BY cs.id ASC
        LIMIT 1`, p.Id); err != nil {
            panic(NewHttpException(404, "No constructing ship found", err))
    }
    return scg
}

func (p *Planet) getHangarShips() []Ship {
    ships := make([]Ship, 0)
    if err := Database.
        Model(&ships).
        Column("Model").
        Where("construction_state_id IS NULL").
        Where("hangar_id = ?", p.Id).
        Select(); err != nil {
        panic(NewHttpException(404, "Planet not found", err))
    }
    return ships
}

func (p *Planet) getHangarShipsByModel(modelId int, quantity int) []Ship {
    ships := make([]Ship, 0)
    if err := Database.
        Model(&ships).
        Column("Hangar", "Fleet").
        Where("construction_state_id IS NULL").
        Where("hangar_id = ?", p.Id).
        Where("model_id = ?", modelId).
        Limit(quantity).
        Select(); err != nil {
        panic(NewHttpException(404, "Planet not found", err))
    }
    return ships
}

func (p *Planet) getHangarShipGroups() []ShipGroup {
    ships := make([]ShipGroup, 0)

    if err := Database.
        Model((*Ship)(nil)).
        ColumnExpr("model.id, model.name, model.type, model.frame_slug, count(*) AS quantity").
        Join("INNER JOIN ship__models as model ON model.id = ship.model_id").
        Group("model.id").
        Where("ship.construction_state_id IS NULL").
        Where("ship.hangar_id = ?", p.Id).
        Select(&ships); err != nil {
            panic(NewHttpException(404, "fleet not found", err))
    }
    return ships
}

func (p *Player) payShipCost(prices []Price, storage *Storage, quantity uint8) uint8 {
    var points uint8
    for _, price := range prices {
        switch price.Type {
            case PriceTypeMoney:
                if !p.updateWallet(-(int32(price.Amount) * int32(quantity))) {
                    panic(NewHttpException(400, "Not enough money", nil))
                }
                break
            case PriceTypePoints:
                points = uint8(price.Amount) * quantity
                break
            case PriceTypeResources:
                amount := uint16(price.Amount) * uint16(quantity)
                if !storage.hasResource(price.Resource, amount) {
                    panic(NewHttpException(400, "Not enough resources", nil))
                }
                storage.storeResource(price.Resource, -int16(amount))
                break
        }
    }
    return points
}

func CheckShipsBuildingState() {
    defer CatchException()

    ships := make([]Ship, 0)
    if err := Database.
        Model(&ships).
        Column("ship.*", "ConstructionState", "Hangar", "Hangar.Settings").
        Order("ship.construction_state_id ASC").
        Where("ship.construction_state_id IS NOT NULL").
        Select(); err != nil {
        panic(NewException("Constructing ships could not be retrieved", err))
    }
    currentPlanetId := uint16(0)
    remainingPoints := uint8(0)
    for _, ship := range ships {
        if currentPlanetId != ship.HangarId {
            currentPlanetId = ship.HangarId
            remainingPoints = ship.Hangar.Settings.MilitaryPoints
        }
        if remainingPoints < 1 {
            continue
        }
        neededPoints := ship.ConstructionState.Points - ship.ConstructionState.CurrentPoints
        if neededPoints <= remainingPoints {
            remainingPoints -= neededPoints
            ship.finishConstruction()
        } else {
            ship.ConstructionState.CurrentPoints += remainingPoints
            ship.ConstructionState.update()
            remainingPoints = 0
        }
    }
}

func (s *Ship) finishConstruction() {
    s.ConstructionStateId = 0
    s.update()
}

func (sc *ShipConstructionState) update() {
	if err := Database.Update(sc); err != nil {
		panic(NewException("Ship Construction State could not be udpated", err))
	}
}

func getShip(id uint16) *Ship {
    ship := &Ship{}
    if err := Database.
        Model(ship).
        Column("ship.*", "Hangar", "Fleet", "Model","Hangar.Player","Fleet.Location", "Fleet.Location.Player","Fleet.Player").
        Where("construction_state_id IS NULL").
        Where("ship.id = ?", id).
        Select(); err != nil {
            panic(NewHttpException(404, "ship not found", err))
    }
    return ship
}

func getShipsByIds(ids []uint16) []*Ship {
    ships := make([]*Ship, 0)
    if err := Database.
        Model(&ships).
        Column("ship.*", "Hangar", "Fleet", "Model", "Hangar.Player", "Fleet.Location", "Fleet.Location.Player", "Fleet.Player").
        Where("construction_state_id IS NULL").
        WhereIn("ship.id IN ?", ids).
        Select(); err != nil {
            panic(NewHttpException(404, "ship not found", err))
    }
    return ships
}

func (s *Ship) update() {
    if err := Database.Update(s); err != nil {
        panic(NewException("ship could not be updated", err))
    }
}

func removeShipsByIds(shipIds []uint32) {
    ships := make([]Ship, 0)
    if _, err := Database.Model(&ships).WhereIn("ship.id IN (?)", shipIds).Delete(); err != nil {
        panic(NewException("Ships could not be removed", err))
    }
}