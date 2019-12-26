package api

import(
	"math"
)

type(
	PlanetTerritory struct {
		tableName struct{} `pg:"map__planet_territories"`

		TerritoryId uint16 `json:"-" pg:",pk"`
		Territory *Territory `json:"territory"`
		PlanetId uint16 `json:"-" pg:",pk"`
		Planet *Planet `json:"planet"`
		Status string `json:"status"`
		MilitaryInfluence uint16 `json:"military_influence" pg:",notnull,use_zero"`
		EconomicInfluence uint16 `json:"economic_influence" pg:",notnull,use_zero"`
		CulturalInfluence uint16 `json:"cultural_influence" pg:",notnull,use_zero"`
		PoliticalInfluence uint16 `json:"political_influence" pg:",notnull,use_zero"`
		ReligiousInfluence uint16 `json:"religious_influence" pg:",notnull,use_zero"`
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

func (p *Planet) getTerritories() []*PlanetTerritory {
	territories := make([]*PlanetTerritory, 0)
	if err := Database.Model(&territories).Relation("Territory.Planet.Player.Faction").Where("planet_territory.planet_id = ?", p.Id).Select(); err != nil {
		panic(NewException("Could not retrieve planet territories", err))
	}
	return territories
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

func (p *Planet) addInfluence(player *Player, points uint16, modifier func(pt *PlanetTerritory, influence uint16)) {
	for _, pt := range p.getTerritories() {
		influence := points
		if pt.Territory.Planet.Player.FactionId != player.FactionId {
			influence = uint16(math.Ceil(float64(points) / 5))
		}
		if pt.Territory.Planet.PlayerId == player.Id {
			influence = uint16(math.Ceil(float64(points) * 1.5))
		}
		modifier(pt, influence)
		pt.update()
	}
	p.System.checkTerritories()
}

func (pt *PlanetTerritory) update() {
	if err := Database.Update(pt); err != nil {
		panic(NewException("Could not update planet territory", err))
	}
}