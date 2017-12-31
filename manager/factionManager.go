package manager

import(
  "kalaxia-game-api/database"
  factionModel "kalaxia-game-api/model/faction"
  mapModel "kalaxia-game-api/model/map"
  serverModel "kalaxia-game-api/model/server"
)

func GetServerFactions(serverId uint16) []*factionModel.Faction {
  factions := make([]*factionModel.Faction, 0)
  if err := database.Connection.Model(&factions).Where("server_id = ?", serverId).Select(); err != nil {
    return nil
  }
  return factions
}

// @TODO this function is random for now.
// It will rely on diplomatic relations when it will be implemented
func GetFactionPlanetChoices(factionId uint16) []*mapModel.Planet {
  planets := make([]*mapModel.Planet, 0)
  faction := &factionModel.Faction{ Id: factionId }
  if err := database.Connection.Select(faction); err != nil {
    return planets
  }
  if _, err := database.
    Connection.
    Query(&planets, `SELECT p.*
      FROM map__maps m
      LEFT JOIN map__systems s ON s.map_id = m.id
      LEFT JOIN map__planets p ON p.system_id = s.id
      LEFT JOIN map__planet_resources pr ON pr.planet_id = p.id
      WHERE p.player_id IS NULL AND m.server_id = ?
      LIMIT 4`, faction.ServerId); err != nil {
    return planets
  }
  for _, planet := range planets {
      planet.Resources = make([]mapModel.PlanetResource, 0)
      if err := database.
        Connection.
        Model(&planet.Resources).
        Where("planet_id = ?", planet.Id).
        Select(); err != nil {
        panic(err)
      }
  }
  return planets
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
