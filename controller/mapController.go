package controller

import (
  "net/http"
  "encoding/json"
  "github.com/gorilla/context"
  "kalaxia-game-api/manager"
  "kalaxia-game-api/model/player"
)

func GetMap(w http.ResponseWriter, r *http.Request) {
  player := context.Get(r, "player").(*model.Player)
  gameMap := manager.GetMapByServerId(player.ServerId)

  w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&gameMap); err != nil {
    panic(err)
  }
}
