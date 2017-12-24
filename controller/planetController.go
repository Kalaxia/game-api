package controller

import (
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
  "kalaxia-game-api/manager"
  "strconv"
)

func GetPlanet(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  id, _ := strconv.ParseUint(vars["id"], 10, 16)
  planet := manager.GetPlanet(uint16(id))

  w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&planet); err != nil {
    panic(err)
  }
}
