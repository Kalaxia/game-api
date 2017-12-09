package model

import(
  "time"
  "kalaxia-game-api/model/server"
)

type(
  Player struct {
    Id int16 `json:"id"`
    Username string `json:"_" sql:"type:varchar(180);not null;unique"`
    Pseudo string `json:"pseudo" sql:"type:varchar(180);not null;unique"`
    ServerId int16 `json:"_"`
    Server *model.Server `json:"_"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
  }
  Players []Player
)
