package api

import(
    "time"
    "sort"
    "sync"
)

func CalculatePlanetsProductions() {
    nbPlanets, _ := Database.Model(&Planet{}).Count()
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

func getPlanets(offset int, limit int) []*Planet {
    planets := make([]*Planet, 0)
    if err := Database.
        Model(&planets).
        Column("planet.*", "Player", "Buildings", "Buildings.ConstructionState", "Resources", "Storage", "Settings").
        Limit(limit).
        Offset(offset).
        Order("planet.id ASC").
        Select(); err != nil {
            panic(NewException("Planets could not be retrieved", err))
    }
    return planets
}

func getPlanetsById(ids []uint16) []*Planet {
    planets := make([]*Planet, 0)

    if err := Database.
        Model(&planets).
        Column("System", "Player", "Buildings.ConstructionState", "Resources", "Storage", "Settings").
        WhereIn("planet.id IN (?)", ids).
        Select(); err != nil {
            panic(NewException("Planets could not be retrieved", err))
    }
    return planets
}

func calculatePlanetProduction(planet *Planet, wg *sync.WaitGroup) {
    defer wg.Done()
    defer CatchException(nil)

	planet.produceResources()
	planet.producePoints()
}

func (p *Planet) produceResources() {
    if p.Storage == nil {
        storage := &Storage{
            Capacity: 5000,
            Resources: make(map[string]uint16, 0),
        }
        storage.storeResourceProduction(p)
        if err := Database.Insert(storage); err != nil {
            panic(NewException("Storage could not be created", err))
        }
        p.Storage = storage
        p.StorageId = storage.Id
        p.update()
    } else {
        p.Storage.storeResourceProduction(p)
        p.Storage.update()
    }
}

func (p *Planet) producePoints() {
    p.produceBuildingPoints()
    p.produceMilitaryPoints()
}

func (p *Planet) produceBuildingPoints() {
    buildingPoints := p.Settings.BuildingPoints
    if buildingPoints <= 0 || len(p.Buildings) == 0 {
        return
    }
    // Sort the buildings by construction date
    constructingBuildings := make(map[string]Building, 0)
    var buildingDates []string
    for _, building := range p.Buildings {
        if building.Status == BuildingStatusConstructing {
            date := building.ConstructionState.BuiltAt.Format(time.RFC3339)
            // we use the date as a key for the constructing buildings map
            constructingBuildings[date] = building
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
