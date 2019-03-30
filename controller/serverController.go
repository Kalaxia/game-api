package controller

import(
    "net/http"
    "github.com/gorilla/mux"
    "kalaxia-game-api/manager"
    "kalaxia-game-api/utils"
    "strconv"
)

func CreateServer(w http.ResponseWriter, r *http.Request) {
    data := utils.DecodeJsonRequest(r)
    server := manager.CreateServer(
        data["name"].(string),
        data["type"].(string),
        data["signature"].(string),
    )
    factions := manager.CreateServerFactions(server, data["factions"].([]interface{}))
    manager.GenerateMap(server, factions, uint16(data["map_size"].(float64)))
    utils.SendJsonResponse(w, 201, server)
}

func RemoveServer(w http.ResponseWriter, r *http.Request) {
	serverId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    data := utils.DecodeJsonRequest(r)
    manager.RemoveServer(uint16(serverId), data["signature"].(string))
}