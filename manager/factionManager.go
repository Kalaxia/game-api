package manager

import(
    "kalaxia-game-api/database"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/model"

	"github.com/gosimple/slug"
)

func GetFaction(id uint16) *model.Faction {
  faction := &model.Faction{}
  if err := database.
    Connection.
    Model(faction).
    Column("faction.*", "Relations", "Relations.OtherFaction").
    Where("faction.id = ?", id).
    Select(); err != nil {
    panic(exception.NewHttpException(404, "Faction not found", err))
  }
  return faction
}

func GetServerFactions(serverId uint16) []*model.Faction {
    factions := make([]*model.Faction, 0)
    if err := database.Connection.Model(&factions).Where("server_id = ?", serverId).Select(); err != nil {
        panic(exception.NewHttpException(404, "Server not found", err))
    }
    return factions
}

func GetFactionPlanetChoices(factionId uint16) []*model.Planet {
    planets := make([]*model.Planet, 0)
    if _, err := database.
        Connection.
        Query(&planets, `SELECT p.*, s.id AS system__id, s.x AS system__x, s.y AS system__y
        FROM map__planets p
        INNER JOIN map__systems s ON p.system_id = s.id
        LEFT JOIN diplomacy__relations d ON d.planet_id = p.id
        WHERE p.player_id IS NULL AND d.faction_id = ?
        ORDER BY d.score DESC
        LIMIT 3`, factionId); err != nil {
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
        colors := make(map[string]string, 0)

        for k, v := range data["colors"].(map[string]interface{}) {
            colors[k] = v.(string)
        }
        faction := &model.Faction{
            Name: data["name"].(string),
            Slug: slug.Make(data["name"].(string)),
            Description: data["description"].(string),
            Colors: colors,
            Banner: data["banner"].(string),
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
