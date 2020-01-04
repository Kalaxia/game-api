package api

import(
	"github.com/gorilla/mux"
	"github.com/go-pg/pg/v9/orm"
	"net/http"
	"strconv"
	"time"
)

type(
	FactionWar struct {
		tableName struct{} `pg:"faction__wars"`

		Id uint32 `json:"id"`
		FactionId uint16 `json:"-"`
		Faction *Faction `json:"faction"`
		TargetId uint16 `json:"-"`
		Target *Faction `json:"target"`
		CasusBelli []*FactionCasusBelli `json:"casus_belli" pg:"fk:war_id"`
		CreatedAt time.Time `json:"created_at"`
		EndedAt *time.Time `json:"ended_at"`
	}

	FactionCasusBelli struct {
		tableName struct{} `pg:"faction__casus_belli"`

		Id uint64 `json:"id"`
		FactionId uint16 `json:"-"`
		Faction *Faction `json:"faction"`
		VictimId uint16 `json:"-"`
		Victim *Faction `json:"victim"`
		PlayerId uint16 `json:"-"`
		Player *Player `json:"player"`
		Type string `json:"type"`
		Data map[string]interface{} `json:"data"`
		WarId uint32 `json:"-"`
		War *FactionWar `json:"war"`
		CreatedAt time.Time `json:"created_at"`
	}
)

const(
	CasusBelliTypePlunder = "plunder"
	CasusBelliTypeConquest = "conquest"
	CasusBelliTypeCombat = "combat"
)

func GetFactionUnansweredCasusBelli(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseUint(mux.Vars(r)["faction_id"], 10, 16)
	faction := getFaction(uint16(id))
	targetRawId := r.URL.Query()["faction_id"][0]
	targetId, _ := strconv.ParseUint(targetRawId, 10, 16)
	target := getFaction(uint16(targetId))

    SendJsonResponse(w, 200, faction.getUnansweredCasusBelli(target))
}

func GetFactionCasusBelli(w http.ResponseWriter, r *http.Request) {
	factionId, _ := strconv.ParseUint(mux.Vars(r)["faction_id"], 10, 16)
	faction := getFaction(uint16(factionId))
	casusBelliId, _ := strconv.ParseUint(mux.Vars(r)["casus_belli_id"], 10, 16)

    SendJsonResponse(w, 200, faction.getCasusBelli(casusBelliId))
}

func GetFactionWars(w http.ResponseWriter, r *http.Request) {
	factionId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	faction := getFaction(uint16(factionId))

    SendJsonResponse(w, 200, faction.getWars())
}

func GetFactionWar(w http.ResponseWriter, r *http.Request) {
	factionId, _ := strconv.ParseUint(mux.Vars(r)["faction_id"], 10, 16)
	faction := getFaction(uint16(factionId))
	warId, _ := strconv.ParseUint(mux.Vars(r)["war_id"], 10, 16)

    SendJsonResponse(w, 200, faction.getWar(uint32(warId)))
}

func (f *Faction) createWar(target *Faction) *FactionWar {
	war := &FactionWar{
		FactionId: f.Id,
		Faction: f,
		TargetId: target.Id,
		Target: target,
		CreatedAt: time.Now(),
	}
	if err := Database.Insert(war); err != nil {
		panic(NewException("Could not create faction war", err))
	}
	for _, cb := range f.getUnansweredCasusBelli(target) {
		cb.WarId = war.Id
		cb.War = war
		cb.update()
	}
	f.notify(NotificationTypeFaction, "faction.war.war_declared", map[string]interface{}{
		"enemy_name": target.Name,
	})
	target.notify(NotificationTypeFaction, "faction.war.war_alert", map[string]interface{}{
		"enemy_name": f.Name,
	})
	return war
}

func (f *Faction) endWar(target *Faction) {
	war := f.getCurrentWarWith(target)
	e := time.Now()
	war.EndedAt = &e
	war.update()

	f.notify(NotificationTypeFaction, "faction.war.war_ended", map[string]interface{}{
		"enemy_name": target.Name,
	})
	target.notify(NotificationTypeFaction, "faction.war.war_ended", map[string]interface{}{
		"enemy_name": f.Name,
	})
}

func (f *Faction) createCasusBelli(victim *Faction, p *Player, cType string, data map[string]interface{}) *FactionCasusBelli {
	casusBelli := &FactionCasusBelli{
		FactionId: f.Id,
		Faction: f,
		VictimId: victim.Id,
		Victim: victim,
		PlayerId: p.Id,
		Player: p,
		Type: cType,
		Data: data,
		CreatedAt: time.Now(),
	}
	if err := Database.Insert(casusBelli); err != nil {
		panic(NewException("Could not create casus belli", err))
	}
	return casusBelli
}

func (f *Faction) getCasusBelli(id uint64) *FactionCasusBelli {
	fcb := &FactionCasusBelli{}
	if err := Database.Model(fcb).Relation("Victim").Relation("Player").Where("faction_casus_belli.faction_id = ?", f.Id).Where("faction_casus_belli.id = ?", id).Select(); err != nil {
		panic(NewHttpException(404, "Faction Casus Belli not found", err))
	}
	return fcb
}

func (f *Faction) getUnansweredCasusBelli(faction *Faction) []*FactionCasusBelli {
	fcb := make([]*FactionCasusBelli, 0)
	if err := Database.Model(&fcb).Relation("Victim").Relation("Player").Where("victim_id = ?", f.Id).Where("faction_casus_belli.faction_id = ?", faction.Id).Where("war_id IS NULL").Select(); err != nil {
		panic(NewException("Could not retrieve faction casus belli", err))
	}
	return fcb
}

func (f *Faction) getWars() []*FactionWar {
	wars := make([]*FactionWar, 0)
	if err := Database.Model(&wars).Relation("Faction").Relation("Target").Where("faction_id = ?", f.Id).WhereOr("target_id = ?", f.Id).Select(); err != nil {
		panic(NewException("Could not retrieve faction wars", err))
	}
	return wars
}

func (f *Faction) getWar(id uint32) *FactionWar {
	war := &FactionWar{}
	if err := Database.Model(war).Relation("Target").Where("faction_war.id = ?", id).Select(); err != nil {
		panic(NewException("Could not retrieve faction war", err))
	}
	return war
}

func (f *Faction) getCurrentWarWith(target *Faction) *FactionWar {
	war := &FactionWar{}
	if err := Database.Model(war).
		WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			return q.Where("faction_id = ? AND target_id = ?", f.Id, target.Id).WhereOr("faction_id = ? AND target_id = ?", target.Id, f.Id), nil
		}).Where("ended_at IS NULL").Select(); err != nil {
		panic(NewException("Could not retrieve faction war", err))
	}
	return war
}

func (fcb *FactionCasusBelli) update() {
	if err := Database.Update(fcb); err != nil {
		panic(NewException("Could not update faction casus belli", err))
	}
}

func (fw *FactionWar) update() {
	if err := Database.Update(fw); err != nil {
		panic(NewException("Could not update faction war", err))
	}
}