package model

const PLANET_TYPE_ARCTIC = "arctic"
const PLANET_TYPE_DESERT = "desert"
const PLANET_TYPE_OCEANIC = "oceanic"
const PLANET_TYPE_TEMPERATE = "temperate"
const PLANET_TYPE_TROPICAL = "tropical"
const PLANET_TYPE_ROCKY = "rocky"
const PLANET_TYPE_VOLCANIC = "volcanic"

type(
  Planet struct {
    TableName struct{} `sql:"map__planets"`

    Id uint16 `json:"id"`
    Name string `json:"name"`
    Type string `json:"type"`
    SystemId uint16 `json:"-"`
    System *System `json:"system"`
    OrbitId uint16 `json:"-"`
    Orbit *SystemOrbit `json:"orbit"`
  }
  Planets []*Planet
)
