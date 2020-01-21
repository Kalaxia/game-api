package api

import(
	"testing"
)

func TestProcessPopulationGrowth(t *testing.T) {
	planet := &Planet{ Population: 1000000 }

	planet.processPopulationGrowth()

	if planet.Population != 1010000 {
		t.Errorf("Population should equal 1010000, got %d", planet.Population)
	}
}

func TestCalculatePopulationGrowth(t *testing.T) {
	planet := &Planet{ Population: 1000000 }
	
	if growth := planet.calculatePopulationGrowth(); growth / 0.01 == 1 {
		t.Errorf("Growth coefficient should equal 0.01, got %f", growth)
	}
}

func TestCalculatePopulationGrowthRate(t *testing.T) {
	planet := &Planet{ Population: 1000000 }

	if rate := planet.calculatePopulationGrowthRate(); rate != float64(0.015) {
		t.Errorf("Growth rate should equal 0.015, got %f", rate)
	}
}

func TestCalculatePopulationDeclineRate(t *testing.T) {
	planet := &Planet{ Population: 1000000 }

	if rate := planet.calculatePopulationDeclineRate(); rate != float64(0.005) {
		t.Errorf("Growth rate should equal 0.005, got %f", rate)
	}
}

func TestCalculateTaxes(t *testing.T) {
	InitDatabaseMock()

	planet := &Planet{
		Population: 1000000,
		TaxRate: 2,
		Player: &Player{
			Wallet: 100,
		},
		PublicOrder: 60,
	}

	planet.calculateTaxes()

	if planet.Player.Wallet != 300 {
		t.Errorf("Player wallet should equal 300, got %d", planet.Player.Wallet)
	}
	if planet.PublicOrder != 61 {
		t.Errorf("Public order should equal 61, got %d", planet.PublicOrder)
	}
}

func TestUpdatePublicOrder(t *testing.T) {
	planet := &Planet{ PublicOrder: 50 }

	if planet.updatePublicOrder(5); planet.PublicOrder != 55 {
		t.Errorf("Public order should equal 55, got %d", planet.PublicOrder)
	}
	if planet.updatePublicOrder(60); planet.PublicOrder != 100 {
		t.Errorf("Public order should equal 100, got %d", planet.PublicOrder)
	}
	if planet.updatePublicOrder(-105); planet.PublicOrder != 0 {
		t.Errorf("Public order should equal 0, got %d", planet.PublicOrder)
	}
}

func TestUpdateTaxRate(t *testing.T) {
	InitDatabaseMock()
	planet := &Planet{ TaxRate: 2 }

	if planet.updateTaxRate(3); planet.TaxRate != 3 {
		t.Errorf("Planet tax rate should equal 3, got %d", planet.TaxRate)
	}
}

func TestCalculatePublicOrderGrowth(t *testing.T) {
	planet := &Planet{ TaxRate: 4, PublicOrder: 50 }
	if growth := planet.calculatePublicOrderGrowth(); growth != -1 {
		t.Errorf("Planet public order growth should equal -1, got %d", growth)
	}
}