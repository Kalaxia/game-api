package controller

import (
    "net/http"
    "github.com/gorilla/context"
    "kalaxia-game-api/manager"
    "kalaxia-game-api/model"
    "kalaxia-game-api/utils"
)

func GetMap(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*model.Player)
    utils.SendJsonResponse(w, 200, manager.GetMapByServerId(player.ServerId))
}
