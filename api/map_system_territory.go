package api

import (
	"math"
)

type(
	SystemTerritory struct {
		tableName struct{} `pg:"map__system_territories"`

		TerritoryId uint16 `json:"-" pg:",pk"`
		Territory *Territory `json:"territory"`
		SystemId uint16 `json:"-" pg:",pk"`
		System *System `json:"system"`
		Status string `json:"status"`
	}

	SystemTerritoryCacheItem struct {
		System *System `json:"system"`
		Status string `json:"status"`
		Radius float64 `json:"radius"`
		TotalInfluence uint16 `json:"total_influence"`
		PoliticalInfluence uint16 `json:"political_influence"`
		MilitaryInfluence uint16 `json:"military_influence"`
		EconomicInfluence uint16 `json:"economic_influence"`
		CulturalInfluence uint16 `json:"cultural_influence"`
		ReligiousInfluence uint16 `json:"religious_influence"`
	}
)

func (s *System) addTerritory(t *Territory) {
	newSystemTerritory := &SystemTerritory{
		TerritoryId: t.Id,
		Territory: t,
		SystemId: s.Id,
		System: s,
		Status: s.getNewTerritoryStatus(t),
	}
	if err := Database.Insert(newSystemTerritory); err != nil {
		panic(NewException("Could not create system territory", err))
	}
	s.Territories = append(s.Territories, newSystemTerritory)
	s.addPlanetTerritories(t)
	if newSystemTerritory.Status == TerritoryStatusPledge {
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
		if p.Player == nil {
			continue
		}
		p.addTerritory(t)
	}
} 

func (s *System) getNewTerritoryStatus(t *Territory) string {
	for _, t := range s.Territories {
		if t.Status == TerritoryStatusPledge {
			return TerritoryStatusContest
		}
	}
	if !s.hasRivalPlanets(t) && s.hasFriendlyPlanets(t) {
		return TerritoryStatusPledge
	}
	return TerritoryStatusContest
}

func (s *System) hasRivalPlanets(t *Territory) bool {
	for _, p := range s.Planets {
		if f := p.getFaction(); f != nil && f.Id != t.Planet.Player.Faction.Id {
			return true
		}
	}
	return false
}

func (s *System) hasFriendlyPlanets(t *Territory) bool {
	for _, p := range s.Planets {
		if f := p.getFaction(); f != nil && f.Id == t.Planet.Player.Faction.Id {
			return true
		}
	}
	return false
}

func (s *System) checkTerritories() {
	s.Territories = s.getTerritories()
	if len(s.Territories) == 0 {
		return
	}
	var dominantTerritory *SystemTerritory
	dominantInfluence := uint16(0)

	for _, st := range s.Territories {
		st.Status = TerritoryStatusContest
		if influence := st.getTotalInfluence(); influence > dominantInfluence {
			dominantInfluence = influence
			dominantTerritory = st
		}
	}
	if dominantTerritory == nil {
		return
	}
	dominantTerritory.Status = TerritoryStatusPledge
	s.Faction = dominantTerritory.Territory.Planet.Player.Faction
	s.FactionId = s.Faction.Id
	s.update()

	for _, st := range s.Territories {
		st.update()
		st.Territory.expand()
	}
}

func (s *System) getTerritories() []*SystemTerritory {
	territories := make([]*SystemTerritory, 0)
	if err := Database.Model(&territories).Relation("Territory.Planet.Player.Faction").Relation("Territory.Planet.System").Where("system_territory.system_id = ?", s.Id).Select(); err != nil {
		panic(NewException("Could not retrieve system territories", err))
	}
	return territories
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
	if err := Database.Model(&territories).Relation("System").Where("territory_id = ?", t.Id).Select(); err != nil {
		panic(NewException("Could not retrieve system territories", err))
	}
	return territories
}

func (st *SystemTerritory) getPlanetTerritories() []*PlanetTerritory {
	territories := make([]*PlanetTerritory, 0)
	if err := Database.Model(&territories).Relation("Planet.Player").Where("Planet.system_id = ?", st.SystemId).Where("territory_id = ?", st.TerritoryId).Select(); err != nil {
		panic(NewException("Could not retrieve system planet territories", err))
	}
	return territories
}

func (st *SystemTerritory) getRadius() float64 {
	return math.Sqrt(float64(st.getTotalInfluence() / 10) / math.Pi)
}

func (st *SystemTerritory) checkForIncludedSystems() (hasNewSystem bool) {
	minX, maxX, minY, maxY := st.getCoordLimits()

	systems := make([]*System, 0)
	if err := Database.
		Model(&systems).
		Relation("Territories").
		Relation("Planets.Player.Faction").
		Where("system.x <= ?", maxX).
		Where("system.x >= ?", minX).
		Where("system.y <= ?", maxY).
		Where("system.y >= ?", minY).
		Where("system.id != ?", st.System.Id).
		Select(); err != nil {
		panic(NewException("Could not retrieve included systems", err))
	}
	for _, s := range systems {
		if st.canAnnex(s) {
			hasNewSystem = true
			s.addTerritory(st.Territory)
		}
	}
	return
}

func (st *SystemTerritory) getCoordLimits() (minX, maxX, minY, maxY float64) {
	radius := st.getRadius()

	minX = float64(st.System.X) - radius
    maxX = float64(st.System.X) + radius
    minY = float64(st.System.Y) - radius
    maxY = float64(st.System.Y) + radius

	return
}

func (st *SystemTerritory) isSystemIn(s *System) bool {
	return st.getRadius() >= getDistance(float64(st.System.X), float64(s.X), float64(st.System.Y), float64(s.Y))
}

func (st *SystemTerritory) canAnnex(s *System) bool {
	return st.isSystemIn(s) && !st.hasSystem(s)
}

func (st *SystemTerritory) hasSystem(s *System) bool {
	for _, t := range s.Territories {
		if t.TerritoryId == st.TerritoryId {
			return true
		}
	}
	return false
}

func (st *SystemTerritory) generateCache() *SystemTerritoryCacheItem {
	stci := &SystemTerritoryCacheItem{
		System: st.System,
		Status: st.Status,
		Radius: st.getRadius(),
		TotalInfluence: 0,
		PoliticalInfluence: 0,
		MilitaryInfluence: 0,
		EconomicInfluence: 0,
		CulturalInfluence: 0,
		ReligiousInfluence: 0,
	}
	for _, pt := range st.getPlanetTerritories() {
		stci.TotalInfluence += pt.PoliticalInfluence + pt.MilitaryInfluence + pt.EconomicInfluence + pt.CulturalInfluence + pt.ReligiousInfluence
		stci.PoliticalInfluence += pt.PoliticalInfluence
		stci.MilitaryInfluence += pt.MilitaryInfluence
		stci.EconomicInfluence += pt.EconomicInfluence
		stci.CulturalInfluence += pt.CulturalInfluence
		stci.ReligiousInfluence += pt.ReligiousInfluence
	}
	return stci
}

func (st *SystemTerritory) update() {
	if err := Database.Update(st); err != nil {
		panic(NewException("Could not update system territory", err))
	}
}