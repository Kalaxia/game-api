package manager

import(
  "kalaxia-game-api/database"
  serverModel "kalaxia-game-api/model/server"
  factionModel "kalaxia-game-api/model/faction"
  mapModel "kalaxia-game-api/model/map"
  "kalaxia-game-api/utils"
)

func GenerateMap(server *serverModel.Server, factions []*factionModel.Faction, size uint16) *mapModel.Map {
  gameMap := &mapModel.Map{
    Server: server,
    ServerId: server.Id,
    Size: size,
  }
  if err := database.Connection.Insert(gameMap); err != nil {
    panic(err)
  }
  utils.GenerateMapSystems(gameMap, factions)
  return gameMap
}

func GetMapByServerId(serverId uint16) *mapModel.Map {
  gameMap := &mapModel.Map{ServerId: serverId}
  if err := database.Connection.Model(gameMap).Where("server_id = ?", serverId).Select(); err != nil {
    return nil
  }
  gameMap.Systems = GetMapSystems(gameMap.Id)
  return gameMap
}
