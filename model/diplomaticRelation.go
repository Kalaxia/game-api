package model

type(
  DiplomaticRelation struct {
    TableName struct{} `json:"-" sql:"diplomacy__relations"`

    Planet *Planet `json:"planet"`
    PlanetId uint16 `json:"-"`
    Faction *Faction `json:"faction"`
    FactionId uint16 `json:"-"`
    Player *Player `json:"player"`
    PlayerId uint16 `json:"-"`
    Score int `json:"score"`
  }
)
