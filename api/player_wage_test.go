package api

import(
	"testing"
)

func TestGetServiceWage(t *testing.T) {
	planet := &Planet{ Settings: &PlanetSettings{ ServicesPoints: 10 }}

	if wage := planet.getServiceWage(); wage != 75 {
		t.Errorf("Planet wage should equal 75, got %d", wage)
	}
}