package base

import(
    "errors"
    "time"
    "encoding/json"
    "kalaxia-game-api/database"
    "kalaxia-game-api/galaxy"
    "kalaxia-game-api/utils"
)

func GetPlanetBuildings(planetId uint16) ([]Building, []BuildingPlan) {
    buildings := make([]Building, 0)
    database.Connection.Model(&buildings).Where("planet_id = ?", planetId).Select()
    return buildings, getAvailableBuildings(buildings)
}

func getAvailableBuildings(buildings []Building) []BuildingPlan {
  availableBuildings := make([]BuildingPlan, 0)

  for buildingName, buildingPlan := range buildingPlansData {
      if len(buildingPlan.ParentName) == 0 {
          buildingPlan.Name = buildingName
          availableBuildings = append(availableBuildings, buildingPlan)
      }
  }
  return availableBuildings
}

func CreateBuilding(planet *galaxy.Planet, name string) Building {
    buildingPlan, isset := buildingPlansData[name]
    if !isset {
        panic(errors.New("Unknown building plan"))
    }
    buildingType, isset := buildingTypesData[buildingPlan.Type]
    if !isset {
        panic(errors.New("Unknown building type"))
    }
    building := Building{
        Name: name,
        Type: &BuildingType{
            Name: buildingPlan.Type,
            Color: buildingType.Color,
        },
        TypeName: buildingPlan.Type,
        Planet: planet,
        PlanetId: planet.Id,
        Status: BUILDING_STATUS_CONSTRUCTING,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
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
    building := Building{
        Id: id,
    }
    if err := database.Connection.Select(&building); err != nil {
        panic(err)
    }
    building.Status = BUILDING_STATUS_OPERATIONAL
    if err := database.Connection.Update(&building); err != nil {
        panic(err)
    }
}

func scheduleConstructions() {
    constructingBuildings := getConstructingBuildings()
    now := time.Now()

    for _, building := range constructingBuildings {
        go func(building Building) {
            plan := buildingPlansData[building.Name]
            endedAt := building.CreatedAt.Add(time.Second * time.Duration(plan.Duration))
            if endedAt.After(now) {
                utils.Scheduler.AddTask(uint(time.Until(endedAt).Seconds()), func() {
                    b.FinishConstruction(building.Id)
                })
            } else {
                b.FinishConstruction(building.Id)
            }
        }(building)
    }
}

func getConstructingBuildings() []Building {
    buildings := make([]Building, 0)
    database.Connection.Model(&buildings).Where("status = ?", BUILDING_STATUS_CONSTRUCTING).Select()
    return buildings
}
