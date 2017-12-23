package manager

import(
  "kalaxia-game-api/database"
  serverModel "kalaxia-game-api/model/server"
  mapModel "kalaxia-game-api/model/map"
  "kalaxia-game-api/utils"
)

func GenerateMap(server *serverModel.Server, size uint16) mapModel.Map {
  gameMap := mapModel.Map{
    Server: server,
    ServerId: server.Id,
    Size: size,
  }
  if err := database.Connection.Insert(&gameMap); err != nil {
    panic(err)
  }
  utils.GenerateMapSystems(&gameMap)
  return gameMap
}
