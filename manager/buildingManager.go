package manager

import(
    "time"
    "encoding/json"
    "io/ioutil"
    "kalaxia-game-api/database"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/model"
    "kalaxia-game-api/utils"
)

var buildingPlansData model.BuildingPlansData

func init() {
    defer utils.CatchException()
    buildingsDataJSON, err := ioutil.ReadFile("../kalaxia-game-api/resources/buildings.json")
    if err != nil {
        panic(exception.NewException("Can't open buildings configuration file", err))
    }
    if err := json.Unmarshal(buildingsDataJSON, &buildingPlansData); err != nil {
        panic(exception.NewException("Can't read buildings configuration file", err))
    }
}

func GetPlanetBuildings(planetId uint16) ([]model.Building, []model.BuildingPlan) {
    buildings := make([]model.Building, 0)
    if err := database.Connection.Model(&buildings).Where("building.planet_id = ?", planetId).Order("id").Column("building.*", "ConstructionState").Select(); err != nil {
        panic(exception.NewHttpException(500, "Something nasty happened", err))
    }
    return buildings, getAvailableBuildings(buildings)
}

func getAvailableBuildings(buildings []model.Building) []model.BuildingPlan {
    availableBuildings := make([]model.BuildingPlan, 0)

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

func CreateBuilding(planet *model.Planet, name string) model.Building {
    buildingPlan, isset := buildingPlansData[name]
    if !isset {
        panic(exception.NewHttpException(400, "unknown building plan", nil))
    }
    constructionState := createConstructionState(planet.Player, buildingPlan)
    building := model.Building{
        Name: name,
        Type: buildingPlan.Type,
        Planet: planet,
        PlanetId: planet.Id,
        Status: model.BuildingStatusConstructing,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        ConstructionState: constructionState,
        ConstructionStateId: constructionState.Id,
    }
    if err := database.Connection.Insert(&building); err != nil {
      panic(exception.NewHttpException(500, "Building could not be created", err))
    }
    planet.AvailableBuildings = getAvailableBuildings(append(planet.Buildings, building))
    return building
}

func CancelBuilding(planet *model.Planet, id uint32) {
    building := model.Building{}
    if err := database.Connection.Model(&building).Column("building.*", "ConstructionState").Where("building.id = ?", id).Select(); err != nil {
        panic(exception.NewHttpException(404, "Building not found", err))
    }
    if building.PlanetId != planet.Id {
        panic(exception.NewHttpException(400, "Building does not belong to the given planet", nil))
    }
    if err := database.Connection.Delete(&building); err != nil {
        panic(exception.NewException("Building could not be removed", err))
    }
}

func createConstructionState(player *model.Player, buildingPlan model.BuildingPlan) *model.ConstructionState {
    points := uint8(0)
    for _, price := range buildingPlan.Price {
        if price.Type == model.PRICE_TYPE_POINTS {
            points = uint8(price.Amount)
        } else if price.Type == model.PRICE_TYPE_MONEY {
            if !UpdatePlayerWallet(player, -int32(price.Amount)) {
                panic(exception.NewHttpException(400, "The player has not enough money", nil))
            }
            UpdatePlayer(player)
        }
    }
    constructionState := &model.ConstructionState {
        Points: points,
        CurrentPoints: 0,
        BuiltAt: time.Now(),
    }
    if err := database.Connection.Insert(constructionState); err != nil {
      panic(exception.NewHttpException(500, "Construction State could not be created", err))
    }
    return constructionState
}

func spendBuildingPoints(building model.Building, buildingPoints uint8) uint8 {
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
    if err := database.Connection.Update(building.ConstructionState); err != nil {
        panic(exception.NewException("Construction State could not be updated", err))
    }
    if building.ConstructionState.CurrentPoints == building.ConstructionState.Points {
        checkConstructionState(building.Id)
    }
    return buildingPoints
}

func checkConstructionState(id uint32) {
    building := &model.Building{}
    if err := database.Connection.Model(building).Column("building.*", "ConstructionState").Where("building.id = ?", id).Select(); err != nil {
        panic(exception.NewException("Building not found", err))
    }
    if building.ConstructionState.CurrentPoints == building.ConstructionState.Points {
        finishConstruction(building)
    }
}

func finishConstruction(building *model.Building) {
    building.Status = model.BuildingStatusOperational
    building.ConstructionStateId = 0
    if err := database.Connection.Update(building); err != nil {
        panic(exception.NewException("Building could not be updated", err))
    }
    if err := database.Connection.Delete(building.ConstructionState); err != nil {
        panic(exception.NewException("Construction State could not be removed", err))
    }
}