package api

import(
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "math"
    "net/http"
    "strconv"
)

type(
  System struct {
    tableName struct{} `pg:"map__systems"`

    Id uint16 `json:"id"`
    MapId uint16 `json:"-"`
    Map *Map `json:"-"`
    FactionId uint16 `json:"-"`
    Faction *Faction `json:"faction"`
    Territories []*SystemTerritory `json:"-" pg:"fk:system_id"`
    Planets []Planet `json:"planets"`
    X uint16 `json:"coord_x"`
    Y uint16 `json:"coord_y"`
    Orbits []SystemOrbit `json:"orbits"`
  }
  SystemOrbit struct {
    tableName struct{} `pg:"map__system_orbits"`

    Id uint16 `json:"id"`
    Radius uint16 `json:"radius"`
    SystemId uint16 `json:"-"`
    System *System `json:"system"`
  }
)

func GetSystem(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)

    SendJsonResponse(w, 200, getSystem(uint16(id)))
}

func GetSectorSystems(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)
    starmap := getServerMap(player.ServerId)
    sectorId, _ := strconv.ParseUint(r.FormValue("sector"), 10, 16)

    SendJsonResponse(w, 200, starmap.getSectorSystems(uint16(sectorId)))
}

func (m *Map) getSystems() []System {
    systems := make([]System, 0)
    if err := Database.Model(&systems).Where("map_id = ?", m.Id).Select(); err != nil {
        panic(NewHttpException(404, "Map not found", err))
    }
    return systems
}

func (m *Map) getSectorSystems(sector uint16) []System {
    sectorsPerLine := m.Size / m.SectorSize
    lineNumber := uint16(math.Ceil(float64((sector - 1) / sectorsPerLine)))

    systems := make([]System, 0)
    if err := Database.
        Model(&systems).
        Relation("Faction").
        Relation("Planets.Player.Faction").
        Where("map_id = ?", m.Id).
        Where("x >= ?", (sector - ((lineNumber * sectorsPerLine) + 1)) * m.SectorSize).
        Where("x <= ?", (sector - (lineNumber * sectorsPerLine)) * m.SectorSize).
        Where("y >= ?", lineNumber * sectorsPerLine).
        Where("y <= ?", (lineNumber + 1) * sectorsPerLine).
        Select(); err != nil {
        panic(NewHttpException(404, "Map not found", err))
    }
    return systems
}

func getSystem(id uint16) *System {
    system := System{Id: id}
    if err := Database.Select(&system); err != nil {
        return nil
    }
    system.Planets = system.getPlanets()
    return &system
}

func (s *System) update() {
    if err := Database.Update(s); err != nil {
        panic(NewException("Could not update system", err))
    }
}