package api

import(
    "github.com/gorilla/context"
    "net/http"
)

type(
	Map struct {
	  TableName struct{} `json:"-" sql:"map__maps"`
  
	  Id uint16 `json:"-"`
	  ServerId uint16 `json:"-"`
	  Server *Server `json:"-"`
      Systems []System `json:"systems" sql:"-"`
      Territories []Territory `json:"territories"`
	  Size uint16 `json:"size"`
	  SectorSize uint16 `json:"sector_size" sql:"-"`
	}
)  

func GetMap(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)
    SendJsonResponse(w, 200, getServerMap(player.ServerId))
}

func (s *Server) generateMap(factions []*Faction, size uint16) *Map {
    gameMap := &Map{
        Server: s,
        ServerId: s.Id,
        Size: size,
    }
    if err := Database.Insert(gameMap); err != nil {
        panic(NewException("Map could not be created", err))
    }
    gameMap.generateSystems(factions)
    return gameMap
}

func getServerMap(serverId uint16) *Map {
    gameMap := &Map{}
    if err := Database.Model(gameMap).Where("server_id = ?", serverId).Select(); err != nil {
        return nil
    }
    gameMap.Systems = gameMap.getSystems()
    gameMap.SectorSize = 10
    return gameMap
}
