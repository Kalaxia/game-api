package manager

import(
    "errors"
    "time"
    "encoding/json"
    "io/ioutil"
    "kalaxia-game-api/database"
    "kalaxia-game-api/model"
    "kalaxia-game-api/utils"
)

var buildingPlansData model.BuildingPlansData

func init() {
    buildingsDataJSON, err := ioutil.ReadFile("../kalaxia-game-api/resources/buildings.json")
    if err != nil {
        panic(err)
    }
    if err := json.Unmarshal(buildingsDataJSON, &buildingPlansData); err != nil {
        panic(err)
    }
    scheduleConstructions()
}

func GetPlanetBuildings(planetId uint16) ([]model.Building, []model.BuildingPlan) {
    buildings := make([]model.Building, 0)
    _ = database.Connection.Model(&buildings).Where("planet_id = ?", planetId).Select()
    return buildings, getAvailableBuildings(buildings)
}

// FIXME some building must need other buildings
func getAvailableBuildings(buildings []model.Building) []model.BuildingPlan {
    availableBuildings := make([]model.BuildingPlan, 0)

    for buildingName, buildingPlan := range buildingPlansData {
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
        panic(errors.New("unknown building plan"))
    }
    building := model.Building{
        Name: name,
        Type: buildingPlan.Type,
        Planet: planet,
        PlanetId: planet.Id,
        Status: model.BuildingStatusConstructing,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        BuiltAt: time.Now().Add(time.Second * time.Duration(buildingPlan.Duration)),
    }
    if err := database.Connection.Insert(&building); err != nil {
      panic(err)
    }
    utils.Scheduler.AddTask(buildingPlan.Duration, func() {
        FinishConstruction(building.Id)
    })
    return building
}

func FinishConstruction(id uint32) {
    building := model.Building{
        Id: id,
    }
    if err := database.Connection.Select(&building); err != nil {
        panic(err)
    }
    building.Status = model.BuildingStatusOperational
    if err := database.Connection.Update(&building); err != nil {
        panic(err)
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
                    FinishConstruction(building.Id)
                })
            } else {
                FinishConstruction(building.Id)
            }
        }(building)
    }
}

func getConstructingBuildings() []model.Building {
    buildings := make([]model.Building, 0)
    _ = database.
        Connection.
        Model(&buildings).
        Where("status = ?", model.BuildingStatusConstructing).
        Select()
    return buildings
}
