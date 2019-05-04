package api

func getPlayerMock(faction *Faction) *Player {
	return &Player{
		Id: 1,
		Username: "Tester",
		Pseudo: "Testinator",
		Wallet: 1000,
		Gender: GenderMale,
		Faction: faction,
		FactionId: faction.Id,
	}
}