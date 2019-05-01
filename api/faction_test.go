package api

import(
	"reflect"
	"testing"
)

func TestCreateServerFactions(t *testing.T) {
	InitDatabaseMock()
	reflect.ValueOf(Database).Elem().FieldByName("NextId").SetUint(1)
	factions := getServerMock().createFactions(getFactionsMock())

	if len(factions) != 3 {
		t.Errorf("Factions array should have three items")
	}
	faction := factions[0]
	if faction.Name != "Les Kalankars" {
		t.Errorf("Faction name should be the kalankars")
	}
	if faction.Slug != "les-kalankars" {
		t.Errorf("Faction should have been slugified")
	}
	if faction.Colors["main"] != "#2D2D2D" {
		t.Errorf("Faction main color should be set")
	}
	if faction.Server.Id != 1 {
		t.Errorf("Factions should have Server struct set properly")
	}
	if len(faction.Relations) != 2 {
		t.Errorf("Faction should have relations set with the two others")
	}
	if faction.Relations[0].Faction != faction {
		t.Errorf("Faction relation should have its own Faction structure stored")
	}
	if faction.Relations[0].State != RelationNeutral {
		t.Errorf("Faction relation should be neutral")
	}
}

func getServerMock() *Server {
	return &Server{
		Id: 1,
		Name: "The test server",
		Type: "Multiplayer",
		Signature: "qd15sd14ff64d5sqd",
	}
}

func getFactionsMock() []interface{} {
	factions := make([]interface{}, 3)
	factions[0] = map[string]interface{}{
		"colors": map[string]interface{} {
			"main": "#2D2D2D",
		},
		"name": "Les Kalankars",
		"banner": "kalankars.png",
		"description": "L'empire du bon vieux temps",
	}
	factions[1] = map[string]interface{}{
		"colors": map[string]interface{} {
			"main": "#2D2D2D",
		},
		"name": "L'Ascendance Valkar",
		"banner": "valkars.png",
		"description": "L'empire comme on l'aime",
	}
	factions[2] = map[string]interface{}{
		"colors": map[string]interface{} {
			"main": "#2D2D2D",
		},
		"name": "Les Adranites",
		"banner": "adranites.png",
		"description": "Les envahisseurs inconnus",
	}
	return factions
}