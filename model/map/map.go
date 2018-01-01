package model

import(
  "kalaxia-game-api/model/server"
)

type(
  Map struct {
    TableName struct{} `json:"-" sql:"map__maps"`

    Id uint16 `json:"-"`
    ServerId uint16 `json:"-"`
    Server *model.Server `json:"-"`
    Systems []System `json:"systems" sql:"-"`
    Size uint16 `json:"size"`
  }
)
