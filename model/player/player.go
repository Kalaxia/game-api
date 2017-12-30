package model

import(
  "time"
  "kalaxia-game-api/model/server"
)

type(
  Player struct {
    Id uint16 `json:"id"`
    Username string `json:"-" sql:"type:varchar(180);not null;unique"`
    Pseudo string `json:"pseudo" sql:"type:varchar(180);not null;unique"`
    ServerId uint16 `json:"-"`
    Server *model.Server `json:"-"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
  }
  Players []Player
)
