package shipManager

import(
    "kalaxia-game-api/manager"
	"kalaxia-game-api/model"
)

func ConquerPlanet(fleet *model.Fleet, planet *model.Planet) {
	manager.UpdatePlanetOwner(planet, fleet.Player)
}