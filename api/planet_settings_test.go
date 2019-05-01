package api

import(
	"testing"
)

func TestAffectPopulationPoints(t *testing.T) {
	InitDatabaseMock()
	planet := getPlayerPlanetMock()
	planet.affectPopulationPoints(getSettingsMock())

	if planet.Settings.ServicesPoints != 3 {
		t.Errorf("Services points should be set to 3")
	}
	if planet.Settings.MilitaryPoints != 5 {
		t.Errorf("Military points should be 5")
	}
	if planet.Settings.BuildingPoints != 7 {
		t.Errorf("Building points should be 7")
	}
	if planet.Settings.ResearchPoints != 5 {
		t.Errorf("Research points should be 5")
	}
}

func TestAffectPopulationPointsWithTooMuchPoints(t *testing.T) {
	InitDatabaseMock()
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("An error should have been thrown")
		}
		if exception, ok := r.(*HttpException); ok {
			if exception.Code != 400 {
				t.Errorf("The error code should be 400")
			}
			if exception.Message != "Not enough population points" {
				t.Errorf("The error should say there is not enough population points")
			}
		} else {
			t.Errorf("The error should be a HTTP Exception")
		}
	}()
	planet := getPlayerPlanetMock()
	settings := getSettingsMock()
	settings.MilitaryPoints = 15

	planet.affectPopulationPoints(settings)
}

func getSettingsMock() *PlanetSettings {
	return &PlanetSettings{
		ServicesPoints: 3,
		MilitaryPoints: 5,
		BuildingPoints: 7,
		ResearchPoints: 5,
	}
}