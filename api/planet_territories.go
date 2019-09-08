package api

type(
	PlanetTerritory struct {
		TableName struct{} `json:"-" sql:"map__planet_territories"`

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

func (p *Planet) addTerritory(t *Territory, status string) {
	pt := &PlanetTerritory{
		TerritoryId: t.Id,
		Territory: t,
		PlanetId: p.Id,
		Planet: p,
		Status: status,
		MilitaryInfluence: p.calculateMilitaryInfluence(),
		EconomicInfluence: p.calculateEconomicInfluence(),
		PoliticalInfluence: p.calculatePoliticalInfluence(),
		ReligiousInfluence: p.calculateReligiousInfluence(),
		CulturalInfluence: p.calculateCulturalInfluence(),
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