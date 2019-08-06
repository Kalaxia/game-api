package api

import(
	"testing"
)

func TestPayPrice(t *testing.T) {
	InitDatabaseMock()

	planet := getPlayerPlanetMock(getPlayerMock(getFactionMock()))
	buildingPlan := BuildingPlan{
		Price: []Price{
			Price{
				Type: PriceTypePoints,
				Amount: 10,
			},
			Price{
				Type: PriceTypeMoney,
				Amount: 500,
			},
		},
	}

	points := planet.payPrice(buildingPlan.Price, 1)

	if points != 10 {
		t.Errorf("The price in points should equal 10, not %d", points)
	}
	if planet.Player.Wallet != 500 {
		t.Errorf("Player wallet should equal 500, not %d", planet.Player.Wallet)
	}
}

func TestSpendPoints(t *testing.T) {
	InitDatabaseMock()
	building := &Building{
		ConstructionStateId: 1,
		ConstructionState: &PointsProduction{
			Id: 1,
			Points: 10,
			CurrentPoints: 2,
		},
	}
	if points := building.spendPoints(5); points != 0 {
		t.Errorf("Remaining points should equal 0, not %d", points)
	}
	if building.ConstructionState.CurrentPoints != 7 {
		t.Errorf("Current points should equal 7, not %d", building.ConstructionState.CurrentPoints)
	}
	if points := building.spendPoints(5); points != 2 {
		t.Errorf("Remaining points should equal 2, not %d", points)
	}
	if building.Status != BuildingStatusOperational {
		t.Errorf("Building status should be operational")
	}
	if building.ConstructionStateId != 0 {
		t.Errorf("Building construction state ID should equal 0")
	}
}