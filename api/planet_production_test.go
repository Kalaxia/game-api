package api

import(
	"reflect"
	"testing"
	_ "kalaxia-game-api/api/mock"
)

func TestCalculatePlanetResourcesProduction(t *testing.T) {
	reflect.ValueOf(Database).Elem().FieldByName("NextId").SetUint(1)
	planet := getPlayerPlanetMock()
	planet.Storage = getStorageMock()

	planet.produceResources()

	if len(planet.Storage.Resources) != 3 {
		t.Fatalf("Storage should contain 3 different resources")
	}
	if planet.Storage.Resources["ore"] != 450 {
		t.Errorf("Storage should contain 450 ore")
	}
	if planet.Storage.Resources["cristal"] != 3100 {
		t.Errorf("Storage should contain 3100 cristal")
	}
}

func TestCalculatePlanetResourcesProductionWithoutStorage(t *testing.T) {
	reflect.ValueOf(Database).Elem().FieldByName("NextId").SetUint(1)
	planet := getPlayerPlanetMock()

	planet.produceResources()

	if planet.Storage == nil {
		t.Fatalf("Storage should have been created")
	}
	if len(planet.Storage.Resources) != 2 {
		t.Fatalf("Storage should contain 2 different resources")
	}
	if planet.Storage.Resources["cristal"] != 600 {
		t.Errorf("Storage should contain 600 cristal")
	}
}