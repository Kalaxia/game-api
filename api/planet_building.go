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

type(
	Building struct {
		TableName struct{} `json:"-" sql:"map__planet_buildings"`

		Id uint32 `json:"id"`
		Name string `json:"name"`
		Type string `json:"type" sql:"type"`
		Planet *Planet `json:"planet"`
		PlanetId uint16 `json:"-"`
		ConstructionState *ConstructionState `json:"construction_state"`
		ConstructionStateId uint32 `json:"-"`
		Status string `json:"status"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		
	}
	BuildingPlan struct {
		Name string `json:"name"`
		ParentName string `json:"parent"`
		Type string `json:"type"`
		Picture string `json:"picture"`
		Price []Price `json:"price"`
	}
	BuildingPlansData map[string]BuildingPlan

	ConstructionState struct {
		TableName struct{} `json:"-" sql:"map__planet_construction_states"`

		Id uint32 `json:"id"`
		BuiltAt time.Time `json:"built_at"`
		CurrentPoints uint8 `json:"current_points" sql:",notnull"`
		Points uint8 `json:"points"`
	}
)

func InitPlanetConstructions() {
    defer CatchException()
    buildingsDataJSON, err := ioutil.ReadFile("../kalaxia-game-api/resources/buildings.json")
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
    planet := getPlayerPlanet(uint16(id), player.Id)

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
    planet := getPlayerPlanet(uint16(planetId), player.Id)

    if uint16(planetId) != planet.Id {
        panic(NewHttpException(403, "Forbidden", nil))
    }
    planet.cancelBuilding(uint32(buildingId))

    w.WriteHeader(204)
    w.Write([]byte(""))
}

func (p *Planet) getBuildings() ([]Building, []BuildingPlan) {
    buildings := make([]Building, 0)
    if err := Database.Model(&buildings).Where("building.planet_id = ?", p.Id).Order("id").Column("building.*", "ConstructionState").Select(); err != nil {
        panic(NewHttpException(500, "buildings.internal_error", err))
    }
    return buildings, getAvailableBuildings(buildings)
}

func getAvailableBuildings(buildings []Building) []BuildingPlan {
    availableBuildings := make([]BuildingPlan, 0)

    for buildingName, buildingPlan := range buildingPlansData {
        existing := false
        for _, building := range buildings {
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

func (p *Planet) createBuilding(name string) Building {
    buildingPlan, isset := buildingPlansData[name]
    if !isset {
        panic(NewHttpException(400, "unknown building plan", nil))
    }
    constructionState := createConstructionState(p.Player, buildingPlan)
    building := Building{
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
    p.AvailableBuildings = getAvailableBuildings(append(p.Buildings, building))
    return building
}

func (p *Planet) cancelBuilding(id uint32) {
    building := &Building{}
    if err := Database.Model(building).Column("building.*", "ConstructionState").Where("building.id = ?", id).Select(); err != nil {
        panic(NewHttpException(404, "Building not found", err))
    }
    if building.PlanetId != p.Id {
        panic(NewHttpException(400, "Building does not belong to the given planet", nil))
    }
    if err := Database.Delete(building); err != nil {
        panic(NewException("Building could not be removed", err))
    }
}

func createConstructionState(player *Player, buildingPlan BuildingPlan) *ConstructionState {
    points := uint8(0)
    for _, price := range buildingPlan.Price {
        if price.Type == PriceTypePoints {
            points = uint8(price.Amount)
        } else if price.Type == PriceTypeMoney {
            if !player.updateWallet(-int32(price.Amount)) {
                panic(NewHttpException(400, "The player has not enough money", nil))
            }
            player.update()
        }
    }
    constructionState := &ConstructionState {
        Points: points,
        CurrentPoints: 0,
        BuiltAt: time.Now(),
    }
    if err := Database.Insert(constructionState); err != nil {
      panic(NewHttpException(500, "Construction State could not be created", err))
    }
    return constructionState
}

func spendBuildingPoints(building Building, buildingPoints uint8) uint8 {
    missingPoints := building.ConstructionState.Points - building.ConstructionState.CurrentPoints
    if missingPoints == 0 {
        checkConstructionState(building.Id)
        return buildingPoints
    }
    if missingPoints > buildingPoints {
        building.ConstructionState.CurrentPoints += buildingPoints
        buildingPoints = 0
    } else {
        building.ConstructionState.CurrentPoints += missingPoints
        buildingPoints -= missingPoints
    }
    if err := Database.Update(building.ConstructionState); err != nil {
        panic(NewException("Construction State could not be updated", err))
    }
    if building.ConstructionState.CurrentPoints == building.ConstructionState.Points {
        checkConstructionState(building.Id)
    }
    return buildingPoints
}

func checkConstructionState(id uint32) {
    building := &Building{}
    if err := Database.Model(building).Column("building.*", "ConstructionState").Where("building.id = ?", id).Select(); err != nil {
        panic(NewException("Building not found", err))
    }
    if building.ConstructionState.CurrentPoints == building.ConstructionState.Points {
        building.finishConstruction()
    }
}

func (b *Building) finishConstruction() {
    b.Status = BuildingStatusOperational
    b.ConstructionStateId = 0
    if err := Database.Update(b); err != nil {
        panic(NewException("Building could not be updated", err))
    }
    if err := Database.Delete(b.ConstructionState); err != nil {
        panic(NewException("Construction State could not be removed", err))
    }
}