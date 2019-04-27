package shipManager

import(
    "kalaxia-game-api/manager"
	"kalaxia-game-api/model"
)

func ConquerPlanet(fleet *model.Fleet, planet *model.Planet) {
	var lastCombat *model.FleetCombat
	for _, f := range GetOrbitingFleets(planet) {
		if planet.Player.Faction.Id == f.Player.Faction.Id {
			lastCombat = Engage(fleet, &f)

			if !lastCombat.IsVictory {
				break
			}
		}
	}
	if lastCombat != nil && lastCombat.IsVictory {
		manager.UpdatePlanetOwner(planet, fleet.Player)
	}
}