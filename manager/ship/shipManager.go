package shipManager

import(
    //"time"
    "encoding/json"
    "io/ioutil"
    //"kalaxia-game-api/database"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/model"
    "kalaxia-game-api/utils"
)

var framesData map[string]model.ShipFrame
var modulesData map[string]model.ShipModule

func init() {
    defer utils.CatchException()
    framesDataJSON, err := ioutil.ReadFile("../kalaxia-game-api/resources/ship_frames.json")
    if err != nil {
        panic(exception.NewException("Can't open ship frames configuration file", err))
    }
    if err := json.Unmarshal(framesDataJSON, &framesData); err != nil {
        panic(exception.NewException("Can't read ship frames configuration file", err))
    }
    modulesDataJSON, err := ioutil.ReadFile("../kalaxia-game-api/resources/ship_modules.json")
    if err != nil {
        panic(exception.NewException("Can't open ship modules configuration file", err))
    }
    if err := json.Unmarshal(modulesDataJSON, &modulesData); err != nil {
        panic(exception.NewException("Can't read ship modules configuration file", err))
    }
    //scheduleShipBuildings()
}

// func CreateShip(player *model.Player, planet *model.Planet, data map[string]interface{}) *model.Ship {
//
//
//
//     ship := &model.Ship{
//         Name: data["name"].(string),
//         Model: &model.ShipModel{
//
//         },
//         CreatedAt: time.Now(),
//     }
//
//     if err := database.Connection.Insert(&ship); err != nil {
//       panic(exception.NewHttpException(500, "Ship could not be created", err))
//     }
//     utils.Scheduler.AddTask(buildingDuration, func() {
//         checkShipBuildingState(ship.Id)
//     })
//     return building
//     return ship
// }
//
// func checkShipBuildingState(uint id) {
//
// }
