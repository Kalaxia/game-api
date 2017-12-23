package utils

import(
  "math/rand"
  "kalaxia-game-api/database"
  "kalaxia-game-api/model/map"
)

const MIN_PLANETS_PER_SYSTEM = 3

func GenerateMapSystems(gameMap *model.Map) {
  generationProbability := 0
  for i := uint16(0); i < gameMap.Size; i++ {
    for j := uint16(0); j < gameMap.Size; j++ {
      random := rand.Intn(100)
      if random > generationProbability {
        generationProbability += 1
        continue
      }
      generateSystem(gameMap, i, j)
      generationProbability = 0
    }
    generationProbability = 0
  }
}

func generateSystem(gameMap *model.Map, x uint16, y uint16) {
  system := &model.System{
    Map: gameMap,
    MapId: gameMap.Id,
    X: x,
    Y: y,
  }
  if err := database.Connection.Insert(system); err != nil {
    panic(err)
  }
  nbOrbits := rand.Intn(5) + MIN_PLANETS_PER_SYSTEM
  for i := 0; i < nbOrbits; i++ {
    orbit := &model.SystemOrbit{
      Radius: uint16(i * 1000 + rand.Intn(500)),
      System: system,
      SystemId: system.Id,
    }
    if err := database.Connection.Insert(orbit); err != nil {
      panic(err)
    }
    system.Orbits = append(system.Orbits, orbit)
    generatePlanet(system, orbit)
  }
}

func generatePlanet(system *model.System, orbit *model.SystemOrbit) *model.Planet {
  planet := &model.Planet{
    Name: "Régalion V",
    Type: choosePlanetType(orbit),
    System: system,
    SystemId: system.Id,
    Orbit: orbit,
    OrbitId: orbit.Id,
  }
  if err := database.Connection.Insert(planet); err != nil {
    panic(err)
  }
  return planet
}

func choosePlanetType(orbit *model.SystemOrbit) string {
  return model.PLANET_TYPE_ROCKY
}
