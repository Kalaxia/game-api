package model

import(
  mapModel "kalaxia-game-api/model/map"
  factionModel "kalaxia-game-api/model/faction"
  playerModel "kalaxia-game-api/model/player"
)

type(
  DiplomaticRelation struct {
    TableName struct{} `json:"-" sql:"diplomacy__relations"`

    Planet *mapModel.Planet `json:"planet"`
    PlanetId uint16 `json:"-"`
    Faction *factionModel.Faction `json:"faction"`
    FactionId uint16 `json:"-"`
    Player *playerModel.Player `json:"player"`
    PlayerId uint16 `json:"-"`
    Score int `json:"score"`
  }
)
