package api

import(
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"net/http"
	"strconv"
)

type Faction struct {
	TableName struct{} `json:"-" sql:"faction__factions"`

	Id uint16 `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Description string `json:"description"`
	Colors map[string]string `json:"colors"`
	Banner string `json:"banner"`
	ServerId uint16 `json:"-"`
	Server *Server `json:"-"`
	Relations []*FactionRelation `json:"relations"`
}

func GetFactions(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)

    SendJsonResponse(w, 200, getServerFactions(player.ServerId))
}

func GetFaction(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)

    SendJsonResponse(w, 200, getFaction(uint16(id)))
}

func GetFactionPlanetChoices(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)

    SendJsonResponse(w, 200, getFactionPlanetChoices(uint16(id)))
}

func GetFactionMembers(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)

    SendJsonResponse(w, 200, getFactionMembers(uint16(id)))
}

func getFaction(id uint16) *Faction {
	faction := &Faction{}
	if err := Database.
		Model(faction).
		Column("faction.*", "Relations", "Relations.OtherFaction").
		Where("faction.id = ?", id).
		Select(); err != nil {
		panic(NewHttpException(404, "Faction not found", err))
	}
	return faction
}
  
func getServerFactions(serverId uint16) []*Faction {
	factions := make([]*Faction, 0)
	if err := Database.Model(&factions).Where("server_id = ?", serverId).Select(); err != nil {
		panic(NewHttpException(404, "Server not found", err))
	}
	return factions
}

func getFactionPlanetChoices(factionId uint16) []*Planet {
	planets := make([]*Planet, 0)
	if _, err := Database.
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
		planet.Resources = make([]PlanetResource, 0)
		if err := Database.
			Model(&planet.Resources).
			Where("planet_id = ?", planet.Id).
			Select(); err != nil {
				panic(NewException("Planet resources not found", err))
		}
		planet.Relations = planet.getPlanetRelations()
	}
	return planets
}

func (s *Server) createFactions(factions []interface{}) []*Faction {
	results := make([]*Faction, 0)
	for _, factionData := range factions {
		data := factionData.(map[string]interface{})
		colors := make(map[string]string, 0)

		for k, v := range data["colors"].(map[string]interface{}) {
			colors[k] = v.(string)
		}
		faction := &Faction{
			Name: data["name"].(string),
			Slug: slug.Make(data["name"].(string)),
			Description: data["description"].(string),
			Colors: colors,
			Banner: data["banner"].(string),
			ServerId: s.Id,
			Server: s,
		}
		if err := Database.Insert(faction); err != nil {
			panic(NewException("Server faction could not be created", err))
		}
		results = append(results, faction)
	}
	createFactionsRelations(results)
	return results
}

func createFactionsRelations(factions []*Faction) {
	for _, faction := range factions {
		for _, otherFaction := range factions {
			if faction.Id == otherFaction.Id {
				continue
			}
			relation := &FactionRelation{
				FactionId: faction.Id,
				Faction: faction,
				OtherFactionId: otherFaction.Id,
				OtherFaction: otherFaction,
				State: RelationNeutral,
			}
			if err := Database.Insert(relation); err != nil {
				panic(NewException("Faction relation could not be created", err))
			}
			faction.Relations = append(faction.Relations, relation)
		}
	}
}

func getFactionMembers(factionId uint16) []*Player {
	members := make([]*Player, 0)
	if err := Database.Model(&members).Where("faction_id = ?", factionId).Select(); err != nil {
		return members
	}
	return members
}
  