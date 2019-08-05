package api

import(
    "time"
    "sort"
    "sync"
)

type(
    PointsProduction struct {
		TableName struct{} `json:"-" sql:"map__planet_point_productions"`

		Id uint32 `json:"id"`
		CurrentPoints uint8 `json:"current_points" sql:",notnull"`
		Points uint8 `json:"points"`
    }
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

func calculatePlanetProduction(planet *Planet, wg *sync.WaitGroup) {
    defer wg.Done()
    defer CatchException(nil)

	planet.produceResources()
	planet.producePoints()
}

func (p *Planet) produceResources() {
    p.Storage.storeResourceProduction(p)
    p.Storage.update()
}

func (p *Planet) producePoints() {
    p.produceBuildingPoints()
    p.produceMilitaryPoints()
}

func (p *Planet) createPointsProduction(points uint8) *PointsProduction {
    pp := &PointsProduction{
        CurrentPoints: 0,
        Points: points,
    }
    if err := Database.Insert(pp); err != nil {
        panic(NewException("Could not create points production", err))
    }
    return pp
}

func (pp *PointsProduction) isCompleted() bool {
    return pp.Points == pp.CurrentPoints
}

func (pp *PointsProduction) getMissingPoints() uint8 {
    return pp.Points - pp.CurrentPoints
}

func (pp *PointsProduction) update() {
	if err := Database.Update(pp); err != nil {
		panic(NewException("Points production could not be udpated", err))
	}
}

func (pp *PointsProduction) delete() {
    if err := Database.Delete(pp); err != nil {
        panic(NewException("Could not delete points production", err))
    }
}

func (p *Planet) produceBuildingPoints() {
    buildingPoints := p.Settings.BuildingPoints
    if buildingPoints <= 0 || len(p.Buildings) == 0 {
        return
    }
    // Sort the buildings by construction date
    constructingBuildings := make(map[string]*Building, 0)
    var buildingDates []string
    for _, building := range p.Buildings {
        if building.Status != BuildingStatusConstructing {
            continue
        }
        date := building.CreatedAt.Format(time.RFC3339)
        // we use the date as a key for the constructing buildings map
        constructingBuildings[date] = &building
        buildingDates = append(buildingDates, date)
    }
    // Here we sort the dates
    sort.Strings(buildingDates)
    for _, date := range buildingDates {
        if buildingPoints <= 0 {
            break
        }
        buildingPoints = constructingBuildings[date].spendPoints(buildingPoints)
    }
}
