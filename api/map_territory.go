package api

import(
	"net/http"
	"github.com/gorilla/context"
	"math"
	"time"
)

const(
	TerritoryActionCreation="creation"
	TerritoryActionConquest="conquest"
	TerritoryStatusPledge="pledge"
	TerritoryStatusContest="contest"
)

type(
	Territory struct {
		TableName struct{} `json:"-" sql:"map__territories"`

		Id uint16 `json:"id"`
		MapId uint16 `json:"-"`
		Map *Map `json:"-"`
		PlanetId uint16 `json:"-"`
		Planet *Planet `json:"planet"`
		MilitaryInfluence uint16 `json:"military_influence"`
		EconomicInfluence uint16 `json:"economic_influence"`
		CulturalInfluence uint16 `json:"cultural_influence"`
		PoliticalInfluence uint16 `json:"political_influence"`
		ReligiousInfluence uint16 `json:"religious_influence"`
		Coordinates []*Coordinates `json:"coordinates" sql:"type:jsonb" pq:",use_zero"`
		History []*TerritoryHistory `json:"history"`
	}

	TerritoryHistory struct {
		TableName struct{} `json:"-" sql:"map__territory_histories"`

		Id uint16 `json:"id"`
		TerritoryId uint16 `json:"-"`
		Territory *Territory `json:"-"`
		PlayerId uint16 `json:"-"`
		Player *Player `json:"player"`
		Action string `json:"action"`
		Data map[string]interface{} `json:"data" pg:",use_zero"`
		HappenedAt time.Time `json:"happened_at"`
	}

	PlanetTerritory struct {
		TableName struct{} `json:"-" sql:"map__planet_territories"`

		TerritoryId uint16 `json:"-"`
		Territory *Territory `json:"territory"`
		PlanetId uint16 `json:"-"`
		Planet *Planet `json:"planet"`
		Status string `json:"status"`
	}

	SystemTerritory struct {
		TableName struct{} `json:"-" sql:"map__system_territories"`

		TerritoryId uint16 `json:"-"`
		Territory *Territory `json:"territory"`
		SystemId uint16 `json:"-"`
		System *System `json:"system"`
		Status string `json:"status"`
	}
)

func GetMapTerritories(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)
	starmap := getServerMap(player.ServerId)
	
	SendJsonResponse(w, 200, starmap.getTerritories())
}

func (m *Map) getTerritories() []*Territory {
	territories := make([]*Territory, 0)
	if err := Database.Model(&territories).Relation("Planet.Player.Faction").Where("map_id = ?", m.Id).Select(); err != nil {
		panic(NewException("Could not retrieve map territories", err))
	}
	return territories
}

func (p *Planet) createTerritory() *Territory {
	t := &Territory{
		MapId: p.System.MapId,
		Map: p.System.Map,
		PlanetId: p.Id,
		Planet: p,
		MilitaryInfluence: 20,
		EconomicInfluence: 20,
		PoliticalInfluence: 20,
		ReligiousInfluence: 20,
		CulturalInfluence: 20,
		Coordinates: make([]*Coordinates, 0),
	}
	if err := Database.Insert(t); err != nil {
		panic(NewException("Could not create territory", err))
	}
	t.addHistory(p.Player, TerritoryActionCreation, make(map[string]interface{}, 0))
	p.addTerritory(t, TerritoryStatusPledge)
	p.System.addTerritory(t)
	t.expand()
	return t
}

func (p *Planet) addTerritory(t *Territory, status string) {
	pt := &PlanetTerritory{
		TerritoryId: t.Id,
		Territory: t,
		PlanetId: p.Id,
		Planet: p,
		Status: status,
	}
	if err := Database.Insert(pt); err != nil {
		panic(NewException("Could not create planet territory", err))
	}
	p.Territories = append(p.Territories, pt)
}

func (s *System) addTerritory(t *Territory) {
	st := &SystemTerritory{
		TerritoryId: t.Id,
		Territory: t,
		SystemId: s.Id,
		System: s,
		Status: TerritoryStatusContest,
	}
	if err := Database.Insert(st); err != nil {
		panic(NewException("Could not create system territory", err))
	}
	s.Territories = append(s.Territories, st)
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

func (t *Territory) getTotalInfluence() uint16 {
	return t.MilitaryInfluence + t.PoliticalInfluence + t.EconomicInfluence + t.ReligiousInfluence + t.CulturalInfluence
}

func (t *Territory) expand() {
	influence := t.getTotalInfluence()
	radius := math.Floor(float64(influence) / 100)

	centerX := float64(t.Planet.System.X)
	centerY := float64(t.Planet.System.Y)
	t.Coordinates = make([]*Coordinates, 0)

	for i := float64(0); i < 6; i++ {
		angle := float64(i * (float64(2) * math.Pi) / 6)
		coords := &Coordinates{
			X: centerX + (radius * math.Cos(angle)),
			Y: centerY + (radius * math.Sin(angle)),
		}
		t.Coordinates = append(t.Coordinates, coords)
	}
	t.update()
}

func (t *Territory) update() {
	if err := Database.Update(t); err != nil {
		panic(NewException("Could not update territory", err))
	}
}