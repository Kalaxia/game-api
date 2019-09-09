package api

import(
	"testing"
)

func TestGetCoordLimits(t *testing.T) {
	territory := getTerritoryMock()
	minX, maxX, minY, maxY := territory.getCoordLimits()

	if minX != 4 {
		t.Errorf("Minimum X coord should be 4, not %f", minX)
	}
	if maxX != 15 {
		t.Errorf("Maximum X coord should be 15, not %f", maxX)
	}
	if minY != 8 {
		t.Errorf("Minimum Y coord should be 8, not %f", minY)
	}
	if maxY != 16 {
		t.Errorf("Maximum Y coord should be 16, not %f", maxY)
	}
}

func TestIsSystemIn(t *testing.T) {
	system := &System{
		X: 12,
		Y: 10,
	}
	territory := getTerritoryMock()
	if !territory.isSystemIn(system) {
		t.Errorf("The system should be included in the territory")
	}
	system.X = 40
	if territory.isSystemIn(system) {
		t.Errorf("The system should not be included in the territory")
	}
}

func getTerritoryMock() *Territory {
	return &Territory{
		Coordinates: CoordinatesSlice{
			&Coordinates{
				X: 4,
				Y: 8,
			},
			&Coordinates{
				X: 8,
				Y: 16,
			},
			&Coordinates{
				X: 15,
				Y: 12,
			},
			&Coordinates{
				X: 13,
				Y: 9,
			},
		},
	}
}