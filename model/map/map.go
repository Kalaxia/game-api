package model

import(
  "kalaxia-game-api/model/server"
)

type(
  Map struct {
    TableName struct{} `sql:"map__maps"`

    Id uint16
    ServerId int16
    Server *model.Server
    Size uint16
  }
)
