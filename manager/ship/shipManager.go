package shipManager

import(
    "time"
    "encoding/json"
    "io/ioutil"
    "kalaxia-game-api/database"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/manager"
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
    utils.Scheduler.AddHourlyTask(func () { checkShipBuildingState() })
}

func CreateShip(player *model.Player, planet *model.Planet, data map[string]interface{}) *model.Ship {
    modelId := uint32(data["model"].(map[string]interface{})["id"].(float64))
    quantity := uint8(data["quantity"].(float64))
    shipModel := GetShipModel(player.Id, modelId)

    points := PayShipCost(shipModel.Price, &player.Wallet, planet.Storage, quantity)

    constructionState := model.ShipConstructionState{
        CurrentPoints: 0,
        Points: points,
    }
    ship := model.Ship{
        ModelId: shipModel.Id,
        Model: shipModel,
        HangarId: planet.Id,
        Hangar: planet,
        CreatedAt: time.Now(),
    }
    for i := uint8(0); i < quantity; i++ {
        cs:= constructionState
        s := ship
        if err := database.Connection.Insert(&cs); err != nil {
            panic(exception.NewHttpException(500, "Ship construction state could not be created", err))
        }
        s.ConstructionState = &cs
        s.ConstructionStateId = cs.Id
        if err := database.Connection.Insert(&s); err != nil {
            panic(exception.NewHttpException(500, "Ship could not be created", err))
        }
    }
    manager.UpdatePlanetStorage(planet)
    manager.UpdatePlayer(player)
    return &ship
}

func PayShipCost(prices []model.Price, wallet *uint32, storage *model.Storage, quantity uint8) uint8 {
    var points uint8
    for _, price := range prices {
        switch price.Type {
            case model.PRICE_TYPE_MONEY:
                if price.Amount > *wallet {
                    panic(exception.NewHttpException(400, "Not enough money", nil))
                }
                *wallet -= price.Amount
                break
            case model.PRICE_TYPE_POINTS:
                points = uint8(price.Amount)
                break
            case model.PRICE_TYPE_RESOURCE:
                amount := uint16(price.Amount) * uint16(quantity)
                if !storage.HasResource(price.Resource, amount) {
                    panic(exception.NewHttpException(400, "Not enough resources", nil))
                }
                storage.SubstractResource(price.Resource, amount)
                break
        }
    }
    return points
}

func checkShipBuildingState() {

}
