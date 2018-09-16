package tradeController

import(
    "net/http"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"

    "kalaxia-game-api/exception"
    "kalaxia-game-api/manager"
    "kalaxia-game-api/manager/trade"
    "kalaxia-game-api/model"
    "kalaxia-game-api/utils"
    "strconv"
)

func CreateOffer(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*model.Player)
    planetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := manager.GetPlanet(uint16(planetId), player.Id)

    if planet.Player.Id != player.Id {
        panic(exception.NewHttpException(403, "You do not control this planet", nil))
    }
    utils.SendJsonResponse(w, 201, tradeManager.CreateOffer(planet, utils.DecodeJsonRequest(r)))
}

func SearchOffers(w http.ResponseWriter, r *http.Request) {
    utils.SendJsonResponse(w, 200, tradeManager.SearchOffers(utils.DecodeJsonRequest(r)))
}
