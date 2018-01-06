package utils

import(
  "encoding/json"
  "io/ioutil"
  "math/rand"
  "kalaxia-game-api/database"
  mapModel "kalaxia-game-api/model/map"
  resourceModel "kalaxia-game-api/model/resource"
  factionModel "kalaxia-game-api/model/faction"
  diplomacyModel "kalaxia-game-api/model/diplomacy"
)

const MIN_PLANETS_PER_SYSTEM = 3

var planetsData mapModel.PlanetsData
var resourcesData resourceModel.ResourcesData
var factions []*factionModel.Faction

func GenerateMapSystems(gameMap *mapModel.Map, gameFactions []*factionModel.Faction) {
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
    panic(err)
  }
  resourcesDataJSON, err := ioutil.ReadFile("/go/src/kalaxia-game-api/resources/resources.json")
  if err != nil {
    panic(err)
  }
  if err := json.Unmarshal(planetsDataJSON, &planetsData); err != nil {
    panic(err)
  }
  if err := json.Unmarshal(resourcesDataJSON, &resourcesData); err != nil {
    panic(err)
  }
}

func generateSystem(gameMap *mapModel.Map, x uint16, y uint16) {
  system := &mapModel.System{
    Map: gameMap,
    MapId: gameMap.Id,
    X: x,
    Y: y,
  }
  if err := database.Connection.Insert(system); err != nil {
    panic(err)
  }
  nbOrbits := rand.Intn(5) + MIN_PLANETS_PER_SYSTEM
  for i := 1; i <= nbOrbits; i++ {
    go func(i int, system *mapModel.System) {
      orbit := &mapModel.SystemOrbit{
        Radius: uint16(i * 100 + rand.Intn(100)),
        System: system,
        SystemId: system.Id,
      }
      if err := database.Connection.Insert(orbit); err != nil {
        panic(err)
      }
      system.Orbits = append(system.Orbits, *orbit)
      generatePlanet(system, orbit)
    } (i, system)
  }
}

func generatePlanet(system *mapModel.System, orbit *mapModel.SystemOrbit) *mapModel.Planet {
  planetType := choosePlanetType(orbit)
  planet := &mapModel.Planet{
    Name: "Régalion V",
    Type: planetType,
    System: system,
    SystemId: system.Id,
    Orbit: orbit,
    OrbitId: orbit.Id,
  }
  if err := database.Connection.Insert(planet); err != nil {
    panic(err)
  }
  planet.Resources = choosePlanetResources(planet, planetType)
  system.Planets = append(system.Planets, *planet)
  choosePlanetRelations(planet)
  return planet
}

func choosePlanetType(orbit *mapModel.SystemOrbit) string {
  coeff := int(orbit.Radius) * rand.Intn(3) + rand.Intn(100)
  switch {
    case coeff < 300:
      return mapModel.PLANET_TYPE_VOLCANIC
    case coeff < 400:
      return mapModel.PLANET_TYPE_ROCKY
    case coeff < 500:
      return mapModel.PLANET_TYPE_DESERT
    case coeff < 600:
      return mapModel.PLANET_TYPE_TROPICAL
    case coeff < 700:
      return mapModel.PLANET_TYPE_TEMPERATE
    case coeff < 800:
      return mapModel.PLANET_TYPE_OCEANIC
    default:
      return mapModel.PLANET_TYPE_ARCTIC
  }
}

func choosePlanetResources(planet *mapModel.Planet, planetType string) []mapModel.PlanetResource {
  resources := make([]mapModel.PlanetResource, 0)
  for name, density := range planetsData[planetType].Resources {
    go generatePlanetResource(&resources, name, density, planet)
  }
  return resources
}

func generatePlanetResource(resources *[]mapModel.PlanetResource, name string, density uint8, planet *mapModel.Planet) {
  finalDensity := density + uint8(rand.Intn(30)) - uint8(rand.Intn(30))
  if finalDensity <= 0 { return }
  if finalDensity > 100 { finalDensity = 100 }
  planetResource := &mapModel.PlanetResource{
    Name: name,
    Density: finalDensity,
    Planet: planet,
    PlanetId: planet.Id,
  }
  if err := database.Connection.Insert(planetResource); err != nil {
    panic(err)
  }
  *resources = append(*resources, *planetResource)
}

func choosePlanetRelations(planet *mapModel.Planet) []diplomacyModel.DiplomaticRelation {
  relations := make([]diplomacyModel.DiplomaticRelation, 0)
  for _, faction := range factions {
    generatePlanetRelation(planet, faction, &relations)
  }
  return relations
}

func generatePlanetRelation(planet *mapModel.Planet, faction *factionModel.Faction, relations *[]diplomacyModel.DiplomaticRelation) {
  score := rand.Intn(500) - rand.Intn(500)

  if score > -50 && score < 50 {
    score = 0
  }
  relation := &diplomacyModel.DiplomaticRelation{
    Planet: planet,
    PlanetId: planet.Id,
    Faction: faction,
    FactionId: faction.Id,
    Score: score,
  }
  if err := database.Connection.Insert(relation); err != nil {
    panic(err)
  }
  *relations = append(*relations, *relation)
}
