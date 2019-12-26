package api

import(
	"encoding/json"
	"net/http"
	"github.com/gorilla/context"
	"time"
	"strconv"
)

const(
	TerritoryActionCreation="creation"
	TerritoryActionPlanetGained="planet_gained"
	TerritoryActionPlanetLost="planet_lost"
	TerritoryActionSystemGained="system_gained"
	TerritoryActionSystemLost="system_lost"
	TerritoryActionConquest="conquest"
	TerritoryStatusPledge="pledge"
	TerritoryStatusContest="contest"
)

type(
	Territory struct {
		tableName struct{} `pg:"map__territories"`

		Id uint16 `json:"id"`
		MapId uint16 `json:"-"`
		Map *Map `json:"-"`
		PlanetId uint16 `json:"-"`
		Planet *Planet `json:"planet"`
		History []*TerritoryHistory `json:"history"`
	}

	TerritoryHistory struct {
		tableName struct{} `pg:"map__territory_histories"`

		Id uint16 `json:"id"`
		TerritoryId uint16 `json:"-"`
		Territory *Territory `json:"-"`
		PlayerId uint16 `json:"-"`
		Player *Player `json:"player"`
		Action string `json:"action"`
		Data map[string]interface{} `json:"data" pg:",use_zero"`
		HappenedAt time.Time `json:"happened_at"`
	}

	TerritoryCacheItem struct {
		Hash string `json:"-"`
		Id uint16 `json:"id"`
		Planet *Planet `json:"planet"`
		Systems []*SystemTerritoryCacheItem `json:"systems"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

func GetMapTerritories(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)
	starmap := getServerMap(player.ServerId)
	
	SendJsonResponse(w, 200, starmap.getTerritoriesCache())
}

func InitTerritoriesCache() {
	for _, m := range getAllMaps() {
		RedisClient.Del("map_" + strconv.FormatUint(uint64(m.Id), 10) + "_territories")
		for _, t := range m.getTerritories() {
			t.generateCache().store(m.Id)
		}
	}
}

func (m *Map) getTerritories() []*Territory {
	territories := make([]*Territory, 0)
	if err := Database.Model(&territories).Relation("Planet.Player.Faction").Relation("Planet.System").Where("territory.map_id = ?", m.Id).Select(); err != nil {
		panic(NewException("Could not retrieve map territories", err))
	}
	return territories
}

func (m *Map) getTerritoriesCache() map[string]string {
	result, err := RedisClient.HGetAll("map_" + strconv.FormatUint(uint64(m.Id), 10) +  "_territories").Result()
	if err != nil {
		panic(NewException("Could not retrieve territories cache", err))
	}
	return result
}

func (p *Planet) createTerritory() *Territory {
	t := &Territory{
		MapId: p.System.MapId,
		Map: p.System.Map,
		PlanetId: p.Id,
		Planet: p,
	}
	if err := Database.Insert(t); err != nil {
		panic(NewException("Could not create territory", err))
	}
	t.addHistory(p.Player, TerritoryActionCreation, make(map[string]interface{}, 0))
	p.System.addTerritory(t)
	t.expand()
	return t
}

func (t *Territory) addHistory(p *Player, action string, data map[string]interface{}) *TerritoryHistory {
	th := &TerritoryHistory{
		TerritoryId: t.Id,
		Territory: t,
		PlayerId: p.Id,
		Player: p,
		Action: action,
		Data: data,
		HappenedAt: time.Now(),
	}
	if err := Database.Insert(th); err != nil {
		panic(NewException("Could not create territory history", err))
	}
	t.History = append(t.History, th)
	return th
}

func (t *Territory) expand() {
	hasNew := true
	for hasNew == true {
		for _, st := range t.getSystemTerritories() {
			hasNew = hasNew && st.checkForIncludedSystems()
		}
	}
	t.update()
	t.generateCache().store(t.MapId)
}

func (t *Territory) update() {
	if err := Database.Update(t); err != nil {
		panic(NewException("Could not update territory", err))
	}
}

func (tci *TerritoryCacheItem) store(mapId uint16) {
	data, err := json.Marshal(tci)
	if err != nil {
		panic(NewException("Could not store territory cache", err))
	}
	RedisClient.HSet("map_" + strconv.FormatUint(uint64(mapId), 10)+ "_territories", tci.Hash, data)
}

func (t *Territory) generateCache() *TerritoryCacheItem {
	systemTerritories := t.getSystemTerritories() 
	tci := &TerritoryCacheItem{
		Hash: "territory-" + strconv.FormatUint(uint64(t.Id), 10),
		Id: t.Id,
		Planet: t.Planet,
		Systems: make([]*SystemTerritoryCacheItem, 0),
		UpdatedAt: time.Now(),
	}
	for _, st := range systemTerritories {
		if st.Status == TerritoryStatusPledge {
			tci.Systems = append(tci.Systems, st.generateCache())
		}
	}
	return tci
}