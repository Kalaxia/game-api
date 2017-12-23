package model

type(
  System struct {
    TableName struct{} `json:"-" sql:"map__systems"`

    Id uint16 `json:"id"`
    MapId uint16 `json:"-"`
    Map *Map `json:"-"`
    Planets Planets `json:"planets"`
    X uint16 `json:"coord_x"`
    Y uint16 `json:"coord_y"`
    Orbits SystemOrbits `json:"orbits"`
  }
  SystemOrbit struct {
    TableName struct{} `json:"-" sql:"map__system_orbits"`

    Id uint16 `json:"id"`
    Radius uint16 `json:"radius"`
    SystemId uint16 `json:"-"`
    System *System `json:"system"`
  }
  SystemOrbits []*SystemOrbit
)
