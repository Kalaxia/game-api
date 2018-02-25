package faction

import (
  "net/http"
  "encoding/json"
  "github.com/gorilla/context"
  "github.com/gorilla/mux"
  "kalaxia-game-api/player"
  "strconv"
)

func GetFactionsAction(w http.ResponseWriter, r *http.Request) {
    currentPlayer := context.Get(r, "player").(*player.Player)
    factions := GetServerFactions(currentPlayer.ServerId)

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(&factions); err != nil {
        panic(err)
    }
}

func GetFactionPlanetChoicesAction(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    id, _ := strconv.ParseUint(vars["id"], 10, 16)
    choices := GetFactionPlanetChoices(uint16(id))

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(&choices); err != nil {
        panic(err)
    }
}
