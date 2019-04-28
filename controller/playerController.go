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
    currentPlayer := context.Get(r, "player").(*model.Player)

    utils.SendJsonResponse(w, 200, manager.GetPlayer(uint16(id), currentPlayer.Id == uint16(id)))
}

func GetPlayerPlanets(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)

    utils.SendJsonResponse(w, 200, manager.GetPlayerPlanets(uint16(id)))
}

func UpdateCurrentPlanet(w http.ResponseWriter, r *http.Request) {
    data := utils.DecodeJsonRequest(r)
    player := context.Get(r, "player").(*model.Player)

    manager.UpdateCurrentPlanet(player, uint16(data["planet_id"].(float64)))
    
    w.WriteHeader(204)
    w.Write([]byte(""))
}

func RegisterPlayer(w http.ResponseWriter, r *http.Request) {
    data := utils.DecodeJsonRequest(r)
    player := context.Get(r, "player").(*model.Player)
    if player.IsActive == true {
        panic(exception.NewHttpException(http.StatusForbidden, "Player account is already active", nil))
    }
    manager.RegisterPlayer(
        player,
        data["pseudo"].(string),
        data["gender"].(string),
        data["avatar"].(string),
        uint16(data["faction_id"].(float64)),
        uint16(data["planet_id"].(float64)),
    )
    w.WriteHeader(204)
    w.Write([]byte(""))
}
