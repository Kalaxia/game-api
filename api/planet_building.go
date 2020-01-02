package api

import(
    "encoding/json"
    "io/ioutil"
    "net/http"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
	"strconv"
	"time"
)

var buildingPlansData BuildingPlansData

const BuildingStatusConstructing = "constructing"
const BuildingStatusOperational = "operational"
const BuildingStatusDestroying = "destroying"

const BuildingTypeTerritorialControl = "territorial-control"
const BuildingTypeTrade = "trade"
const BuildingTypeShipyard = "shipyard"
const BuildingTypeTechno = "techno"
const BuildingTypeResource = "resource"

type(
	Building struct {
		tableName struct{} `pg:"map__planet_buildings"`

		Id uint32 `json:"id"`
		Name string `json:"name"`
		Type string `json:"type" pg:"type"`
		Planet *Planet `json:"planet"`
		PlanetId uint16 `json:"-"`
		ConstructionState *PointsProduction `json:"construction_state"`
        ConstructionStateId uint32 `json:"-"`
        Compartments []*BuildingCompartment `json:"compartments"`
		Status string `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		
    }
    BuildingCompartment struct {
        tableName struct{} `pg:"map__planet_building_compartments"`

        Id uint32 `json:"id"`
        Name string `json:"name"`
        BuildingId uint32 `json:"-"`
        Building *Building `json:"building"`
        ConstructionStateId uint32 `json:"-"`
        ConstructionState *PointsProduction `json:"construction_state"`
        Status string `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
    }

	BuildingPlan struct {
		Name string `json:"name"`
		ParentName string `json:"parent"`
        Type string `json:"type"`
        Resources []string `json:"resources"`
        Picture string `json:"picture"`
        Compartments []BuildingCompartmentPlan `json:"compartments"`
		Price []Price `json:"price"`
	}
    BuildingPlansData map[string]BuildingPlan
    
    BuildingCompartmentPlan struct {
        Name string `json:"name"`
        Modifiers []Modifier `json:"modifiers"`
        Price []Price `json:"price"`
    }
)

func InitPlanetConstructions() {
    defer CatchException(nil)
    buildingsDataJSON, err := ioutil.ReadFile("/go/src/kalaxia-game-api/resources/buildings.json")
    if err != nil {
        panic(NewException("Can't open buildings configuration file", err))
    }
    if err := json.Unmarshal(buildingsDataJSON, &buildingPlansData); err != nil {
        panic(NewException("Can't read buildings configuration file", err))
    }
}

func CreateBuilding(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    player := context.Get(r, "player").(*Player)

    id, _ := strconv.ParseUint(vars["id"], 10, 16)
    planet := player.getPlanet(uint16(id))

    if uint16(id) != planet.Id {
        panic(NewHttpException(403, "Forbidden", nil))
    }
    data := DecodeJsonRequest(r)
    SendJsonResponse(w, 201, planet.createBuilding(data["name"].(string)))
}

func CancelBuilding(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    player := context.Get(r, "player").(*Player)

    planetId, _ := strconv.ParseUint(vars["planet-id"], 10, 16)
    buildingId, _ := strconv.ParseUint(vars["building-id"], 10, 16)
    planet := player.getPlanet(uint16(planetId))

    if uint16(planetId) != planet.Id {
        panic(NewHttpException(403, "Forbidden", nil))
    }
    planet.cancelBuilding(uint32(buildingId))

    w.WriteHeader(204)
    w.Write([]byte(""))
}

func CreateBuildingCompartment(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    player := context.Get(r, "player").(*Player)

    planetId, _ := strconv.ParseUint(vars["planet-id"], 10, 16)
    buildingId, _ := strconv.ParseUint(vars["building-id"], 10, 16)
    planet := player.getPlanet(uint16(planetId))

    data := DecodeJsonRequest(r)

    SendJsonResponse(w, 201, planet.getBuilding(uint32(buildingId)).createCompartment(data["name"].(string)))
}

func (p *Planet) getBuildings() ([]*Building, []BuildingPlan) {
    p.Buildings = make([]*Building, 0)
    if err := Database.Model(&p.Buildings).Where("building.planet_id = ?", p.Id).Order("id").Relation("ConstructionState").Relation("Compartments.ConstructionState").Select(); err != nil {
        panic(NewHttpException(500, "buildings.internal_error", err))
    }
    return p.Buildings, p.getAvailableBuildings()
}

func (p *Planet) getBuilding(id uint32) *Building {
    for _, b := range p.Buildings {
        if b.Id == id {
            b.Planet = p
            return b
        }
    }
    panic(NewHttpException(404, "planets.buildings.not_found", nil))
}

func (p *Planet) getAvailableBuildings() []BuildingPlan {
    availableBuildings := make([]BuildingPlan, 0)

    for buildingName, buildingPlan := range buildingPlansData {
        existing := false
        for _, building := range p.Buildings {
            if building.Name == buildingName {
                existing = true
            }
        }
        if existing == true {
            continue
        }
        if len(buildingPlan.ParentName) == 0 {
            buildingPlan.Name = buildingName
            availableBuildings = append(availableBuildings, buildingPlan)
        }
    }
    return availableBuildings
}

func (p *Planet) createBuilding(name string) *Building {
    buildingPlan, isset := buildingPlansData[name]
    if !isset {
        panic(NewHttpException(400, "unknown building plan", nil))
    }
    constructionState := p.createPointsProduction(p.payPrice(buildingPlan.Price, 1))
    building := &Building{
        Name: name,
        Type: buildingPlan.Type,
        Planet: p,
        PlanetId: p.Id,
        Status: BuildingStatusConstructing,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        ConstructionState: constructionState,
        ConstructionStateId: constructionState.Id,
    }
    if err := Database.Insert(&building); err != nil {
      panic(NewHttpException(500, "Building could not be created", err))
    }
    p.Buildings = append(p.Buildings, building)
    p.AvailableBuildings = p.getAvailableBuildings()
    return building
}

func (p *Planet) cancelBuilding(id uint32) {
    building := &Building{}
    if err := Database.Model(building).Relation("ConstructionState").Where("building.id = ?", id).Select(); err != nil {
        panic(NewHttpException(404, "Building not found", err))
    }
    if building.PlanetId != p.Id {
        panic(NewHttpException(400, "Building does not belong to the given planet", nil))
    }
    if err := Database.Delete(building); err != nil {
        panic(NewException("Building could not be removed", err))
    }
}

func (b *Building) createCompartment(name string) *BuildingCompartment {
    compartmentPlan := b.getCompartmentPlan(name)
    if compartmentPlan == nil {
        panic(NewHttpException(400, "planets.buildings.compartments.invalid", nil))
    }
    cs := b.Planet.createPointsProduction(b.Planet.payPrice(compartmentPlan.Price, 1))
    compartment := &BuildingCompartment{
        Name: name,
        Status: BuildingStatusConstructing,
        Building: b,
        BuildingId: b.Id,
        ConstructionState: cs,
        ConstructionStateId: cs.Id,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    if err := Database.Insert(compartment); err != nil {
        panic(NewException("Could not create building compartment", err))
    }
    return compartment
}

func (b *Building) getCompartmentPlan(name string) *BuildingCompartmentPlan {
    for _, plan := range buildingPlansData[b.Name].Compartments {
        if plan.Name == name {
            return &plan
        }
    }
    return nil
}

func (b *Building) update() {
    b.UpdatedAt = time.Now()
    if err := Database.Update(b); err != nil {
        panic(NewException("Building could not be updated", err))
    }
}

func (b *Building) finishConstruction() {
    b.Status = BuildingStatusOperational
    b.ConstructionStateId = 0

    b.Planet.Player.notify(NotificationTypeBuilding, "planet.buildings.notifications.construction_success", map[string]interface{}{
        "building_name": b.Name,
        "planet_id": b.Planet.Id,
        "planet_name": b.Planet.Name,
    })

    b.update()
    b.ConstructionState.delete()
    b.apply()
}

func (b *Building) apply() {
    switch (b.Type) {
        case BuildingTypeTerritorialControl:
            planet := getPlanet(b.PlanetId)
            planet.createTerritory()
            break;
    }
}

func (c *BuildingCompartment) finishConstruction() {
    c.Status = BuildingStatusOperational
    c.ConstructionStateId = 0

    c.update()
    c.Building.update()

    c.Building.Planet.Player.notify(NotificationTypeBuilding, "planet.buildings.notifications.compartment_success", map[string]interface{}{
        "compartment_name": c.Name,
        "building_name": c.Building.Name,
        "planet_id": c.Building.Planet.Id,
        "planet_name": c.Building.Planet.Name,
    })
    
    c.ConstructionState.delete()
}

func (c *BuildingCompartment) update() {
    c.UpdatedAt = time.Now()
    if err := Database.Update(c); err != nil {
        panic(NewException("Building Compartment could not be updated", err))
    }
}