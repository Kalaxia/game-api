package model

type(
  Faction struct {
    TableName struct{} `json:"-" sql:"faction__factions"`

    Id uint16 `json:"id"`
    Name string `json:"name"`
    Description string `json:"description"`
    Color string `json:"color"`
    ServerId uint16 `json:"-"`
    Server *Server `json:"-"`
  }
)
