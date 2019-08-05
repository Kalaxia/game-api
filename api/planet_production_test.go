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

func TestCreatePointsProduction(t *testing.T) {
	InitDatabaseMock()
	planet := getPlayerPlanetMock(getPlayerMock(getFactionMock()))
	pp := planet.createPointsProduction(10)

	if pp.CurrentPoints != 0 {
		t.Errorf("Current points should equal 0, not %d", pp.CurrentPoints)
	}
	if pp.Points != 10 {
		t.Errorf("Points should equal 10, not %d", pp.Points)
	}
}

func TestIsCompleted(t *testing.T) {
	pp := &PointsProduction{
		CurrentPoints: 10,
		Points: 10,
	}
	if !pp.isCompleted() {
		t.Errorf("Points production should be complete")
	}
	pp.CurrentPoints = 8
	if pp.isCompleted() {
		t.Errorf("Points production should not be complete")
	}
}

func TestGetMissingPoints(t *testing.T) {
	pp := &PointsProduction{
		CurrentPoints: 8,
		Points: 12,
	}
	if missingPoints := pp.getMissingPoints(); missingPoints != 4 {
		t.Errorf("Missing points should equal 4, not %d", missingPoints)
	}
}