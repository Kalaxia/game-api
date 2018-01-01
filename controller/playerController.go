package controller

import (
  "io"
  "io/ioutil"
  "net/http"
  "encoding/json"
  "github.com/gorilla/context"
  "kalaxia-game-api/manager"
  "kalaxia-game-api/model/player"
  "strconv"
)

func GetCurrentPlayer(w http.ResponseWriter, r *http.Request) {
  player := context.Get(r, "player")
  w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(player); err != nil {
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
  factionId, _ := strconv.ParseUint(data["faction_id"], 10, 16)
  planetId, _ := strconv.ParseUint(data["planet_id"], 10, 16)
  manager.RegisterPlayer(
    context.Get(r, "player").(*model.Player),
    uint16(factionId),
    uint16(planetId),
  )
  w.Write([]byte(""))
}
