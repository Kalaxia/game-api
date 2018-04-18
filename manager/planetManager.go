package manager

import(
    "kalaxia-game-api/database"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/model"
    "kalaxia-game-api/utils"
    "sync"
    "math"
)

func init() {
    utils.Scheduler.AddHourlyTask(func () { CalculatePlanetsProductions() })
}

func GetSystemPlanets(id uint16) []model.Planet {
    var planets []model.Planet
    if err := database.
        Connection.
        Model(&planets).
        Column("planet.*", "Orbit", "Player", "Player.Faction").
        Where("planet.system_id = ?", id).
        Select(); err != nil {
            panic(exception.NewHttpException(404, "System not found", err))
    }
    return planets
}

func GetPlayerPlanets(id uint16) []model.Planet {
    var planets []model.Planet
    if err := database.
        Connection.
        Model(&planets).
        Column("planet.*", "Resources").
        Where("planet.player_id = ?", id).
        Select(); err != nil {
            panic(exception.NewHttpException(404, "Player not found", err))
    }
    return planets
}

func GetPlanet(id uint16, playerId uint16) *model.Planet {
    var planet model.Planet
    if err := database.
        Connection.
        Model(&planet).
        Column("planet.*", "Player", "Settings", "Relations", "Relations.Player", "Relations.Player.Faction", "Relations.Faction", "Resources", "System", "Storage").
        Where("planet.id = ?", id).
        Select(); err != nil {
            return nil
    }
    if planet.Player != nil && playerId == planet.Player.Id {
        getPlanetOwnerData(&planet)
    }
    return &planet
}

func getPlanetOwnerData(planet *model.Planet) {
    planet.Buildings, planet.AvailableBuildings = GetPlanetBuildings(planet.Id)
    planet.NbBuildings = 3
}

func CalculatePlanetsProductions() {
    nbPlanets, _ := database.Connection.Model(&model.Planet{}).Count()
    limit := 20

    var wg sync.WaitGroup

    for offset := 0; offset < nbPlanets; offset += limit {
        planets := getPlanets(offset, limit)

        for _, planet := range planets {
            wg.Add(1)
            go calculatePlanetProduction(planet, &wg)
        }
        wg.Wait()
    }
}

func getPlanets(offset int, limit int) []model.Planet {
    var planets []model.Planet
    if err := database.
        Connection.
        Model(&planets).
        Column("planet.*", "Player", "Buildings", "Resources", "Storage").
        Limit(limit).
        Offset(offset).
        Select(); err != nil {
            panic(exception.NewException("Planets could not be retrieved", err))
    }
    return planets
}

func calculatePlanetProduction(planet model.Planet, wg *sync.WaitGroup) {
    defer wg.Done()
    defer utils.CatchException()
    if planet.Storage == nil {
        storage := &model.Storage{
            Capacity: 5000,
            Resources: make(map[string]uint16, 0),
        }
        addResourcesToStorage(planet, storage)
        if err := database.Connection.Insert(storage); err != nil {
            panic(exception.NewException("Storage could not be created", err))
        }
        planet.Storage = storage
        planet.StorageId = storage.Id
        if err := database.Connection.Update(&planet); err != nil {
            panic(exception.NewException("Planet storage could not be updated", err))
        }
    } else {
        addResourcesToStorage(planet, planet.Storage)
        if err := database.Connection.Update(planet.Storage); err != nil {
            panic(exception.NewException("Planet storage could not be updated", err))
        }
    }
}

func addResourcesToStorage(planet model.Planet, storage *model.Storage) {
    for _, resource := range planet.Resources {
        var currentStock uint16
        var newStock uint16
        var isset bool
        if currentStock, isset = storage.Resources[resource.Name]; !isset {
            currentStock = 0
        }
        if newStock = currentStock + uint16(resource.Density) * 10; newStock > storage.Capacity {
            newStock = storage.Capacity
        }
        storage.Resources[resource.Name] = newStock
    }
}

func UpdatePlanetSettings(planet *model.Planet, settings *model.PlanetSettings) error {
    if settings.ServicesPoints +
    settings.BuildingPoints +
    settings.MilitaryPoints +
    settings.ResearchPoints > calculatePopulationPoints(planet) {
        panic(exception.NewHttpException(400, "Not enough population points", nil))
    }
    planet.Settings.ServicesPoints = settings.ServicesPoints
    planet.Settings.BuildingPoints = settings.BuildingPoints
    planet.Settings.MilitaryPoints = settings.MilitaryPoints
    planet.Settings.ResearchPoints = settings.ResearchPoints

    if err := database.Connection.Update(planet.Settings); err != nil {
        panic(exception.NewException("Planet settings could not be updated", nil))
    }
    return nil
}

func calculatePopulationPoints(planet *model.Planet) uint8 {
    return uint8(math.Ceil(float64(planet.Population / 100000)))
}
