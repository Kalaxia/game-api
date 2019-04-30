package api

import(
	"reflect"
	"testing"
	_ "kalaxia-game-api/api/mock"
)

func getPlayerPlanetMock() *Planet {
	return &Planet{
		Id: 1,
		Name: "Régalion V",
		Type: PlanetTypeArctic,
		Population: 2000000,
		Settings: &PlanetSettings{
			ServicesPoints: 5,
			MilitaryPoints: 7,
			BuildingPoints: 3,
			ResearchPoints: 2,
		},
		Resources: []PlanetResource{
			PlanetResource{
				Name: "cristal",
				Density: 60,
				PlanetId: 1,
			},
			PlanetResource{
				Name: "ore",
				Density: 45,
				PlanetId: 1,
			},
		},
	}
}