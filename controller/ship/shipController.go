package shipController

import (
    "net/http"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/manager/ship"
    "kalaxia-game-api/manager"
    "kalaxia-game-api/model"
    "kalaxia-game-api/utils"
    "strconv"
)

func CreateShip(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*model.Player)
    planetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := manager.GetPlayerPlanet(uint16(planetId), player.Id)

    if planet.Player.Id != player.Id {
        panic(exception.NewHttpException(403, "You do not control this planet", nil))
    }
    utils.SendJsonResponse(w, 201, shipManager.CreateShip(
        player,
        planet,
        utils.DecodeJsonRequest(r)),
    )
}

func GetCurrentlyConstructingShips(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*model.Player)
    planetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := manager.GetPlayerPlanet(uint16(planetId), player.Id)

    if planet.Player.Id != player.Id {
        panic(exception.NewHttpException(403, "You do not control this planet", nil))
    }
    utils.SendJsonResponse(w, 200, shipManager.GetCurrentlyConstructingShips(planet))
}

func GetConstructingShips(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*model.Player)
    planetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := manager.GetPlayerPlanet(uint16(planetId), player.Id)

    if planet.Player.Id != player.Id {
        panic(exception.NewHttpException(403, "You do not control this planet", nil))
    }
    utils.SendJsonResponse(w, 200, shipManager.GetConstructingShips(planet))
}

func GetHangarShips(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*model.Player)
    planetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := manager.GetPlayerPlanet(uint16(planetId), player.Id)

    if planet.Player.Id != player.Id {
        panic(exception.NewHttpException(403, "You do not control this planet", nil))
    }
    utils.SendJsonResponse(w, 200, shipManager.GetHangarShips(planet))
}

func GetHangarShipGroups(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*model.Player)
    planetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := manager.GetPlayerPlanet(uint16(planetId), player.Id)

    if planet.Player.Id != player.Id {
        panic(exception.NewHttpException(403, "You do not control this planet", nil))
    }
    utils.SendJsonResponse(w, 200, shipManager.GetHangarShipGroups(planet))
}
