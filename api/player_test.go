package api

import(
	"testing"
)

// func TestGetPlayer(t *testing.T) {
// 	InitDatabaseMock()

// 	player := getPlayer(1, true)

// 	if id := player.Id; id != 1 {
// 		t.Errorf("Player id should equal 1, got %d", id)
// 	}
// }

func TestCreatePlayer(t *testing.T) {
	InitDatabaseMock()
	s := &Server{ Id: 1 }

	p := s.createPlayer("tester")

	if p.Username != "tester" {
		t.Errorf("Player username should be 'tester', got %s", p.Username)
	}
	if p.Pseudo != "tester" {
		t.Errorf("Player pseudo should be 'tester', got %s", p.Pseudo)
	}
	if p.Server.Id != s.Id {
		t.Errorf("Player server id should equal 1, got %d", p.Server.Id)
	}
}

func TestUpdateWallet(t *testing.T) {
	player := &Player{ Wallet: 10000 }

	if !player.updateWallet(2500) {
		t.Errorf("Should have been able to update player wallet")
	}
	if player.Wallet != 12500 {
		t.Errorf("Player wallet should equal 12500, got %d", player.Wallet)
	}
	if !player.updateWallet(-7500) {
		t.Errorf("Should have been able to update player wallet")
	}
	if player.Wallet != 5000 {
		t.Errorf("Player wallet should equal 5000, got %d", player.Wallet)
	}
	if player.updateWallet(-7500) {
		t.Errorf("Should not have been able to update player wallet")
	}
	if player.Wallet != 5000 {
		t.Errorf("Player wallet should equal 5000, got %d", player.Wallet)
	}
}

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