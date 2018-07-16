package shipController

import (
    "net/http"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "kalaxia-game-api/manager/ship"
    "kalaxia-game-api/manager"
    "kalaxia-game-api/model"
    "kalaxia-game-api/utils"
    "strconv"
)

func CreateShip(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*model.Player)
    planetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)

    utils.SendJsonResponse(w, 200, shipManager.CreateShip(
        player,
        manager.GetPlanet(uint16(planetId), player.Id),
        utils.DecodeJsonRequest(r)),
    )
}
