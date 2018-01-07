package model

import(
  playerModel "kalaxia-game-api/model/player"
)

const PLANET_TYPE_ARCTIC = "arctic"
const PLANET_TYPE_DESERT = "desert"
const PLANET_TYPE_OCEANIC = "oceanic"
const PLANET_TYPE_TEMPERATE = "temperate"
const PLANET_TYPE_TROPICAL = "tropical"
const PLANET_TYPE_ROCKY = "rocky"
const PLANET_TYPE_VOLCANIC = "volcanic"

type(
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
    Player *playerModel.Player `json:"player"`
    Resources []PlanetResource `json:"resources"`
    Relations []interface{} `json:"relations" sql:"-"`
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
)
