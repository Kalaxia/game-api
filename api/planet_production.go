package api

import(
    "math"
    "sync"
)

const PlanetProductionResourceDensityCoeff = 10

type(
    ResourceProduction struct {
        Name string `json:"name"`
        Density uint8 `json:"density"`
        BaseQuantity uint16 `json:"base_quantity"`
        FinalQuantity uint16 `json:"final_quantity"`
        Percent int8 `json:"percent"`
    }

    PointsProduction struct {
		tableName struct{} `json:"-" pg:"map__planet_point_productions"`

		Id uint32 `json:"id"`
		CurrentPoints uint8 `json:"current_points" pg:",notnull,use_zero"`
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
        Relation("Player").
        Relation("Buildings.ConstructionState").
        Relation("Buildings.Compartments.ConstructionState").
        Relation("Resources").
        Relation("Storage").
        Relation("Settings").
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
    for _, rp := range p.getProducedResources() {
        p.Storage.storeResource(rp.Name, int16(rp.FinalQuantity))
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
    if b.Status != BuildingStatusOperational {
        return points
    }
    switch (b.Type) {
        case BuildingTypeShipyard:
            return b.produceShips(points)
        default:
            return points
    }
}

func (b *Building) getResourceModifiers() map[string]int8 {
    resourceModifiers := make(map[string]int8, 0)

    for _, c := range b.Compartments {
        if c.Status != BuildingStatusOperational {
            continue
        }
        plan := b.getCompartmentPlan(c.Name)
        for _, m := range plan.Modifiers {
            if m.Type != "resource" {
                continue
            }
            resourceModifiers[m.Resource] = m.Percent
        }
    }
    return resourceModifiers
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

func (p *Planet) getProducedResources() map[string]*ResourceProduction {
    resourcesProduction := make(map[string]*ResourceProduction, len(p.Resources))

    for _, r := range p.Resources {
        resourcesProduction[r.Name] = &ResourceProduction{
            Name: r.Name,
            Density: r.Density,
            BaseQuantity: 0,
            FinalQuantity: 0,
            Percent: 0,
        }
    }
    for _, b := range p.Buildings {
        if b.Type != BuildingTypeResource {
            continue
        }
        resourcesProduction = b.getProducedResources(resourcesProduction)
    }
    for _, rp := range resourcesProduction {
        rp.calculateFinalQuantity()
    }
    return resourcesProduction
}

func (b *Building) getProducedResources(resourcesProduction map[string]*ResourceProduction) map[string]*ResourceProduction {
    modifiers := b.getResourceModifiers()
    for _, resourceName := range buildingPlansData[b.Name].Resources {
        if rp, ok := resourcesProduction[resourceName]; ok {
            rp.BaseQuantity += uint16(rp.Density) * PlanetProductionResourceDensityCoeff
            if percent, ok := modifiers[resourceName]; ok {
                rp.Percent += percent
            }
        }
    }
    return resourcesProduction
}

func (rp *ResourceProduction) calculateFinalQuantity() {
    rp.FinalQuantity = rp.BaseQuantity + uint16(math.Floor(float64(rp.BaseQuantity) * (float64(rp.Percent) / 100)))
}