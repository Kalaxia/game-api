package api

import(
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "net/http"
    "strconv"
)

type(
    ShipConstructionGroup struct {
        tableName struct{} `json:"-" pg:"ship__construction_groups"`

        Id uint32 `json:"id"`
        LocationId uint16 `json:"-"`
        Location *Planet `json:"location"`
        ModelId uint `json:"-"`
        Model *ShipModel `json:"model"`
        ConstructionStateId uint32 `json:"-"`
        ConstructionState *PointsProduction `json:"construction_state"`
        Quantity uint8 `json:"quantity"`
    }
)

func CreateShips(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)
    planetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := player.getPlanet(uint16(planetId))

    if planet.Player.Id != player.Id {
        panic(NewHttpException(403, "You do not control this planet", nil))
    }
    SendJsonResponse(w, 201, planet.createShips(DecodeJsonRequest(r)))
}

func GetConstructingShips(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)
    planetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := player.getPlanet(uint16(planetId))

    if planet.Player.Id != player.Id {
        panic(NewHttpException(403, "You do not control this planet", nil))
    }
    SendJsonResponse(w, 200, planet.getConstructingShips())
}

func GetHangarGroups(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)
    planetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := player.getPlanet(uint16(planetId))

    if planet.Player.Id != player.Id {
        panic(NewHttpException(403, "You do not control this planet", nil))
    }
    SendJsonResponse(w, 200, planet.getHangarGroups())
}

func (p *Planet) createShips(data map[string]interface{}) *ShipConstructionGroup {
    modelId := uint32(data["model"].(map[string]interface{})["id"].(float64))
    quantity := uint8(data["quantity"].(float64))
    shipModel := p.Player.getShipModel(modelId)

    constructionState := p.createPointsProduction(p.payPrice(shipModel.Price, quantity))
    cg := &ShipConstructionGroup{
        Location: p,
        LocationId: p.Id,
        ConstructionState: constructionState,
        ConstructionStateId: constructionState.Id,
        Model: shipModel,
        ModelId: shipModel.Id,
        Quantity: quantity,
    }
    if err := Database.Insert(cg); err != nil {
        panic(NewException("Could not create ship construction group", err))
    }
    p.Storage.update()
    return cg
}

func (p *Planet) getConstructingShips() []*ShipConstructionGroup {
    groups := make([]*ShipConstructionGroup, 0)
    if err := Database.
        Model(&groups).
        Relation("ConstructionState").
        Relation("Model").
        Where("location_id = ?", p.Id).
        Order("construction_state_id ASC").
        Select(); err != nil {
            panic(NewHttpException(404, "No constructing ship found", err))
    }
    for _, scg := range groups {
        scg.Location = p
        scg.LocationId = p.Id
    }
    return groups
}

func (cg *ShipConstructionGroup) finishConstruction() {
    cg.Location.addShips(cg.Model, cg.Quantity)
    cg.ConstructionState.delete()
}