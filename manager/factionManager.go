package manager

import(
    "kalaxia-game-api/database"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/model"
)

func GetFaction(id uint16) *model.Faction {
  faction := &model.Faction{}
  if err := database.
    Connection.
    Model(faction).
    Column("faction.*", "Relations", "Relations.OtherFaction").
    Where("faction.id = ?", id).
    Select(); err != nil {
    panic(exception.NewHttpException(404, "Faction not found", nil))
  }
  return faction
}

func GetServerFactions(serverId uint16) []*model.Faction {
    factions := make([]*model.Faction, 0)
    if err := database.Connection.Model(&factions).Where("server_id = ?", serverId).Select(); err != nil {
        panic(exception.NewHttpException(404, "Server not found", nil))
    }
    return factions
}

func GetFactionPlanetChoices(factionId uint16) []*model.Planet {
    planets := make([]*model.Planet, 0)
    faction := &model.Faction{ Id: factionId }
    if err := database.Connection.Select(faction); err != nil {
        return planets
    }
    if _, err := database.
        Connection.
        Query(&planets, `SELECT p.*
        FROM map__maps m
        LEFT JOIN map__systems s ON s.map_id = m.id
        LEFT JOIN map__planets p ON p.system_id = s.id
        LEFT JOIN diplomacy__relations d ON d.planet_id = p.id
        WHERE p.player_id IS NULL AND m.server_id = ? AND d.faction_id = ?
        ORDER BY d.score DESC
        LIMIT 4`, faction.ServerId, faction.Id); err != nil {
            return planets
    }
    for _, planet := range planets {
        planet.Resources = make([]model.PlanetResource, 0)
        if err := database.
            Connection.
            Model(&planet.Resources).
            Where("planet_id = ?", planet.Id).
            Select(); err != nil {
                panic(exception.NewException("Planet resources not found", err))
        }
        planet.Relations = GetPlanetRelations(planet.Id)
    }
    return planets
}

func CreateServerFactions(server *model.Server, factions []interface{}) []*model.Faction {
    results := make([]*model.Faction, 0)
    for _, factionData := range factions {
        data := factionData.(map[string]interface{})
        faction := &model.Faction{
            Name: data["name"].(string),
            Description: data["description"].(string),
            Color: data["color"].(string),
            ServerId: server.Id,
            Server: server,
        }
        if err := database.Connection.Insert(faction); err != nil {
            panic(exception.NewException("Server faction could not be created", err))
        }
        results = append(results, faction)
    }
    createFactionsRelations(results)
    return results
}

func createFactionsRelations(factions []*model.Faction) {
    for _, faction := range factions {
        for _, otherFaction := range factions {
            if faction.Id == otherFaction.Id {
                continue
            }
            relation := &model.FactionRelation{
                FactionId: faction.Id,
                Faction: faction,
                OtherFactionId: otherFaction.Id,
                OtherFaction: otherFaction,
                State: model.RELATION_NEUTRAL,
            }
            if err := database.Connection.Insert(relation); err != nil {
                panic(exception.NewException("Faction relation could not be created", err))
            }
            faction.Relations = append(faction.Relations, relation)
        }
    }
}

func GetFactionMembers(factionId uint16) []*model.Player {
    members := make([]*model.Player, 0)
    if err := database.Connection.Model(&members).Where("faction_id = ?", factionId).Select(); err != nil {
        return members
    }
    return members
}
