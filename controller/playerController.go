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

func GetCurrentPlayer(w http.ResponseWriter, r *http.Request) {
    utils.SendJsonResponse(w, 200, context.Get(r, "player"))
}

func GetPlayer(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)

    utils.SendJsonResponse(w, 200, manager.GetPlayer(uint16(id)))
}

func GetPlayerPlanets(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)

    utils.SendJsonResponse(w, 200, manager.GetPlayerPlanets(uint16(id)))
}

func RegisterPlayer(w http.ResponseWriter, r *http.Request) {
    data := utils.DecodeJsonRequest(r)
    player := context.Get(r, "player").(*model.Player)
    if player.IsActive == true {
        panic(exception.NewHttpException(http.StatusForbidden, "", nil))
    }
    factionId, _ := strconv.ParseUint(data["faction_id"].(string), 10, 16)
    planetId, _ := strconv.ParseUint(data["planet_id"].(string), 10, 16)
    manager.RegisterPlayer(
        player,
        uint16(factionId),
        uint16(planetId),
    )
    w.WriteHeader(204)
    w.Write([]byte(""))
}
