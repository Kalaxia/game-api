package shipController

import(
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
	"kalaxia-game-api/exception"
	"kalaxia-game-api/model"
	"kalaxia-game-api/manager/ship"
	"kalaxia-game-api/utils"
	"net/http"
	"strconv"
)

func GetCombatReport(w http.ResponseWriter, r *http.Request) {
	player := context.Get(r, "player").(*model.Player)
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)

	report := shipManager.GetCombatReport(uint16(id))

	if report.Attacker.Player.Id != player.Id && report.Defender.Player.Id != player.Id {
		panic(exception.NewHttpException(403, "You do not own this combat report", nil))
	}
	utils.SendJsonResponse(w, 200, report)
}

func GetCombatReports(w http.ResponseWriter, r *http.Request) {
	player := context.Get(r, "player").(*model.Player)

	utils.SendJsonResponse(w, 200, shipManager.GetCombatReports(player))
}