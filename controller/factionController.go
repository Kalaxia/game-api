package controller

import (
  "net/http"
  "encoding/json"
  "github.com/gorilla/context"
  "kalaxia-game-api/manager"
  "kalaxia-game-api/model/player"
)

func GetFactions(w http.ResponseWriter, r *http.Request) {
  player := context.Get(r, "player").(*model.Player)
  factions := manager.GetServerFactions(player.ServerId)

  w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&factions); err != nil {
    panic(err)
  }
}
