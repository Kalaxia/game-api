package api

const(
    RelationAlly = "ally"
    RelationNeutral = "neutral"
    RelationHostile = "hostile"
)

type(
  DiplomaticRelation struct {
    tableName struct{} `json:"-" pg:"diplomacy__relations"`

    Planet *Planet `json:"planet"`
    PlanetId uint16 `json:"-"`
    Faction *Faction `json:"faction"`
    FactionId uint16 `json:"-"`
    Player *Player `json:"player"`
    PlayerId uint16 `json:"-"`
    Score int `json:"score"`
  }
  FactionRelation struct {
      tableName struct{} `json:"-" pg:"diplomacy__factions"`

      Faction *Faction `json:"-"`
      FactionId uint16 `json:"-"`
      OtherFaction *Faction `json:"faction"`
      OtherFactionId uint16 `json:"-"`

      State string `json:"state"`
  }
)

func (p *Planet) getPlanetRelations() []DiplomaticRelation {
  relations := make([]DiplomaticRelation, 0)
  if err := Database.
    Model(&relations).
    Relation("Faction").
    Relation("Player.Faction").
    Where("diplomatic_relation.planet_id = ?", p.Id).
    Select(); err != nil {
    panic(NewException("Planet not found", err))
  }
  return relations
}

func (p *Planet) increasePlayerRelation(player *Player, score int) {
    relation := &DiplomaticRelation{}
    if err := Database.
        Model(relation).
        Where("planet_id = ?", p.Id).
        Where("player_id = ?", player.Id).
        Select(); err != nil {
            p.createPlayerRelation(player, score)
            return
    }
    relation.Score += score
    if _, err := Database.
        Model(relation).
        Set("score = ?score").
        Where("planet_id = ?", p.Id).
        Where("player_id = ?", player.Id).
        Update(); err != nil {
            panic(NewException("Planet relation could not be updated", err))
    }
}

func (p *Planet) createPlayerRelation(player *Player, score int) {
    relation := &DiplomaticRelation{
        Planet: p,
        PlanetId: p.Id,
        Player: player,
        PlayerId: player.Id,
        Score: score,
    }
    if err := Database.Insert(relation); err != nil {
        panic(NewException("Planet relation could not be created", err))
    }
}
