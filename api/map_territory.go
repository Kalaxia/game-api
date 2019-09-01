package api

import(
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
		Data map[string]interface{} `json:"data"`
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
	}
	if err := Database.Insert(t); err != nil {
		panic(NewException("Could not create territory", err))
	}
	t.addHistory(p.Player, TerritoryActionCreation, make(map[string]interface{}))
	p.addTerritory(t, TerritoryStatusPledge)
	p.System.addTerritory(t)
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