package shipController

import (
    "net/http"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "kalaxia-game-api/manager/ship"
    "kalaxia-game-api/model"
    "kalaxia-game-api/utils"
    "strconv"
)

func GetPlayerShipModels(w http.ResponseWriter, r *http.Request) {
    utils.SendJsonResponse(w, 200, shipManager.GetShipPlayerModels(
        context.Get(r, "player").(*model.Player).Id,
    ))
}

func GetShipModel(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)

    utils.SendJsonResponse(w, 200, shipManager.GetShipModel(
        context.Get(r, "player").(*model.Player).Id,
        uint(id),
    ))
}

func CreateShipModel(w http.ResponseWriter, r *http.Request) {
    utils.SendJsonResponse(w, 201, shipManager.CreateShipModel(
        context.Get(r, "player").(*model.Player),
        utils.DecodeJsonRequest(r),
    ))
}
