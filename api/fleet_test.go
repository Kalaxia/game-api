package api

// import (
// 	"testing"
// 	"reflect"
// )

// func TestCreateFleet(t *testing.T) {
// 	InitDatabaseMock()
// 	reflect.ValueOf(Database).Elem().FieldByName("NextId").SetUint(1)

// 	player := getPlayerMock(getFactionMock())
// 	planet := getPlayerPlanetMock(player)

// 	fleet := player.createFleet(planet)

// 	if structName := reflect.TypeOf(fleet).Elem().Name(); structName != "Fleet" {
// 		t.Fatalf("Result should be a Fleet pointer, have %s", structName)
// 	}
// 	if fleet.Player != player {
// 		t.Errorf("Fleet should be assigned to the given player")
// 	}
// 	if fleet.Place.Planet != planet {
// 		t.Errorf("Fleet should be located on the given planet")
// 	}
// 	if fleet.Journey != nil {
// 		t.Errorf("Fleet should not be on journey")
// 	}
// }