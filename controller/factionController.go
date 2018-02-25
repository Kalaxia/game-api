package controller

import (
    "net/http"
    "encoding/json"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "kalaxia-game-api/manager"
    "kalaxia-game-api/model"
    "strconv"
)

func GetFactions(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*model.Player)
    factions := manager.GetServerFactions(player.ServerId)

    w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(&factions); err != nil {
        panic(err)
    }
}

func GetFactionPlanetChoices(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    id, _ := strconv.ParseUint(vars["id"], 10, 16)
    choices := manager.GetFactionPlanetChoices(uint16(id))

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(&choices); err != nil {
        panic(err)
    }
}
