package manager

import(
	"testing"
	_ "kalaxia-game-api/database/mock"
	"kalaxia-game-api/exception"
	"kalaxia-game-api/model"
)

func TestUpdatePlanetSettings(t *testing.T) {
	planet := getPlanetMock()
	UpdatePlanetSettings(planet, getSettingsMock())

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

func TestUpdatePlanetSettingsWithTooMuchPoints(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("An error should have been thrown")
		}
		if exception, ok := r.(*exception.HttpException); ok {
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
	planet := getPlanetMock()
	settings := getSettingsMock()
	settings.MilitaryPoints = 15

	UpdatePlanetSettings(planet, settings)
}

func getPlanetMock() *model.Planet {
	return &model.Planet{
		Id: 1,
		Name: "Régalion V",
		Type: model.PlanetTypeArctic,
		Population: 2000000,
		Settings: &model.PlanetSettings{
			ServicesPoints: 5,
			MilitaryPoints: 7,
			BuildingPoints: 3,
			ResearchPoints: 2,
		},
	}
}

func getSettingsMock() *model.PlanetSettings {
	return &model.PlanetSettings{
		ServicesPoints: 3,
		MilitaryPoints: 5,
		BuildingPoints: 7,
		ResearchPoints: 5,
	}
}