package galaxy

import(
    "kalaxia-game-api/player"
    "kalaxia-game-api/server"
)

const PLANET_TYPE_ARCTIC = "arctic"
const PLANET_TYPE_DESERT = "desert"
const PLANET_TYPE_OCEANIC = "oceanic"
const PLANET_TYPE_TEMPERATE = "temperate"
const PLANET_TYPE_TROPICAL = "tropical"
const PLANET_TYPE_ROCKY = "rocky"
const PLANET_TYPE_VOLCANIC = "volcanic"

type(
    Map struct {
        TableName struct{} `json:"-" sql:"map__maps"`

        Id uint16 `json:"-"`
        ServerId uint16 `json:"-"`
        Server *interface{} `json:"-"`
        Systems []System `json:"systems" sql:"-"`
        Size uint16 `json:"size"`
    }
    Planet struct {
        TableName struct{} `json:"-" sql:"map__planets"`

        Id uint16 `json:"id"`
        Name string `json:"name"`
        Type string `json:"type"`
        SystemId uint16 `json:"-"`
        System *System `json:"system"`
        OrbitId uint16 `json:"-"`
        Orbit *SystemOrbit `json:"orbit"`
        PlayerId uint16 `json:"-"`
        Player *interface{} `json:"player"`
        Resources []PlanetResource `json:"resources"`
        Relations []interface{} `json:"relations" sql:"-"`
        Buildings []interface{} `json:"buildings" sql:"-"`
        NbBuildings uint8 `json:"nb_buildings" sql:"-"`
        AvailableBuildings []interface{} `json:"available_buildings" sql:"-"`
    }
    PlanetResource struct {
        TableName struct{} `json:"-" sql:"map__planet_resources"`

        Name string `json:"name"`
        Density uint8 `json:"density"`
        PlanetId uint16 `json:"-"`
        Planet *Planet `json:"planet"`
    }
    PlanetData struct {
        Resources map[string]uint8
    }
    PlanetsData map[string]PlanetData

    System struct {
        TableName struct{} `json:"-" sql:"map__systems"`

        Id uint16 `json:"id"`
        MapId uint16 `json:"-"`
        Map *Map `json:"-"`
        Planets []Planet `json:"planets"`
        X uint16 `json:"coord_x"`
        Y uint16 `json:"coord_y"`
        Orbits []SystemOrbit `json:"orbits"`
    }
    SystemOrbit struct {
        TableName struct{} `json:"-" sql:"map__system_orbits"`

        Id uint16 `json:"id"`
        Radius uint16 `json:"radius"`
        SystemId uint16 `json:"-"`
        System *System `json:"system"`
    }
)
