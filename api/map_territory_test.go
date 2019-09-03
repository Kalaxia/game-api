package api

import(
	"testing"
)

func TestConvexHull(t *testing.T) {
	coords := CoordinatesSlice{
		&Coordinates{
			X: 47.03,
			Y: 23.12,
		},
		&Coordinates{
			X: 80,
			Y: 70,
		},
		&Coordinates{
			X: 40.16,
			Y: 20,
		},
		&Coordinates{
			X: 55.17,
			Y: 37.03,
		},
		&Coordinates{
			X: 96.23,
			Y: 14.07,
		},
		&Coordinates{
			X: 15,
			Y: 20,
		},
	}
	territory := &Territory{
		Coordinates: coords,
	}
	territory.convexHull()

	if nbCoords := len(territory.Coordinates); nbCoords != 6 {
		t.Errorf("Territory should have 6 Coordinates, not %d", nbCoords)
	}
}