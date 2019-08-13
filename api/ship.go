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
        HangarId uint16 `json:"-"`
        Hangar *Planet `json:"hangar"`
        FleetId uint16 `json:"-"`
        Fleet *Fleet `json:"fleet"`
        ModelId uint `json:"-"` 
        Model *ShipModel `json:"model"`
        CreatedAt time.Time `json:"created_at"`
        ConstructionStateId uint32 `json:"-"`
        ConstructionState *PointsProduction `json:"construction_state"`
        // Combat fields
        Damage uint8 `json:"-" sql:"-"`
    }
    ShipConstructionGroup struct {
        TableName struct{} `json:"-" pg:",discard_unknown_columns"`

        Model *ShipModel `json:"model"`
        ConstructionState *PointsProduction `json:"construction_state"`
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
    SendJsonResponse(w, 201, planet.createShips(DecodeJsonRequest(r)))
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

    if planet.Player != nil && planet.Player.Id != player.Id {
        panic(NewHttpException(403, "You do not control this planet", nil))
    }
    SendJsonResponse(w, 200, planet.getHangarShipGroups())
}

func (p *Planet) createShips(data map[string]interface{}) *ShipConstructionGroup {
    modelId := uint32(data["model"].(map[string]interface{})["id"].(float64))
    quantity := uint8(data["quantity"].(float64))
    shipModel := p.Player.getShipModel(modelId)

    constructionState := p.createPointsProduction(p.payPrice(shipModel.Price, quantity))
	ships := make([]Ship, quantity)
    ship := Ship{
        ModelId: shipModel.Id,
        Model: shipModel,
        HangarId: p.Id,
        Hangar: p,
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
    p.Storage.update()
    return &ShipConstructionGroup{
        ConstructionState: constructionState,
        Model: shipModel,
        Quantity: uint(quantity),
    }
}

func (p *Planet) getConstructingShips() []ShipConstructionGroup {
    scg := make([]ShipConstructionGroup, 0)
    if _, err := Database.Query(&scg, `SELECT m.id as model__id, m.name as model__name, m.type as model__type, m.frame_slug as model__frame_slug,
        pp.id as construction_state__id, pp.points as construction_state__points, pp.current_points as construction_state__current_points, COUNT(pp.id) as quantity
        FROM ship__ships s
        INNER JOIN ship__models m ON s.model_id = m.id
        INNER JOIN map__planet_point_productions pp ON s.construction_state_id = pp.id
        WHERE s.hangar_id = ?
        GROUP BY pp.id, m.id
        ORDER BY pp.id ASC`, p.Id); err != nil {
            panic(NewHttpException(404, "No constructing ship found", err))
    }
    return scg
}

func (p *Planet) getCurrentlyConstructingShips() *ShipConstructionGroup {
    scg := &ShipConstructionGroup{}
    if _, err := Database.
        Query(scg, `SELECT m.id as model__id, m.name as model__name, m.type as model__type, m.frame_slug as model__frame_slug,
        pp.id as construction_state__id, pp.points as construction_state__points, pp.current_points as construction_state__current_points, COUNT(pp.id) as quantity
        FROM ship__ships s
        INNER JOIN ship__models m ON s.model_id = m.id
        INNER JOIN map__planet_point_productions pp ON s.construction_state_id = pp.id
        WHERE s.hangar_id = ?
        GROUP BY pp.id, m.id
        ORDER BY pp.id ASC
        LIMIT 1`, p.Id); err != nil {
            panic(NewHttpException(404, "No constructing ship found", err))
    }
    return scg
}

func (p *Planet) getHangarShips() []Ship {
    ships := make([]Ship, 0)
    if err := Database.
        Model(&ships).
        Relation("Model").
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
        Relation("Hangar").
        Relation("Fleet").
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

func (cg *ShipConstructionGroup) finishConstruction() {
    ships := make([]Ship, 0)

    if err := Database.Model(&ships).Where("construction_state_id = ?", cg.ConstructionState.Id).Select(); err != nil {
        panic(NewException("Construction group ships could not be retrieved", err))
    }
    for _, ship := range ships {
        ship.finishConstruction()
    }
    cg.ConstructionState.delete()
}

func (s *Ship) finishConstruction() {
    s.ConstructionStateId = 0
    s.update()
}

func getShip(id uint16) *Ship {
    ship := &Ship{}
    if err := Database.
        Model(ship).
        Relation("Model").
        Relation("Hangar.Player").
        Relation("Fleet.Location.Player").
        Relation("Fleet.Player").
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
        Relation("Model").
        Relation("Hangar.Player").
        Relation("Fleet.Location.Player").
        Relation("Fleet.Player").
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