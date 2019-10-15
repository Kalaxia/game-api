package api

type(
	PlanetTerritory struct {
		tableName struct{} `json:"-" pg:"map__planet_territories"`

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
	pt := &PlanetTerritory{
		TerritoryId: t.Id,
		Territory: t,
		PlanetId: p.Id,
		Planet: p,
		Status: p.getNewTerritoryStatus(t),
	}
	pt.calculateInfluence()
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

func (pt *PlanetTerritory) calculateInfluence() {
	if pt.Planet.Player.FactionId == pt.Territory.Planet.Player.FactionId {
		pt.MilitaryInfluence = pt.Planet.calculateMilitaryInfluence()
		pt.EconomicInfluence = pt.Planet.calculateEconomicInfluence()
		pt.PoliticalInfluence = pt.Planet.calculatePoliticalInfluence()
		pt.ReligiousInfluence = pt.Planet.calculateReligiousInfluence()
		pt.CulturalInfluence = pt.Planet.calculateCulturalInfluence()
	} else {
		pt.MilitaryInfluence = 0
		pt.EconomicInfluence = 0
		pt.PoliticalInfluence = 0
		pt.ReligiousInfluence = 0
		pt.CulturalInfluence = 0
	}
}

func (p *Planet) calculateMilitaryInfluence() uint16 {
	return 5
}

func (p *Planet) calculateEconomicInfluence() uint16 {
	return 5
}

func (p *Planet) calculatePoliticalInfluence() uint16 {
	return 5
}

func (p *Planet) calculateReligiousInfluence() uint16 {
	return 5
}

func (p *Planet) calculateCulturalInfluence() uint16 {
	return 5
}