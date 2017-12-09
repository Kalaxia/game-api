package controller

import (
  "net/http"
  "encoding/json"
  "github.com/gorilla/context"
)

func GetCurrentPlayer(w http.ResponseWriter, r *http.Request) {
  player := context.Get(r, "player")
  w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(player); err != nil {
    panic(err)
  }
}
