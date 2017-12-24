package manager

import(
  "kalaxia-game-api/database"
  "kalaxia-game-api/model/map"
)

func GetMapSystems(mapId uint16) []model.System {
  var systems []model.System
  if err := database.Connection.Model(&systems).Where("map_id = ?", mapId).Select(); err != nil {
    panic(err)
  }
  return systems
}

func GetSystem(id uint16) *model.System {
  system := model.System{Id: id}
  if err := database.Connection.Select(&system); err != nil {
    return nil
  }
  system.Planets = GetSystemPlanets(id)
  system.Orbits = GetSystemOrbits(id)
  return &system
}

func GetSystemOrbits(id uint16) []model.SystemOrbit {
  var systemOrbits []model.SystemOrbit
  if err := database.Connection.Model(&systemOrbits).Where("system_id = ?", id).Select(); err != nil {
    panic(err)
  }
  return systemOrbits
}
