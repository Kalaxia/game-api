package manager

import(
  "kalaxia-game-api/database"
  "kalaxia-game-api/model"
  "kalaxia-game-api/utils"
)

func GenerateMap(server *model.Server, factions []*model.Faction, size uint16) *model.Map {
    gameMap := &model.Map{
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

func GetMapByServerId(serverId uint16) *model.Map {
    gameMap := &model.Map{ServerId: serverId}
    if err := database.Connection.Model(gameMap).Where("server_id = ?", serverId).Select(); err != nil {
        return nil
    }
    gameMap.Systems = GetMapSystems(gameMap.Id)
    return gameMap
}
