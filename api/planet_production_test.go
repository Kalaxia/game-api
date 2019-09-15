package api

import(
	"reflect"
	"testing"
)

func TestGetPlanetProducedResources(t *testing.T) {
	InitDatabaseMock()
	InitPlanetConstructions()
	reflect.ValueOf(Database).Elem().FieldByName("NextId").SetUint(1)
	planet := getPlayerPlanetMock(getPlayerMock(getFactionMock()))
	planet.Storage = getStorageMock()
	planet.Buildings = []Building{
		Building{
			Planet: planet,
			Type: "resource",
			Name: "ore-mine",
			Compartments: []*BuildingCompartment{
				&BuildingCompartment{
					Name: "ore-well",
					Status: BuildingStatusOperational,
				},
			},
		},
		Building{
			Planet: planet,
			Type: "resource",
			Name: "cristal-synthetiser",
		},
	}
	planetResources := planet.getProducedResources()

	if len := len(planetResources); len != 3 {
		t.Fatalf("Storage should contain 3 different resources, got %d", len)
	}
	if quantity := planetResources["red-ore"].BaseQuantity; quantity != 200 {
		t.Errorf("Storage should contain 200 red-ore, not %d", quantity)
	}
	if percent := planetResources["red-ore"].Percent; percent != 10 {
		t.Errorf("Produced red-ore should be improved by 10 percents, got %d", percent)
	}
	if quantity := planetResources["red-ore"].FinalQuantity; quantity != 220 {
		t.Errorf("Finally produced red-ore should equal 220, not %d", quantity)
	}
	if quantity := planetResources["ore"].BaseQuantity; quantity != 450 {
		t.Errorf("Storage should contain 450 ore, not %d", quantity)
	}
	if percent := planetResources["ore"].Percent; percent != 10 {
		t.Errorf("Produced ore should be improved by 10 percents, got %d", percent)
	}
	if quantity := planetResources["ore"].FinalQuantity; quantity != 495 {
		t.Errorf("Finally produced ore should equal 495, not %d", quantity)
	}
	if quantity := planetResources["cristal"].BaseQuantity; quantity != 600 {
		t.Errorf("Storage should contain 600 cristal, not %d", quantity)
	}
	if percent := planetResources["cristal"].Percent; percent != 0 {
		t.Errorf("Produced cristal should not be affected, got %d", percent)
	}
	if quantity := planetResources["cristal"].FinalQuantity; quantity != 600 {
		t.Errorf("Finally produced cristal should equal 495, not %d", quantity)
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