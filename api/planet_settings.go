package api

import (
    "net/http"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "strconv"
    "math"
)

type(
	PlanetSettings struct {
		TableName struct{} `json:"-" sql:"map__planet_settings"`
  
		Id uint16 `json:"-"`
		ServicesPoints uint8 `json:"services_points" sql:",notnull"`
		BuildingPoints uint8 `json:"building_points" sql:",notnull"`
		MilitaryPoints uint8 `json:"military_points" sql:",notnull"`
		ResearchPoints uint8 `json:"research_points" sql:",notnull"`
	}
)

func AffectPopulationPoints(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)

    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := getPlayerPlanet(uint16(id), player.Id)

    if player.Id != planet.Player.Id {
        panic(NewHttpException(http.StatusForbidden, "", nil))
    }
    data := DecodeJsonRequest(r)
    settings := &PlanetSettings{
        ServicesPoints: uint8(data["services_points"].(float64)),
        BuildingPoints: uint8(data["building_points"].(float64)),
        MilitaryPoints: uint8(data["military_points"].(float64)),
        ResearchPoints: uint8(data["research_points"].(float64)),
    }
    planet.affectPopulationPoints(settings)
    SendJsonResponse(w, 200, settings)
}

func (p *Planet) affectPopulationPoints(settings *PlanetSettings) {
    if settings.ServicesPoints +
    settings.BuildingPoints +
    settings.MilitaryPoints +
    settings.ResearchPoints > calculatePopulationPoints(p) {
        panic(NewHttpException(400, "Not enough population points", nil))
    }
    p.Settings.ServicesPoints = settings.ServicesPoints
    p.Settings.BuildingPoints = settings.BuildingPoints
    p.Settings.MilitaryPoints = settings.MilitaryPoints
    p.Settings.ResearchPoints = settings.ResearchPoints

    if err := Database.Update(p.Settings); err != nil {
        panic(NewException("Planet settings could not be updated", err))
    }
}

func calculatePopulationPoints(planet *Planet) uint8 {
    return uint8(math.Ceil(float64(planet.Population / 100000)))
}
