package manager

import(
    "fmt"
    "encoding/json"
    "io/ioutil"
    "kalaxia-game-api/database"
    "kalaxia-game-api/model/building"
)

var buildingPlansData model.BuildingPlansData
var buildingTypesData model.BuildingTypesData

func init() {
  buildingsDataJSON, err := ioutil.ReadFile("/go/src/kalaxia-game-api/resources/buildings.json")
  if err != nil {
    panic(err)
  }
  buildingTypesJSON, err := ioutil.ReadFile("/go/src/kalaxia-game-api/resources/building_types.json")
  if err != nil {
    panic(err)
  }
  if err := json.Unmarshal(buildingsDataJSON, &buildingPlansData); err != nil {
    panic(err)
  }
  if err := json.Unmarshal(buildingTypesJSON, &buildingTypesData); err != nil {
    panic(err)
  }
  fmt.Println(buildingPlansData)
}

func GetPlanetBuildings(planetId uint16) ([]model.Building, []model.BuildingPlan) {
  buildings := make([]model.Building, 0)
  _ = database.Connection.Model(&buildings).Where("planet_id = ?", planetId).Select()
  return buildings, getAvailableBuildings(buildings)
}

func getAvailableBuildings(buildings []model.Building) []model.BuildingPlan {
  availableBuildings := make([]model.BuildingPlan, 0)

  for buildingName, buildingPlan := range buildingPlansData {
      if len(buildingPlan.ParentName) == 0 {
          buildingPlan.Name = buildingName
          availableBuildings = append(availableBuildings, buildingPlan)
      }
  }

  return availableBuildings;
}
