package shipController

import (
    "net/http"
    "github.com/gorilla/context"
    "kalaxia-game-api/manager/ship"
    "kalaxia-game-api/model"
    "kalaxia-game-api/utils"
)

func GetPlayerShipModels(w http.ResponseWriter, r *http.Request) {
    utils.SendJsonResponse(w, 200, shipManager.GetShipPlayerModels(
        context.Get(r, "player").(*model.Player).Id,
    ))
}

func CreateShipModel(w http.ResponseWriter, r *http.Request) {
    utils.SendJsonResponse(w, 201, shipManager.CreateShipModel(
        context.Get(r, "player").(*model.Player),
        utils.DecodeJsonRequest(r),
    ))
}
