package api

func getPlayerPlanetMock(player *Player) *Planet {
	return &Planet{
		Id: 1,
		Name: "Régalion V",
		Type: PlanetTypeArctic,
		Player: player,
		PlayerId: player.Id,
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