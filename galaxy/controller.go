package galaxy

import (
  "net/http"
  "encoding/json"
  "github.com/gorilla/context"
  "kalaxia-game-api/player"
  "strconv"
)

func GetMapAction(w http.ResponseWriter, r *http.Request) {
    currentPlayer := context.Get(r, "player").(*player.Player)
    gameMap := GetMapByServerId(currentPlayer.ServerId)

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(&gameMap); err != nil {
        panic(err)
    }
}

func GetSystemAction(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    id, _ := strconv.ParseUint(vars["id"], 10, 16)
    system := GetSystem(uint16(id))

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(&system); err != nil {
        panic(err)
    }
}

func GetPlanetAction(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    currentPlayer := context.Get(r, "player").(*player.Player)

    id, _ := strconv.ParseUint(vars["id"], 10, 16)
    planet := GetPlanet(uint16(id), currentPlayer.Id)

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(&planet); err != nil {
        panic(err)
    }
}

func GetPlayerPlanetsAction(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, _ := strconv.ParseUint(vars["id"], 10, 16)

    planets := GetPlayerPlanets(uint16(id))
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(&planets); err != nil {
        panic(err)
    }
}
