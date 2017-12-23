package manager

import(
  "kalaxia-game-api/database"
  "kalaxia-game-api/model/map"
)

func GetMapSystems(mapId uint16) []model.System {
  var systems []model.System
  if err := database.Connection.Model(&systems).Select(); err != nil {
    panic(err)
  }
  return systems
}
