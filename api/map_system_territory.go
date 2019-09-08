package api

type(
	SystemTerritory struct {
		TableName struct{} `json:"-" sql:"map__system_territories"`

		TerritoryId uint16 `json:"-"`
		Territory *Territory `json:"territory"`
		SystemId uint16 `json:"-"`
		System *System `json:"system"`
		Status string `json:"status"`
	}
)

func (s *System) addTerritory(t *Territory) {
	st := &SystemTerritory{
		TerritoryId: t.Id,
		Territory: t,
		SystemId: s.Id,
		System: s,
		Status: s.getNewTerritoryStatus(t),
	}
	if err := Database.Insert(st); err != nil {
		panic(NewException("Could not create system territory", err))
	}
	s.Territories = append(s.Territories, st)
	s.addPlanetTerritories(t)
	if st.Status == TerritoryStatusPledge {
		s.Faction = t.Planet.Player.Faction
		s.FactionId = s.Faction.Id
		s.update()
	}
	t.addHistory(t.Planet.Player, TerritoryActionSystemGained, map[string]interface{}{
		"id": s.Id,
	})
}

func (s *System) addPlanetTerritories(t *Territory) {
	for _, p := range s.getPlanets() {
		if p.Player == nil || p.Player.FactionId != t.Planet.Player.FactionId {
			p.addTerritory(t, TerritoryStatusContest)
			continue
		}
		found := false
		for _, pt := range p.Territories {
			if pt.Status == TerritoryStatusPledge {
				p.addTerritory(t, TerritoryStatusPledge)
				found = true
				break
			}
		}
		if !found {
			p.addTerritory(t, TerritoryStatusContest)
		}
	}
} 

func (s *System) getNewTerritoryStatus(t *Territory) string {
	for _, t := range s.Territories {
		if t.Status == TerritoryStatusPledge {
			return TerritoryStatusContest
		}
	}
	for _, p := range s.getPlanets() {
		if p.Player != nil && p.Player.Faction.Id != t.Planet.Player.Faction.Id {
			return TerritoryStatusContest
		}
	}
	return TerritoryStatusPledge
}

func (st *SystemTerritory) getTotalInfluence() uint16 {
	influence := uint16(0)
	for _, pt := range st.getPlanetTerritories() {
		influence += pt.MilitaryInfluence + pt.PoliticalInfluence + pt.CulturalInfluence + pt.ReligiousInfluence + pt.EconomicInfluence
	}
	return influence
}

func (t *Territory) getSystemTerritories() []*SystemTerritory {
	territories := make([]*SystemTerritory, 0)
	if err := Database.Model(&territories).Where("territory_id = ?", t.Id).Select(); err != nil {
		panic(NewException("Could not retrieve system territories", err))
	}
	return territories
}

func (st *SystemTerritory) getPlanetTerritories() []*PlanetTerritory {
	territories := make([]*PlanetTerritory, 0)
	if err := Database.Model(&territories).Where("territory_id = ?", st.TerritoryId).Select(); err != nil {
		panic(NewException("Could not retrieve system planet territories", err))
	}
	return territories
}