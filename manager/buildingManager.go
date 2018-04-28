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
    buildingsDataJSON, err := ioutil.ReadFile("../kalaxia-game-api/resources/buildings.json")
    if err != nil {
        panic(exception.NewException("Can't open buildings configuration file", err))
    }
    if err := json.Unmarshal(buildingsDataJSON, &buildingPlansData); err != nil {
        panic(exception.NewException("Can't read buildings configuration file", err))
    }
    scheduleConstructions()
}

func GetPlanetBuildings(planetId uint16) ([]model.Building, []model.BuildingPlan) {
    buildings := make([]model.Building, 0)
    if err := database.Connection.Model(&buildings).Where("building.planet_id = ?", planetId).Column("building.*", "ConstructionState").Select(); err != nil {
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
    constructionState := createConstructionState(buildingPlan)
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
    utils.Scheduler.AddTask(buildingPlan.Duration, func() {
        checkConstructionState(building.Id)
    })
    planet.AvailableBuildings = getAvailableBuildings(append(planet.Buildings, building))
    return building
}

func createConstructionState(buildingPlan model.BuildingPlan) *model.ConstructionState {
    points := uint8(0)
    for _, price := range buildingPlan.Price {
        if price.Type == "points" {
            points = uint8(price.Amount)
        }
    }
    constructionState := &model.ConstructionState {
        Points: points,
        CurrentPoints: 0,
        BuiltAt: time.Now().Add(time.Second * time.Duration(buildingPlan.Duration)),
    }
    if err := database.Connection.Insert(constructionState); err != nil {
      panic(exception.NewHttpException(500, "Construction State could not be created", err))
    }
    return constructionState
}

func spendBuildingPoints(building model.Building, buildingPoints uint8) uint8 {
    missingPoints := building.ConstructionState.Points - building.ConstructionState.CurrentPoints
    if missingPoints == 0 {
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
    if time.Now().After(building.ConstructionState.BuiltAt) &&
       building.ConstructionState.CurrentPoints == building.ConstructionState.Points {
        finishConstruction(building)
    }
}

func finishConstruction(building *model.Building) {
    building.Status = model.BuildingStatusOperational
    if err := database.Connection.Update(building); err != nil {
        panic(exception.NewException("Building could not be updated", err))
    }
    if err := database.Connection.Delete(building.ConstructionState); err != nil {
        panic(exception.NewException("Construction State could not be removed", err))
    }
}

func scheduleConstructions() {
    constructingBuildings := getConstructingBuildings()
    now := time.Now()

    for _, building := range constructingBuildings {
        go func(building model.Building) {
            plan := buildingPlansData[building.Name]
            endedAt := building.CreatedAt.Add(time.Second * time.Duration(plan.Duration))
            if endedAt.After(now) {
                utils.Scheduler.AddTask(uint(time.Until(endedAt).Seconds()), func() {
                    checkConstructionState(building.Id)
                })
            } else {
                checkConstructionState(building.Id)
            }
        }(building)
    }
}

func getConstructingBuildings() []model.Building {
    buildings := make([]model.Building, 0)
    _ = database.
        Connection.
        Model(&buildings).
        Where("building.status = ?", model.BuildingStatusConstructing).
        Select("building.*", "building.ConstructionState")
    return buildings
}
