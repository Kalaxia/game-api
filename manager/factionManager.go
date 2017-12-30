package manager

import(
  "kalaxia-game-api/database"
  factionModel "kalaxia-game-api/model/faction"
  serverModel "kalaxia-game-api/model/server"
)

func GetServerFactions(serverId uint16) []*factionModel.Faction {
  factions := make([]*factionModel.Faction, 0)
  if err := database.Connection.Model(&factions).Where("server_id = ?", serverId).Select(); err != nil {
    return nil
  }
  return factions
}


func CreateServerFactions(server *serverModel.Server, factions []interface{}) []*factionModel.Faction {
  results := make([]*factionModel.Faction, 0)
  for _, factionData := range factions {
    data := factionData.(map[string]interface{})
    faction := &factionModel.Faction{
      Name: data["name"].(string),
      Description: data["description"].(string),
      ServerId: server.Id,
      Server: server,
    }
    if err := database.Connection.Insert(faction); err != nil {
      panic(err)
    }
    results = append(results, faction)
  }
  return results
}
