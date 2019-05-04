package api

import(
	"reflect"
	"testing"
)

func TestCalculatePlanetResourcesProduction(t *testing.T) {
	InitDatabaseMock()
	reflect.ValueOf(Database).Elem().FieldByName("NextId").SetUint(1)
	planet := getPlayerPlanetMock(getPlayerMock(getFactionMock()))
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
	InitDatabaseMock()
	reflect.ValueOf(Database).Elem().FieldByName("NextId").SetUint(1)
	planet := getPlayerPlanetMock(getPlayerMock(getFactionMock()))

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