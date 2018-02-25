package base

import(
    "encoding/json"
    "io/ioutil"
)

var buildingPlansData BuildingPlansData
var buildingTypesData BuildingTypesData

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
}
