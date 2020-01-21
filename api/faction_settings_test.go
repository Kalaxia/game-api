package api

import(
	"testing"
)

func TestCreateFactionSettings(t *testing.T) {
	InitDatabaseMock()
	factions := []*Faction{
		&Faction{},
		&Faction{},
	}
	createFactionSettings(factions)

	if len(factions[0].Settings) == 0 {
		t.Errorf("Faction should have 5 settings, got %d", len(factions[0].Settings))
	}
	if factions[1].Settings[0].IsPublic != true {
		t.Errorf("First setting should be public")
	}
	if factions[1].Settings[1].Value != 6 {
		t.Errorf("Motion duration setting value should equal 6, got %d", factions[1].Settings[1].Value)
	}
} 