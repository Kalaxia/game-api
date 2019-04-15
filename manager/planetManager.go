package manager

import(
    "kalaxia-game-api/database"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/model"
    "kalaxia-game-api/utils"
    "time"
    "sort"
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
        Column("planet.*", "System", "Resources", "Settings").
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
        Column("planet.*", "Player", "Player.Faction", "Settings", "Relations", "Relations.Player", "Relations.Player.Faction", "Relations.Faction", "Resources", "System", "Storage").
        Where("planet.id = ?", id).
        Select(); err != nil {
            panic(exception.NewHttpException(404, "Planet not found", err))
    }
    if planet.Player != nil && playerId == planet.Player.Id {
        getPlanetOwnerData(&planet)
    }
    return &planet
}

func getPlanetOwnerData(planet *model.Planet) {
    planet.Buildings, planet.AvailableBuildings = GetPlanetBuildings(planet.Id)
    planet.NbBuildings = 7
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

func getPlanets(offset int, limit int) []*model.Planet {
    var planets []*model.Planet
    if err := database.
        Connection.
        Model(&planets).
        Column("planet.*", "Player", "Buildings", "Buildings.ConstructionState", "Resources", "Storage", "Settings").
        Limit(limit).
        Offset(offset).
        Order("planet.id ASC").
        Select(); err != nil {
            panic(exception.NewException("Planets could not be retrieved", err))
    }
    return planets
}

func GetPlanetsById(ids []uint16) []*model.Planet {
    var planets []*model.Planet
    
    if err := database.
        Connection.
        Model(&planets).
        Column("System", "Player", "Buildings.ConstructionState", "Resources", "Storage", "Settings").
        WhereIn("planet.id IN ?",ids).
        Select(); err != nil {
            panic(exception.NewException("Planets could not be retrieved", err))
    }
    return planets
}

func calculatePlanetProduction(planet *model.Planet, wg *sync.WaitGroup) {
    defer wg.Done()
    defer utils.CatchException()

    calculatePlanetResourcesProduction(planet)
    calculatePointsProduction(planet)
}

func calculatePlanetResourcesProduction(planet *model.Planet) {
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
        if err := database.Connection.Update(planet); err != nil {
            panic(exception.NewException("Planet storage could not be updated", err))
        }
    } else {
        addResourcesToStorage(planet, planet.Storage)
        if err := database.Connection.Update(planet.Storage); err != nil {
            panic(exception.NewException("Planet storage could not be updated", err))
        }
    }
}

func calculatePointsProduction(planet *model.Planet) {
    buildingPoints := planet.Settings.BuildingPoints
    if buildingPoints <= 0 {
        return
    }
    // Sort the buildings by construction date
    constructingBuildings := make(map[string]*model.Building, 0)
    var buildingDates []string
    for _, building := range planet.Buildings {
        if building.Status == model.BuildingStatusConstructing {
            date := building.ConstructionState.BuiltAt.Format(time.RFC3339)
            // we use the date as a key for the constructing buildings map
            constructingBuildings[date] = &building
            buildingDates = append(buildingDates, date)
        }
    }
    // Here we sort the dates
    sort.Strings(buildingDates)
    for _, date := range buildingDates {
        if buildingPoints <= 0 {
            break
        }
        buildingPoints = spendBuildingPoints(constructingBuildings[date], buildingPoints)
    }
}

func addResourcesToStorage(planet *model.Planet, storage *model.Storage) {
    for _, resource := range planet.Resources {
        UpdateStorageResource(storage, resource.Name, int16(resource.Density) * 10)
    }
}

func UpdateStorageResource(storage *model.Storage, resource string, quantity int16) bool {
    var currentStock uint16
    var newStock int16
    var isset bool
    if currentStock, isset = storage.Resources[resource]; !isset {
        currentStock = 0
    }
    if newStock = int16(currentStock) + quantity; newStock > int16(storage.Capacity) {
        newStock = int16(storage.Capacity)
    }
    if newStock < 0 {
        return false
    }
    storage.Resources[resource] = uint16(newStock)
    return true
}

func UpdatePlanetSettings(planet *model.Planet, settings *model.PlanetSettings) {
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
        panic(exception.NewException("Planet settings could not be updated", err))
    }
}

func UpdatePlanetStorage(planet *model.Planet) {
    if err := database.Connection.Update(planet.Storage); err != nil {
        panic(exception.NewException("Planet storage could not be updated", err))
    }
}

func calculatePopulationPoints(planet *model.Planet) uint8 {
    return uint8(math.Ceil(float64(planet.Population / 100000)))
}
