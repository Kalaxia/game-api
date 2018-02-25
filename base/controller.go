package base

import (
  "io"
  "io/ioutil"
  "net/http"
  "encoding/json"
  "github.com/gorilla/context"
  "github.com/gorilla/mux"
  "kalaxia-game-api/galaxy"
  "kalaxia-game-api/player"
  "strconv"
)

func CreateBuildingAction(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    player := context.Get(r, "player").(*player.Player)

    id, _ := strconv.ParseUint(vars["id"], 10, 16)
    planet := galaxy.GetPlanet(uint16(id), player.Id)

    if uint16(id) != planet.Id {
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
    var data map[string]string
    if err = json.Unmarshal(body, &data); err != nil {
        panic(err)
    }
    building := CreateBuilding(planet, data["name"])

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(&building); err != nil {
        panic(err)
    }
}
