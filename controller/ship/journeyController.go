package shipController


import (
    "net/http"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "kalaxia-game-api/exception"
	"kalaxia-game-api/manager/ship"
    "kalaxia-game-api/model"
	"kalaxia-game-api/utils"
    "strconv"
)

func GetJourney (w http.ResponseWriter, r *http.Request){
	
	player := context.Get(r, "player").(*model.Player)
	idJourney, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16);
	fleet := shipManager.GetFleetOnJourney (uint16(idJourney));
	
	
	if player.Id != fleet.Player.Id { // the player does not own the planet
		panic(exception.NewHttpException(http.StatusForbidden, "", nil));
	}
	if !fleet.IsOnJourney() {
		panic(exception.NewHttpException(400, "This journey has ended", nil));
	}
	
	utils.SendJsonResponse(w, 200,fleet.Journey);
}

func CreateJourney (w http.ResponseWriter, r *http.Request) {
	
	
	utils.SendJsonResponse(w, 202,"");
}
