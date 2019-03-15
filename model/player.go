package model

import(
  "time"
)

type(
  Player struct {
    Id uint16 `json:"id"`
    Username string `json:"-" sql:"type:varchar(180);not null;unique"`
    Pseudo string `json:"pseudo" sql:"type:varchar(180);not null;unique"`
    Gender string `json:"gender"`
    Avatar string `json:"avatar"`
    ServerId uint16 `json:"-"`
    Server *Server `json:"-"`
    FactionId uint16 `json:"-"`
    Faction *Faction `json:"faction"`
    IsActive bool `json:"is_active"`
    Wallet uint32 `json:"wallet"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
  }
  Players []Player
)
