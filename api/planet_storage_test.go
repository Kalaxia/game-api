package api

import(
	"testing"
)

func TestHasResource(t *testing.T) {
	storage := getStorageMock()

	if storage.hasResource("cristal", 2700) {
		t.Errorf("Storage has not enough cristal")
	}
	if storage.hasResource("red-ore", 500) {
		t.Errorf("Storage has no red-ore")
	}
	if !storage.hasResource("geode", 1000) {
		t.Errorf("Storage has geode")
	}
}

func TestStoreResourceProduction(t *testing.T) {
	planet := getPlayerPlanetMock(getPlayerMock(getFactionMock()))
	storage := getStorageMock()

	storage.storeResourceProduction(planet)

	if storage.Resources["ore"] != 450 {
		t.Errorf("Storage has 450 ore")
	}
	if storage.Resources["cristal"] != 3100 {
		t.Errorf("Storage has 3100 cristal")
	}
}

func TestStoreResource(t *testing.T) {
	storage := getStorageMock()

	if !storage.storeResource("red-ore", 6000) {
		t.Errorf("Storage could store red-ore")
		if storage.Resources["red-ore"] != 5000 {
			t.Errorf("Storage has 5000 red-ore")
		}
	}
	if storage.storeResource("geode", -2000) {
		t.Errorf("Storage cannot spend 2000 geodes")
	}
}

func getStorageMock() *Storage {
	return &Storage{
		Id: 1,
		Capacity: 5000,
		Resources: map[string]uint16{
			"cristal": 2500,
			"geode": 1450,
		},
	}
}