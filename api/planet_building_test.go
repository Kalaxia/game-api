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