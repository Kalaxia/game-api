package api

import(
  "encoding/json"
  "io/ioutil"
  "math"
  "math/rand"
)

const MinPlanetsPerSystem = 3

var planetsData PlanetsData
var resourcesData ResourcesData
var factions []*Faction
var planetsNameFrequencies []Element

func (m *Map) generateSystems(gameFactions []*Faction) {
    factions = gameFactions
    initializeConfiguration()
    generationProbability := 0
    for x := uint16(0); x < m.Size; x++ {
        for y := uint16(0); y < m.Size; y++ {
            random := rand.Intn(100)
            if random > generationProbability {
                generationProbability += 1
                continue
            }
            go m.generateSystem(x, y)
            generationProbability = 0
        }
        generationProbability = 0
    }
}

func initializeConfiguration() {
    planetsDataJSON, err := ioutil.ReadFile("/go/src/kalaxia-game-api/resources/planet_types.json")
    if err != nil {
		panic(NewException("Could not open planet types configuration file", err))
    }
    resourcesDataJSON, err := ioutil.ReadFile("/go/src/kalaxia-game-api/resources/resources.json")
    if err != nil {
		panic(NewException("Could not open resources configuration file", err))
    }
    if err := json.Unmarshal(planetsDataJSON, &planetsData); err != nil {
		panic(NewException("Could not read planet types configuration file", err))
    }
    if err := json.Unmarshal(resourcesDataJSON, &resourcesData); err != nil {
		panic(NewException("Could not read resources configuration file", err))
    }
    planetNamesJSON, err := ioutil.ReadFile("/go/src/kalaxia-game-api/resources/planet_names.json")
    if err != nil {
        panic(NewException("Could not read names file", err))
    }
    planetNames := make([]string, 0)
    if err := json.Unmarshal(planetNamesJSON, &planetNames); err != nil {
		panic(NewException("Could not read planet types configuration file", err))
    }
    planetsNameFrequencies = generateFrequencies(planetNames)
}

func (m *Map) generateSystem(x uint16, y uint16) {
    system := &System{
        Map: m,
        MapId: m.Id,
        X: x,
        Y: y,
    }
    if err := Database.Insert(system); err != nil {
		panic(NewException("System could not be created", err))
    }
    nbOrbits := rand.Intn(5) + MinPlanetsPerSystem
    for i := 1; i <= nbOrbits; i++ {
        go func(i int, system *System) {
            orbit := &SystemOrbit{
                Radius: uint16(i * 100 + rand.Intn(100)),
                System: system,
                SystemId: system.Id,
            }
            if err := Database.Insert(orbit); err != nil {
        		panic(NewException("Orbit could not be created", err))
            }
            system.Orbits = append(system.Orbits, *orbit)
            system.generatePlanet(orbit)
        } (i, system)
    }
}

func (s *System) generatePlanet(orbit *SystemOrbit) *Planet {
    planet := &Planet{
        Name: generatePlanetName(planetsNameFrequencies),
        Type: orbit.generatePlanetType(),
        System: s,
        SystemId: s.Id,
        Orbit: orbit,
        OrbitId: orbit.Id,
        Relations: make([]*DiplomaticRelation, 0),
    }
    planet.generatePlanetPopulation()
    planet.generateSettings()
    if err := Database.Insert(planet); err != nil {
		panic(NewException("Planet could not be created", err))
    }
    planet.createStorage()
    planet.Resources = planet.generatePlanetResources()
    s.Planets = append(s.Planets, *planet)
    if planet.Population > 0 {
        planet.generatePlanetRelations()
    }
    return planet
}

func (p *Planet) generatePlanetPopulation() {
    population := uint(rand.Intn(25))
    if population < 15 {
        p.Population = 0
        return
    }
    p.Population = population * populationPointRatio
    p.PublicOrder = planetPublicOrderStable
    p.TaxRate = planetTaxRateNormal
}

func (p *Planet) generateSettings() {
    points := p.calculatePopulationPoints()
    equalPart := uint8(math.Floor(float64(points) / 4.0))
    settings := &PlanetSettings{
        ServicesPoints: equalPart,
        BuildingPoints: equalPart,
        MilitaryPoints: equalPart,
        ResearchPoints: equalPart,
    }
    settings.ServicesPoints += uint8(points % 4)
    if err := Database.Insert(settings); err != nil {
		panic(NewException("Planet settings could not be created", err))
    }
    p.Settings = settings
    p.SettingsId = settings.Id
}

func (o *SystemOrbit) generatePlanetType() string {
    coeff := int(o.Radius) * rand.Intn(3) + rand.Intn(100)
    switch {
        case coeff < 200:
            return PlanetTypeVolcanic
        case coeff < 300:
            return PlanetTypeRocky
        case coeff < 400:
            return PlanetTypeDesert
        case coeff < 500:
            return PlanetTypeTropical
        case coeff < 600:
            return PlanetTypeTemperate
        case coeff < 700:
            return PlanetTypeOceanic
        default:
            return PlanetTypeArctic
    }
}

func (p *Planet) generatePlanetResources() []PlanetResource {
    resources := make([]PlanetResource, 0)
    for name, density := range planetsData[p.Type].Resources {
        go p.generatePlanetResource(&resources, name, density)
    }
    return resources
}

func (p *Planet) generatePlanetResource(resources *[]PlanetResource, name string, density uint8) {
    finalDensity := density + uint8(rand.Intn(30)) - uint8(rand.Intn(30))
    if finalDensity <= 0 { return }
    if finalDensity > 100 { finalDensity = 100 }
    planetResource := &PlanetResource{
        Name: name,
        Density: finalDensity,
        Planet: p,
        PlanetId: p.Id,
    }
    if err := Database.Insert(planetResource); err != nil {
		panic(NewException("Planet resource could not be created", err))
    }
    *resources = append(*resources, *planetResource)
}

func (p *Planet) generatePlanetRelations() {
    var bestRelation *DiplomaticRelation
    for _, faction := range factions {
        relation := p.createFactionRelation(faction, p.generatePlanetRelationScore())
        p.Relations = append(p.Relations, relation)
        if bestRelation == nil || relation.Score > bestRelation.Score {
            bestRelation = relation
        }
    }
    if bestRelation.Score > 400 && rand.Intn(100) > 50 {
        p.Faction = bestRelation.Faction
        p.FactionId = bestRelation.FactionId
        p.update()
    }
}

func (p *Planet) generatePlanetRelationScore() int {
    score := rand.Intn(500) - rand.Intn(500)

    if score > -50 && score < 50 {
        return 0
    }
    return score
}