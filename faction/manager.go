package faction

import(
    "kalaxia-game-api/database"
    "kalaxia-game-api/server"
)

func GetFaction(id uint16) *Faction {
    faction := &Faction{Id: id}
    if err := database.Connection.Select(faction); err != nil {
        return nil
    }
    return faction
}

func GetServerFactions(serverId uint16) []*Faction {
    factions := make([]*Faction, 0)
    if err := database.Connection.Model(&factions).Where("server_id = ?", serverId).Select(); err != nil {
        return nil
    }
    return factions
}

func GetFactionPlanetChoices(factionId uint16) []*interface{} {
    planets := make([]*interface{}, 0)
    faction := &Faction{ Id: factionId }
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
        planet.Resources = make([]galaxy.PlanetResource, 0)
        if err := database.
            Connection.
            Model(&planet.Resources).
            Where("planet_id = ?", planet.Id).
            Select(); err != nil {
                panic(err)
        }
        relations := GetPlanetRelations(planet.Id)
        r := make([]interface{}, len(relations))
        for i, v := range relations {
            r[i] = v
        }
        planet.Relations = r
    }
    return planets
}

func CreateServerFactions(server *server.Server, factions []interface{}) []*Faction {
    results := make([]*Faction, 0)
    for _, factionData := range factions {
        data := factionData.(map[string]interface{})
        faction := &Faction{
            Name: data["name"].(string),
            Description: data["description"].(string),
            Color: data["color"].(string),
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


func GetPlanetRelations(planetId uint16) []DiplomaticRelation {
    relations := make([]DiplomaticRelation, 0)
    if err := database.
        Connection.
        Model(&relations).
        Column("diplomatic_relation.*", "Faction", "Player", "Player.Faction").
        Where("diplomatic_relation.planet_id = ?", planetId).
        Select(); err != nil {
            return nil
    }
    return relations
}

func IncreasePlayerRelation(planetId uint16, playerId uint16, score int) {
    var relation DiplomaticRelation
    if err := database.
        Connection.
        Model(&relation).
        Where("planet_id = ?", planetId).
        Where("player_id = ?", playerId).
        Select(); err != nil {
            createPlayerRelation(planetId, playerId, score)
            return
    }
    relation.Score += score
    if _, err := database.
        Connection.
        Model(&relation).
        Set("score = ?score").
        Where("planet_id = ?", planetId).
        Where("player_id = ?", playerId).
        Update(); err != nil {
            panic(err)
    }
}

func createPlayerRelation(planetId uint16, player uint16, score int) {
    relation := &DiplomaticRelation{
        PlanetId: planetId,
        PlayerId: playerId,
        Score: score,
    }
    if err := database.Connection.Insert(relation); err != nil {
        panic(err)
    }
}
