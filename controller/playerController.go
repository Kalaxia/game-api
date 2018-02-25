package controller

import (
    "io"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "kalaxia-game-api/manager"
    "kalaxia-game-api/model"
    "strconv"
)

func GetCurrentPlayer(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player")
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(player); err != nil {
        panic(err)
    }
}

func GetPlayerPlanets(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, _ := strconv.ParseUint(vars["id"], 10, 16)

    planets := manager.GetPlayerPlanets(uint16(id))
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(&planets); err != nil {
        panic(err)
    }
}

func RegisterPlayer(w http.ResponseWriter, r *http.Request) {
    var body []byte
    var err error
    if body, err = ioutil.ReadAll(io.LimitReader(r.Body, 1048576)); err != nil {
        panic(err)
    }
    if err = r.Body.Close(); err != nil {
        panic(err)
    }
    var data map[string]string
    if err = json.Unmarshal(body, &data); err != nil {
        panic(err)
    }
    player := context.Get(r, "player").(*model.Player)
    if player.IsActive == true {
        w.WriteHeader(http.StatusForbidden)
        return
    }
    factionId, _ := strconv.ParseUint(data["faction_id"], 10, 16)
    planetId, _ := strconv.ParseUint(data["planet_id"], 10, 16)
    manager.RegisterPlayer(
        player,
        uint16(factionId),
        uint16(planetId),
    )
    w.Write([]byte(""))
}
