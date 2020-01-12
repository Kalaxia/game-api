package api

import(
	"testing"
)

func TestGetResource(t *testing.T) {
	planet := getPlayerPlanetMock(getPlayerMock(getFactionMock()))

	if geode := planet.getResource("geode"); geode != nil {
		t.Errorf("Planet should not have geode resource")
	}
	ore := planet.getResource("ore")
	if ore == nil {
		t.Errorf("Planet should have ore resource")
	}
	if ore.Density != 45 {
		t.Errorf("Planet ore density should equal 45, not %d", ore.Density)
	}
}

func getPlayerPlanetMock(player *Player) *Planet {
	return &Planet{
		Id: 1,
		Name: "Régalion V",
		Type: PlanetTypeArctic,
		Player: player,
		PlayerId: player.Id,
		Population: 20 * populationPointRatio,
		Settings: &PlanetSettings{
			ServicesPoints: 5,
			MilitaryPoints: 7,
			BuildingPoints: 3,
			ResearchPoints: 2,
		},
		System: &System{
			X: 25,
			Y: 47,
		},
		Resources: []PlanetResource{
			PlanetResource{
				Name: "cristal",
				Density: 60,
				PlanetId: 1,
			},
			PlanetResource{
				Name: "red-ore",
				Density: 20,
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