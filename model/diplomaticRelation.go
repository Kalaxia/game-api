package model

const(
    RELATION_ALLY = "ally"
    RELATION_NEUTRAL = "neutral"
    RELATION_HOSTILE = "hostile"
)

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
  FactionRelation struct {
      TableName struct{} `json:"-" sql:"diplomacy__factions"`

      Faction *Faction `json:"-"`
      FactionId uint16 `json:"-"`
      OtherFaction *Faction `json:"faction"`
      OtherFactionId uint16 `json:"-"`

      State string `json:"state"`
  }
)
