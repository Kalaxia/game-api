package api

import(
	"fmt"
	"reflect"
	"testing"
)

func TestProduceResources(t *testing.T) {
	InitDatabaseMock()
	InitPlanetConstructions()
	reflect.ValueOf(Database).Elem().FieldByName("NextId").SetUint(1)
	planet := getPlayerPlanetMock(getPlayerMock(getFactionMock()))
	planet.Storage = getStorageMock()
	building := &Building{
		Planet: planet,
		Type: "resource",
		Name: "ore-mine",
	}
	building2 := Building{
		Planet: planet,
		Type: "resource",
		Name: "cristal-synthetiser",
	}
	building.produceResources()

	fmt.Println(planet.Storage.Resources)

	if len := len(planet.Storage.Resources); len != 4 {
		t.Fatalf("Storage should contain 3 different resources, got %d", len)
	}
	if planet.Storage.Resources["red-ore"] != 200 {
		t.Errorf("Storage should contain 200 red-ore")
	}
	if planet.Storage.Resources["ore"] != 450 {
		t.Errorf("Storage should contain 450 ore")
	}
	if planet.Storage.Resources["cristal"] != 2500 {
		t.Errorf("Storage should contain 2500 cristal")
	}

	building2.produceResources()

	if planet.Storage.Resources["cristal"] != 3100 {
		t.Errorf("Storage should contain 3100 cristal")
	}
}

func TestGetProducedQuantity(t *testing.T) {
	building := &Building{
		Planet: getPlayerPlanetMock(getPlayerMock(getFactionMock())),
	}

	if quantity := building.getProducedQuantity("red-ore"); quantity != 200 {
		t.Errorf("Produced ore quantity should equal 250, not %d", quantity)
	}
	if quantity := building.getProducedQuantity("emerald"); quantity != 0 {
		t.Errorf("Produced emerald quantity should equal 0, not %d", quantity)
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


func TestSpendPoints(t *testing.T) {
	InitDatabaseMock()
	building := &Building{
		Status: BuildingStatusConstructing,
		ConstructionStateId: 1,
		ConstructionState: &PointsProduction{
			Id: 1,
			Points: 10,
			CurrentPoints: 2,
		},
	}
	if points := building.ConstructionState.spendPoints(5, building.finishConstruction); points != 0 {
		t.Errorf("Remaining points should equal 0, not %d", points)
	}
	if building.ConstructionState.CurrentPoints != 7 {
		t.Errorf("Current points should equal 7, not %d", building.ConstructionState.CurrentPoints)
	}
	if building.Status != BuildingStatusConstructing {
		t.Errorf("Building status should be constructing")
	}
	if points := building.ConstructionState.spendPoints(5, building.finishConstruction); points != 2 {
		t.Errorf("Remaining points should equal 2, not %d", points)
	}
	if building.Status != BuildingStatusOperational {
		t.Errorf("Building status should be operational")
	}
	if building.ConstructionStateId != 0 {
		t.Errorf("Building construction state ID should equal 0")
	}
}