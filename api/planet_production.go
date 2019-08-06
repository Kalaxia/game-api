package api

import(
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
            go planet.calculateProduction(&wg)
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

func (p *Planet) calculateProduction(wg *sync.WaitGroup) {
    defer wg.Done()
    defer CatchException(nil)

    points := p.getAvailablePoints()
    for _, b := range p.Buildings {
        b.Planet = p
        points = b.produce(points)
    }
    p.produceBuildingPoints()
    p.Storage.update()
}

func (p *Planet) getAvailablePoints() map[string]uint8 {
    return map[string]uint8{
        "services": p.Settings.ServicesPoints,
        "building": p.Settings.BuildingPoints,
        "military": p.Settings.MilitaryPoints,
        "research": p.Settings.ResearchPoints,
    }
}

func (b *Building) produce(points map[string]uint8) map[string]uint8 {
    switch (b.Type) {
        case "resource": 
            b.produceResources()
            return points
        case "shipyard":
            return b.produceShips(points)
        default:
            return points
    }
}

func (b *Building) produceResources() {
    for _, resourceName := range buildingPlansData[b.Name].Resources {
        if quantity := int16(b.getProducedQuantity(resourceName)); quantity > 0 {
            b.Planet.Storage.storeResource(resourceName, quantity)
        }
    }
}

func (b *Building) getProducedQuantity(resourceName string) uint16 {
    resource := b.Planet.getResource(resourceName)
    if resource == nil {
        return 0
    }
    return uint16(resource.Density) * 10
}

func (b *Building) produceShips(points map[string]uint8) map[string]uint8 {
    constructionGroups := b.Planet.getConstructingShips()
    if (len(constructionGroups) == 0) {
        return points
    }
    for _, group := range constructionGroups {
        if points["military"] < 1 {
            break
        }
        neededPoints := group.ConstructionState.Points - group.ConstructionState.CurrentPoints
        if neededPoints <= points["military"] {
            points["military"] -= neededPoints
            group.finishConstruction()
        } else {
            group.ConstructionState.CurrentPoints += points["military"]
            group.ConstructionState.update()
            points["military"] = 0
            break
        }
    }
    return points
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

func (pp *PointsProduction) spendPoints(points uint8, callback func()) uint8 {
    if missingPoints := pp.getMissingPoints(); missingPoints <= points {
        callback()
        return points - missingPoints
    }
    pp.CurrentPoints += points
    pp.update()
    return 0
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
    points := p.Settings.BuildingPoints
    if points <= 0 || len(p.Buildings) == 0 {
        return
    }
    for _, b := range p.Buildings {
        if points <= 0 {
            break
        }
        if b.Status == BuildingStatusConstructing {
            points = b.ConstructionState.spendPoints(points, b.finishConstruction)
        } else {
            points = b.constructCompartments(points)
        }
    }
}

func (b *Building) constructCompartments(points uint8) uint8 {
    for _, c := range b.Compartments {
        if points < 1 {
            break;
        }
        if c.Status != BuildingStatusConstructing {
            continue
        }
        points = c.ConstructionState.spendPoints(points, c.finishConstruction)
    }
    return points
}