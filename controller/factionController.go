package controller

import (
    "net/http"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "kalaxia-game-api/manager"
    "kalaxia-game-api/model"
    "kalaxia-game-api/utils"
    "strconv"
)

func GetFactions(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*model.Player)

    utils.SendJsonResponse(w, 200, manager.GetServerFactions(player.ServerId))
}

func GetFaction(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)

    utils.SendJsonResponse(w, 200, manager.GetFaction(uint16(id)))
}

func GetFactionPlanetChoices(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)

    utils.SendJsonResponse(w, 200, manager.GetFactionPlanetChoices(uint16(id)))
}

func GetFactionMembers(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)

    utils.SendJsonResponse(w, 200, manager.GetFactionMembers(uint16(id)))
}
