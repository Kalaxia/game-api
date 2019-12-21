package api

import(
    "github.com/gorilla/context"
    "net/http"
)

type(
	Map struct {
	  tableName struct{} `pg:"map__maps"`
  
	  Id uint16 `json:"-"`
	  ServerId uint16 `json:"-"`
	  Server *Server `json:"-"`
      Systems []System `json:"systems" pg:"-"`
      Territories []Territory `json:"territories"`
	  Size uint16 `json:"size"`
	  SectorSize uint16 `json:"sector_size" pg:"-"`
    }

    Place struct {
        tableName struct{} `pg:"map__places"`

        Id uint32 `json:"id"`
		PlanetId uint16 `json:"-"`
		Planet *Planet `json:"planet"`
		Coordinates *Coordinates `json:"coordinates"`
    }
    
    Coordinates struct {
        X float64 `json:"x"`
        Y float64 `json:"y"`
    }
    CoordinatesSlice []*Coordinates
)  

func GetMap(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)
    SendJsonResponse(w, 200, getServerMap(player.ServerId))
}

func getAllMaps() []*Map {
    maps := make([]*Map, 0)
    if err := Database.Model(&maps).Select(); err != nil {
        panic(NewException("Maps could not be retrieved", err))
    }
    return maps
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

func NewPlace(p *Planet, x, y float64) *Place {
    place := &Place{
        PlanetId: p.Id,
        Planet: p,
        Coordinates: &Coordinates{
            X: x,
            Y: y,
        },
    }
    place.create()
    return place
}

func NewCoordinatesPlace(x, y float64) *Place {
    place := &Place{
        Coordinates: &Coordinates{
            X: x,
            Y: y,
        },
    }
    place.create()
    return place
}

func (p *Place) create() {
    if err := Database.Insert(p); err != nil {
        panic(NewException("Could not create place", err))
    }
}

func (cs CoordinatesSlice) Len() int {
	return len(cs)
}

func (cs CoordinatesSlice) Swap(i, j int) {
	cs[i], cs[j] = cs[j], cs[i]
}

func (cs CoordinatesSlice) Less(i, j int) bool {
	if cs[i].X != cs[j].X {
		return cs[i].X < cs[j].X
	}
	return cs[i].Y < cs[j].Y
}