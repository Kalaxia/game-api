package controller

import (
    "net/http"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "kalaxia-game-api/model"
    "kalaxia-game-api/manager"
    "kalaxia-game-api/utils"
    "strconv"
)

func GetSystem(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)

    utils.SendJsonResponse(w, 200, manager.GetSystem(uint16(id)))
}

func GetSectorSystems(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*model.Player)
    starmap := manager.GetMapByServerId(player.ServerId)
    sectorId, _ := strconv.ParseUint(r.FormValue("sector"), 10, 16)

    utils.SendJsonResponse(w, 200, manager.GetSectorSystems(starmap, uint16(sectorId)))
}