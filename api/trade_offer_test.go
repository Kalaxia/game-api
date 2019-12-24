package api

import(
	"testing"
)

func TestCancelOffer(t *testing.T) {
	InitWsHubMock()
	InitDatabaseMock()

	offer := getResourceOfferMock()
	offer.cancel()

	if cristal := offer.Location.Storage.Resources["cristal"]; cristal != 3000 {
		t.Errorf("Planet storage should contain 3000 cristals, got %d", cristal)
	}
}

func TestApplyTaxes(t *testing.T) {
	finalPrice, gain, purchaseTaxes, salesTaxes := applyTaxes(500, 5, 10)

	if finalPrice != 525 {
		t.Errorf("Final price should equal 525, not %d", finalPrice)
	}
	if gain != 450 {
		t.Errorf("Gain should equal 450, not %d", gain)
	}
	if purchaseTaxes != 25 {
		t.Errorf("Purchase taxes should equal 25, not %d", purchaseTaxes)
	}
	if salesTaxes != 50 {
		t.Errorf("Sales taxes should equal 50, not %d", salesTaxes)
	}
}

func getResourceOfferMock() *ResourceOffer {
	location := getPlayerPlanetMock(getPlayerMock(getFactionMock()))
	location.Storage = getStorageMock()
	return &ResourceOffer{
		Offer: Offer{
			Location: location,
			Price: 10,
		},
		Resource: "cristal",
		Quantity: 500,
	}
}