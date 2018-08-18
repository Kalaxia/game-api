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
    utils.Scheduler.AddHourlyTask(func () { checkShipsBuildingState() })
}

func CreateShip(player *model.Player, planet *model.Planet, data map[string]interface{}) *model.Ship {
    modelId := uint32(data["model"].(map[string]interface{})["id"].(float64))
    quantity := uint8(data["quantity"].(float64))
    shipModel := GetShipModel(player.Id, modelId)

    points := payShipCost(shipModel.Price, &player.Wallet, planet.Storage, quantity)

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

func GetConstructingShips(planet *model.Planet) []model.Ship {
    ships := make([]model.Ship, 0)
    if err := database.
        Connection.
        Model(&ships).
        Column("ConstructionState", "Model").
        Where("construction_state_id IS NOT NULL").
        Where("hangar_id = ?", planet.Id).
        Select(); err != nil {
        panic(exception.NewHttpException(404, "Planet not found", err))
    }
    return ships
}

func GetHangarShips(planet *model.Planet) []model.Ship {
    ships := make([]model.Ship, 0)
    if err := database.
        Connection.
        Model(&ships).
        Column( "Model").
        Where("construction_state_id IS NULL").
        Where("hangar_id = ?", planet.Id).
        Select(); err != nil {
        panic(exception.NewHttpException(404, "Planet not found", err))
    }
    return ships
}

func payShipCost(prices []model.Price, wallet *uint32, storage *model.Storage, quantity uint8) uint8 {
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

func checkShipsBuildingState() {
    defer utils.CatchException()

    var ships []model.Ship
    if err := database.
        Connection.
        Model(&ships).
        Column("ship.*", "ConstructionState", "Hangar", "Hangar.Settings").
        Order("ship.hangar_id ASC").
        Where("ship.construction_state_id IS NOT NULL").
        Select(); err != nil {
        panic(exception.NewException("Constructing ships could not be retrieved", err))
    }
    currentPlanetId := uint16(0)
    remainingPoints := uint8(0)
    for _, ship := range ships {
        if currentPlanetId != ship.HangarId {
            currentPlanetId = ship.HangarId
            remainingPoints = ship.Hangar.Settings.MilitaryPoints
        }
        if remainingPoints < 1 {
            continue
        }
        neededPoints := ship.ConstructionState.Points - ship.ConstructionState.CurrentPoints
        if neededPoints <= remainingPoints {
            remainingPoints -= neededPoints
            finishShipConstruction(&ship)
        } else {
            ship.ConstructionState.CurrentPoints += remainingPoints
            if err := database.Connection.Update(ship.ConstructionState); err != nil {
                panic(exception.NewException("Ship Construction State could not be udpated", err))
            }
            remainingPoints = 0
        }
    }
}

func finishShipConstruction(ship *model.Ship) {
    ship.ConstructionStateId = 0
    if err := database.Connection.Update(ship); err != nil {
        panic(exception.NewException("Ship could not be updated", err))
    }
    if err := database.Connection.Delete(ship.ConstructionState); err != nil {
        panic(exception.NewException("Ship Construction State could not be removed", err))
    }
}

func GetShip(id uint16) *model.Ship{
    /**
     * Get Ship data ( may return incomplete information).
     *  If the player is the owner of the ship all the data are send
     * 
     */
    
    var ship model.Ship
    if err := database.
        Connection.
        Model(&ship).
        Column("ship.*", "Hangar", "Fleet", "Model","Hangar.Player","Fleet.Location", "Fleet.Location.Player","Fleet.Player").
        Where("construction_state_id IS NULL").
        Where("ship.id = ?", id).
        Select(); err != nil {
            panic(exception.NewHttpException(404, "ship not found", err))
    }
    return &ship
}

func GetShipsByIds(ids []uint16) []*model.Ship{
    /**
     * Get Ships data ( may return incomplete information).
     *  If the player is the owner of the ship all the data are send
     * 
     */
    
    var ships []*model.Ship
    if err := database.
        Connection.
        Model(&ships).
        Column("ship.*", "Hangar", "Fleet", "Model","Hangar.Player","Fleet.Location", "Fleet.Location.Player","Fleet.Player").
        Where("construction_state_id IS NULL").
        WhereIn("ship.id IN ?", ids).
        Select(); err != nil {
            panic(exception.NewHttpException(404, "ship not found", err))
    }
    return ships
}



func UpdateShip(ship *model.Ship){
    
    if err := database.Connection.Update(ship); err != nil {
        panic(exception.NewException("ship could not be updated", err))
    }
    
}

func UpdateShips(ships []*model.Ship){
    
    /*
    if _,err := database.Connection.Model(&ships).Update(); err != nil { //< [Exception]: ship could not be updated; [Error]: ERROR #42804 column "hangar_id" is of type integer but expression is of type text
        panic(exception.NewException("ship could not be updated", err))
    }
    */
    
    for _,ship := range ships {
        UpdateShip(ship);
    }
    
}


func IsShipInSamePositionAsFleet (ship model.Ship, fleet model.Fleet ) bool {
    
    return ( ship.Fleet == nil &&  fleet.Location != nil && fleet.Location.Id ==  ship.Hangar.Id ) || // ship in Hangar and hangar same pos as the fleet
	  (ship.Fleet != nil && fleet.Location != nil && ship.Fleet.Location.Id !=  fleet.Location.Id); // ship in fleet  and both fleet are a the same place
    
}
