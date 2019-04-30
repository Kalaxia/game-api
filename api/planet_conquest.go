package api

func (fleet *Fleet) conquerPlanet(planet *Planet) {
	var lastCombat *FleetCombat
	for _, f := range planet.getOrbitingFleets() {
		if planet.Player.Faction.Id == f.Player.Faction.Id {
			lastCombat = fleet.engage(&f)

			if !lastCombat.IsVictory {
				break
			}
		}
	}
	if lastCombat != nil && lastCombat.IsVictory {
		planet.changeOwner(fleet.Player)
		planet.update()
	}
}