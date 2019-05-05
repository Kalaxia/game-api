package api

func (fleet *Fleet) conquerPlanet(planet *Planet) {
	var lastCombat *FleetCombat
	for _, f := range planet.getOrbitingFleets() {
		if planet.Player.Faction.Id == f.Player.Faction.Id && f.Player.Faction.Id != fleet.Player.Faction.Id {
			lastCombat = fleet.engage(&f)

			if lastCombat != nil && !lastCombat.IsVictory {
				break
			}
		}
	}
	if lastCombat == nil || lastCombat.IsVictory {
		fleet.Player.notify(
			NotificationTypeMilitary,
			"notifications.military.planet_conquest",
			map[string]interface{}{
				"planet": planet,
				"opponent": planet.Player,
			},
		)
		if planet.Player != nil {
			planet.Player.notify(
				NotificationTypeMilitary,
				"notifications.military.planet_conquerred",
				map[string]interface{}{
					"planet": planet,
					"opponent": fleet.Player,
				},
			)
		}
		planet.changeOwner(fleet.Player)
		planet.update()
	}
}