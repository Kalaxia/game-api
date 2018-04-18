package controller

import(
    "net/http"
    "kalaxia-game-api/manager"
    "kalaxia-game-api/utils"
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
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte(""))
}
