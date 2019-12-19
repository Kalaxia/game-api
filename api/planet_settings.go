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
		tableName struct{} `json:"-" pg:"map__planet_settings"`
  
		Id uint16 `json:"-"`
		ServicesPoints uint8 `json:"services_points" pg:",notnull,use_zero"`
		BuildingPoints uint8 `json:"building_points" pg:",notnull,use_zero"`
		MilitaryPoints uint8 `json:"military_points" pg:",notnull,use_zero"`
		ResearchPoints uint8 `json:"research_points" pg:",notnull,use_zero"`
	}
)

func AffectPopulationPoints(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)

    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := player.getPlanet(uint16(id))

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
    settings.ResearchPoints > p.calculatePopulationPoints() {
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

func (p *Planet) calculatePopulationPoints() uint8 {
    return uint8(math.Ceil(float64(p.Population / 100000)))
}
