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

func CreateBuilding(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    player := context.Get(r, "player").(*model.Player)

    id, _ := strconv.ParseUint(vars["id"], 10, 16)
    planet := manager.GetPlanet(uint16(id), player.Id)

    if uint16(id) != planet.Id {
        panic(exception.NewHttpException(403, "Forbidden", nil))
    }
    data := utils.DecodeJsonRequest(r)
    utils.SendJsonResponse(w, 201, manager.CreateBuilding(planet, data["name"].(string)))
}

func CancelBuilding(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    player := context.Get(r, "player").(*model.Player)

    planetId, _ := strconv.ParseUint(vars["planet-id"], 10, 16)
    buildingId, _ := strconv.ParseUint(vars["building-id"], 10, 16)
    planet := manager.GetPlanet(uint16(planetId), player.Id)

    if uint16(planetId) != planet.Id {
        panic(exception.NewHttpException(403, "Forbidden", nil))
    }
    manager.CancelBuilding(planet, uint32(buildingId))

    w.WriteHeader(204)
    w.Write([]byte(""))
}