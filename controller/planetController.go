package controller

import (
    "net/http"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/manager"
    "kalaxia-game-api/model"
    "kalaxia-game-api/utils"
    "strconv"
)

func GetPlanet(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*model.Player)
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)

    utils.SendJsonResponse(w, 200, manager.GetPlayerPlanet(uint16(id), player.Id))
}

func UpdatePlanetSettings(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*model.Player)

    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := manager.GetPlayerPlanet(uint16(id), player.Id)

    if player.Id != planet.Player.Id {
        panic(exception.NewHttpException(http.StatusForbidden, "", nil))
    }
    data := utils.DecodeJsonRequest(r)
    settings := &model.PlanetSettings{
        ServicesPoints: uint8(data["services_points"].(float64)),
        BuildingPoints: uint8(data["building_points"].(float64)),
        MilitaryPoints: uint8(data["military_points"].(float64)),
        ResearchPoints: uint8(data["research_points"].(float64)),
    }
    manager.UpdatePlanetSettings(planet, settings)
    utils.SendJsonResponse(w, 200, settings)
}
