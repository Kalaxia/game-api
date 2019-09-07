package api

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
	if st.Status == TerritoryStatusPledge {
		s.Faction = t.Planet.Player.Faction
		s.FactionId = s.Faction.Id
		s.update()
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