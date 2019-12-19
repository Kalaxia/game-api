package api

import(
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"net/http"
	"strconv"
)

type Faction struct {
	tableName struct{} `pg:"faction__factions"`

	Id uint16 `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Description string `json:"description"`
	Colors map[string]string `json:"colors"`
	Banner string `json:"banner"`
	ServerId uint16 `json:"-"`
	Server *Server `json:"-"`
	Relations []*FactionRelation `json:"relations"`
	Settings []*FactionSettings `json:"settings"`
	Wallet int32 `json:"wallet" pg:",notnull,use_zero"`
}

func CalculateFactionsWages() {
	for _, faction := range getAllFactions() {
		go faction.calculateWage()
	}
}

func GetFactions(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)

    SendJsonResponse(w, 200, getServerFactions(player.ServerId))
}

func GetFaction(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	player := context.Get(r, "player").(*Player)
	faction := getFaction(uint16(id))

	faction.Settings = faction.getAllSettings(player.Faction.Id != faction.Id)
	
    SendJsonResponse(w, 200, faction)
}

func GetFactionPlanetChoices(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)

    SendJsonResponse(w, 200, getFactionPlanetChoices(uint16(id)))
}

func GetFactionMembers(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	player := context.Get(r, "player").(*Player)
	faction := getFaction(uint16(id))

	if player.Faction.Id != faction.Id {
		panic(NewHttpException(403, "forbidden", nil))
	}
    SendJsonResponse(w, 200, faction.getMembers())
}

func getFaction(id uint16) *Faction {
	faction := &Faction{}
	if err := Database.
		Model(faction).
		Relation("Relations.OtherFaction").
		Where("faction.id = ?", id).
		Select(); err != nil {
		panic(NewHttpException(404, "Faction not found", err))
	}
	return faction
}

func getAllFactions() []*Faction {
	factions := make([]*Faction, 0)
	if err := Database.Model(&factions).Relation("Relations").Relation("Settings").Select(); err != nil {
		panic(NewException("Could not retrieve factions", err))
	}
	return factions
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
			Wallet: 10000,
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

func (f *Faction) countMembers() int {
	count, err := Database.Model(&Player{}).Where("faction_id = ?", f.Id).Count()
	if err != nil {
		panic(NewException("Could not count faction members", err))
	}
	return count
}

func (f *Faction) getMembers() []*Player {
	members := make([]*Player, 0)
	if err := Database.Model(&members).Where("faction_id = ?", f.Id).Select(); err != nil {
		return members
	}
	return members
}

func (f *Faction) calculateWage() {
	defer CatchException(nil)

	planetTaxes := int32(f.getSettings(FactionSettingsPlanetTaxes).Value)

	for _, planet := range f.getControlledPlanets() {
		if planet.Player.updateWallet(-planetTaxes) {
			f.Wallet += planetTaxes
			planet.Player.update()

		}
	}
	f.update()
}

func (f *Faction) getControlledPlanets() []*Planet {
	planets := make([]*Planet, 0)
	if err := Database.Model(&planets).Relation("Player.Faction").Where("player__faction.id = ?", f.Id).Select(); err != nil {
		panic(NewException("Could not retrieve faction controlled planets", err))
	}
	return planets
}

func (f *Faction) update() {
	if err := Database.Update(f); err != nil {
		panic(NewException("could not update faction", err))
	}
}

func (f *Faction) notify(nType string, content string, data map[string]interface{}) {
	for _, member := range f.getMembers() {
		member.notify(nType, content, data)
	}
}

func (f *Faction) validatePlanetTaxesMotion(taxes int) {
	settings := f.getSettings(FactionSettingsPlanetTaxes)

	if settings.Value == taxes {
		panic(NewHttpException(400, "faction.motions.types.planet_taxes.same_value", nil))
	}
	if taxes < 0 {
		panic(NewHttpException(400, "faction.motions.types.planet_taxes.invalid_value", nil))
	}
}

func (f *Faction) updatePlanetTaxes(taxes int) {
	settings := f.getSettings(FactionSettingsPlanetTaxes)
	settings.Value = taxes
	if err := Database.Update(settings); err != nil {
		panic(NewException("Could not update planet taxes", err))
	}
}