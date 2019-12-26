package api

type(
	PlanetTerritory struct {
		tableName struct{} `pg:"map__planet_territories"`

		TerritoryId uint16 `json:"-"`
		Territory *Territory `json:"territory"`
		PlanetId uint16 `json:"-"`
		Planet *Planet `json:"planet"`
		Status string `json:"status"`
		MilitaryInfluence uint16 `json:"military_influence" pq:",use_zero"`
		EconomicInfluence uint16 `json:"economic_influence" pq:",use_zero"`
		CulturalInfluence uint16 `json:"cultural_influence" pq:",use_zero"`
		PoliticalInfluence uint16 `json:"political_influence" pq:",use_zero"`
		ReligiousInfluence uint16 `json:"religious_influence" pq:",use_zero"`
	}
)

func (p *Planet) addTerritory(t *Territory) {
	status := p.getNewTerritoryStatus(t)
	defaultInfluence := uint16(0)
	if status == TerritoryStatusPledge {
		defaultInfluence = 5
	}
	pt := &PlanetTerritory{
		TerritoryId: t.Id,
		Territory: t,
		PlanetId: p.Id,
		Planet: p,
		Status: status,
		MilitaryInfluence: defaultInfluence,
		EconomicInfluence: defaultInfluence,
		PoliticalInfluence: defaultInfluence,
		ReligiousInfluence: defaultInfluence,
		CulturalInfluence: defaultInfluence,
	}
	if err := Database.Insert(pt); err != nil {
		panic(NewException("Could not create planet territory", err))
	}
	p.Territories = append(p.Territories, pt)
	t.addHistory(p.Player, TerritoryActionPlanetGained, map[string]interface{}{
		"id": p.Id,
		"name": p.Name,
	})
}

func (p *Planet) getNewTerritoryStatus(t *Territory) string {
	if p.Player.FactionId != t.Planet.Player.FactionId {
		return TerritoryStatusContest
	}
	for _, pt := range p.Territories {
		if pt.Status == TerritoryStatusPledge {
			return TerritoryStatusContest
		}
	}
	return TerritoryStatusPledge
}