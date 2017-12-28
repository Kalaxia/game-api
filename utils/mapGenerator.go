package utils

import(
  "math/rand"
  "kalaxia-game-api/database"
  "kalaxia-game-api/model/map"
)

const MIN_PLANETS_PER_SYSTEM = 3

func GenerateMapSystems(gameMap *model.Map) {
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
  for i := 1; i <= nbOrbits; i++ {
    go func(i int, system *model.System) {
      orbit := &model.SystemOrbit{
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
  system.Planets = append(system.Planets, *planet)
  return planet
}

func choosePlanetType(orbit *model.SystemOrbit) string {
  coeff := int(orbit.Radius) * rand.Intn(3) + rand.Intn(100)
  switch {
    case coeff < 300:
      return model.PLANET_TYPE_VOLCANIC
    case coeff < 400:
      return model.PLANET_TYPE_ROCKY
    case coeff < 500:
      return model.PLANET_TYPE_DESERT
    case coeff < 600:
      return model.PLANET_TYPE_TROPICAL
    case coeff < 700:
      return model.PLANET_TYPE_TEMPERATE
    case coeff < 800:
      return model.PLANET_TYPE_OCEANIC
    default:
      return model.PLANET_TYPE_ARCTIC
  }
}
