package controller

import (
    "net/http"
    "encoding/json"
    "io"
    "io/ioutil"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "kalaxia-game-api/manager"
    "kalaxia-game-api/model"
    "strconv"
)

func GetPlanet(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    player := context.Get(r, "player").(*model.Player)

    id, _ := strconv.ParseUint(vars["id"], 10, 16)
    planet := manager.GetPlanet(uint16(id), player.Id)

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(&planet); err != nil {
        panic(err)
    }
}

func UpdatePlanetSettings(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    player := context.Get(r, "player").(*model.Player)

    id, _ := strconv.ParseUint(vars["id"], 10, 16)
    planet := manager.GetPlanet(uint16(id), player.Id)

    if player.Id != planet.Player.Id {
        w.WriteHeader(http.StatusForbidden)
        return
    }

    var body []byte
    var err error
    if body, err = ioutil.ReadAll(io.LimitReader(r.Body, 1048576)); err != nil {
        panic(err)
    }
    if err = r.Body.Close(); err != nil {
        panic(err)
    }
    var settings model.PlanetSettings
    if err = json.Unmarshal(body, &settings); err != nil {
        panic(err)
    }
    if err = manager.UpdatePlanetSettings(planet, &settings); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(&settings); err != nil {
        panic(err)
    }
}
