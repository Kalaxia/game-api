package utils

import(
  "encoding/json"
  "io/ioutil"
  "math/rand"
  "kalaxia-game-api/database"
  "kalaxia-game-api/exception"
  "kalaxia-game-api/model"
)

const MIN_PLANETS_PER_SYSTEM = 3

var planetsData model.PlanetsData
var resourcesData model.ResourcesData
var factions []*model.Faction
var planetsNameFrequencies []Element

func GenerateMapSystems(gameMap *model.Map, gameFactions []*model.Faction) {
    factions = gameFactions
    initializeConfiguration()
    generationProbability := 0
    for x := uint16(0); x < gameMap.Size; x++ {
        for y := uint16(0); y < gameMap.Size; y++ {
            random := rand.Intn(100)
            if random > generationProbability {
                generationProbability += 1
                continue
            }
            go generateSystem(gameMap, x, y)
            generationProbability = 0
        }
        generationProbability = 0
    }
}

func initializeConfiguration() {
    planetsDataJSON, err := ioutil.ReadFile("/go/src/kalaxia-game-api/resources/planet_types.json")
    if err != nil {
		panic(exception.NewException("Could not open planet types configuration file", err))
    }
    resourcesDataJSON, err := ioutil.ReadFile("/go/src/kalaxia-game-api/resources/resources.json")
    if err != nil {
		panic(exception.NewException("Could not open resources configuration file", err))
    }
    if err := json.Unmarshal(planetsDataJSON, &planetsData); err != nil {
		panic(exception.NewException("Could not read planet types configuration file", err))
    }
    if err := json.Unmarshal(resourcesDataJSON, &resourcesData); err != nil {
		panic(exception.NewException("Could not read resources configuration file", err))
    }
    planetNamesJSON, err := ioutil.ReadFile("/go/src/kalaxia-game-api/resources/planet_names.json")
    if err != nil {
        panic(exception.NewException("Could not read names file", err))
    }
    planetNames := make([]string, 0)
    if err := json.Unmarshal(planetNamesJSON, &planetNames); err != nil {
		panic(exception.NewException("Could not read planet types configuration file", err))
    }
    planetsNameFrequencies = generateFrequencies(planetNames)
}

func generateSystem(gameMap *model.Map, x uint16, y uint16) {
    system := &model.System{
        Map: gameMap,
        MapId: gameMap.Id,
        X: x,
        Y: y,
    }
    if err := database.Connection.Insert(system); err != nil {
		panic(exception.NewException("System could not be created", err))
    }
    nbOrbits := rand.Intn(5) + MIN_PLANETS_PER_SYSTEM
    for i := 1; i <= nbOrbits; i++ {
        go func(i int, system *model.System) {
            orbit := &model.SystemOrbit{
                Radius: uint16(i * 100 + rand.Intn(100)),
                System: system,
                SystemId: system.Id,
            }
            if err := database.Connection.Insert(orbit); err != nil {
        		panic(exception.NewException("Orbit could not be created", err))
            }
            system.Orbits = append(system.Orbits, *orbit)
            generatePlanet(system, orbit)
        } (i, system)
    }
}

func generatePlanet(system *model.System, orbit *model.SystemOrbit) *model.Planet {
    planetType := choosePlanetType(orbit)
    settings := generateSettings()
    planet := &model.Planet{
        Name: generatePlanetName(planetsNameFrequencies),
        Type: planetType,
        System: system,
        SystemId: system.Id,
        Orbit: orbit,
        OrbitId: orbit.Id,
        Population: 2000000,
        Settings: settings,
        SettingsId: settings.Id,
    }
    if err := database.Connection.Insert(planet); err != nil {
		panic(exception.NewException("Planet could not be created", err))
    }
    planet.Resources = choosePlanetResources(planet, planetType)
    system.Planets = append(system.Planets, *planet)
    choosePlanetRelations(planet)
    return planet
}

func generateSettings() *model.PlanetSettings {
    settings := &model.PlanetSettings{
        ServicesPoints: 5,
        BuildingPoints: 5,
        MilitaryPoints: 5,
        ResearchPoints: 5,
    }
    if err := database.Connection.Insert(settings); err != nil {
		panic(exception.NewException("Planet settings could not be created", err))
    }
    return settings
}

func choosePlanetType(orbit *model.SystemOrbit) string {
    coeff := int(orbit.Radius) * rand.Intn(3) + rand.Intn(100)
    switch {
        case coeff < 300:
            return model.PlanetTypeVolcanic
        case coeff < 400:
            return model.PlanetTypeRocky
        case coeff < 500:
            return model.PlanetTypeDesert
        case coeff < 600:
            return model.PlanetTypeTropical
        case coeff < 700:
            return model.PlanetTypeTemperate
        case coeff < 800:
            return model.PlanetTypeOceanic
        default:
            return model.PlanetTypeArtic
    }
}

func choosePlanetResources(planet *model.Planet, planetType string) []model.PlanetResource {
    resources := make([]model.PlanetResource, 0)
    for name, density := range planetsData[planetType].Resources {
        go generatePlanetResource(&resources, name, density, planet)
    }
    return resources
}

func generatePlanetResource(resources *[]model.PlanetResource, name string, density uint8, planet *model.Planet) {
    finalDensity := density + uint8(rand.Intn(30)) - uint8(rand.Intn(30))
    if finalDensity <= 0 { return }
    if finalDensity > 100 { finalDensity = 100 }
    planetResource := &model.PlanetResource{
        Name: name,
        Density: finalDensity,
        Planet: planet,
        PlanetId: planet.Id,
    }
    if err := database.Connection.Insert(planetResource); err != nil {
		panic(exception.NewException("Planet resource could not be created", err))
    }
    *resources = append(*resources, *planetResource)
}

func choosePlanetRelations(planet *model.Planet) []model.DiplomaticRelation {
    relations := make([]model.DiplomaticRelation, 0)
    for _, faction := range factions {
        generatePlanetRelation(planet, faction, &relations)
    }
    return relations
}

func generatePlanetRelation(planet *model.Planet, faction *model.Faction, relations *[]model.DiplomaticRelation) {
    score := rand.Intn(500) - rand.Intn(500)

    if score > -50 && score < 50 {
        score = 0
    }
    relation := &model.DiplomaticRelation{
        Planet: planet,
        PlanetId: planet.Id,
        Faction: faction,
        FactionId: faction.Id,
        Score: score,
    }
    if err := database.Connection.Insert(relation); err != nil {
		panic(exception.NewException("Planet relation could not be created", err))
    }
    *relations = append(*relations, *relation)
}
