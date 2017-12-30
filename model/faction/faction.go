package model

import(
  serverModel "kalaxia-game-api/model/server"
)

type(
  Faction struct {
    TableName struct{} `json:"-" sql:"faction__factions"`

    Id uint16 `json:"id"`
    Name string `json:"name"`
    Description string `json:"description"`
    ServerId uint16 `json:"-"`
    Server *serverModel.Server `json:"-"`
  }
)
