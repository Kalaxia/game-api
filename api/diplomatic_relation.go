package api

import "fmt"

const(
    RelationAlly = "ally"
    RelationNeutral = "neutral"
    RelationHostile = "hostile"
)

type(
  DiplomaticRelation struct {
    tableName struct{} `pg:"diplomacy__relations"`

    Planet *Planet `json:"planet"`
    PlanetId uint16 `json:"-"`
    Faction *Faction `json:"faction"`
    FactionId uint16 `json:"-"`
    Player *Player `json:"player"`
    PlayerId uint16 `json:"-"`
    Score int `json:"score"`
  }

  FactionRelation struct {
      tableName struct{} `pg:"diplomacy__factions"`

      Faction *Faction `json:"-"`
      FactionId uint16 `json:"-" pg:",pk"`
      OtherFaction *Faction `json:"faction"`
      OtherFactionId uint16 `json:"-" pg:",pk"`

      State string `json:"state"`
      PurchaseTradeTax uint8 `json:"purchase_trade_tax"`
      SalesTradeTax uint8 `json:"sales_trade_tax"`
  }
)

func (p *Planet) getPlanetRelations() []*DiplomaticRelation {
  relations := make([]*DiplomaticRelation, 0)
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

func (p *Planet) createFactionRelation(faction *Faction, score int) *DiplomaticRelation {
    relation := &DiplomaticRelation{
        Planet: p,
        PlanetId: p.Id,
        Faction: faction,
        FactionId: faction.Id,
        Score: score,
    }
    if err := Database.Insert(relation); err != nil {
        panic(NewException("Planet relation could not be created", err))
    }
    return relation
}

func (f *Faction) getFactionRelation(otherFaction *Faction) *FactionRelation {
    relation := &FactionRelation{}
    if err := Database.Model(relation).Where("faction_id = ?", f.Id).Where("other_faction_id = ?", otherFaction.Id).Select(); err != nil {
        panic(NewException("Faction relation could not be retrieved", err))
    }
    return relation
}

func (f *Faction) validateWarDeclaration(data map[string]interface{}) {
    target := getFaction(uint16(data["id"].(float64)))

    if f.isInWarWith(target) {
        panic(NewHttpException(400, "faction.motions.types.war_declaration.already_in_war", nil))
    }
}

func (f *Faction) validatePeaceTreatySending(data map[string]interface{}) {
    target := getFaction(uint16(data["id"].(float64)))

    if !f.isInWarWith(target) {
        panic(NewHttpException(400, "faction.motions.types.peace_treaty_sending.not_in_war", nil))
    }
}

func (f *Faction) declareWar(data map[string]interface{}) {
    target := getFaction(uint16(data["id"].(float64)))

    relationship := f.getRelationWith(target)
    // If the war has already been declared by the other faction
    if relationship.State == RelationHostile {
        return
    }
    relationship.State = RelationHostile
    targetRelationship := target.getRelationWith(f)
    targetRelationship.State = RelationHostile

    f.createWar(target)

    relationship.update()
    targetRelationship.update()
}

func (f *Faction) sendPeaceTreaty(author *Player, data map[string]interface{}) {
    target := getFaction(uint16(data["id"].(float64)))
    relationship := f.getRelationWith(target)
    // If the peace has already been settled by the other faction
    if relationship.State != RelationHostile {
        return
    }
    target.createMotion(author, MotionTypePeaceTreatyResponse, map[string]interface{}{
        "faction": f,
    })
}

func (f *Faction) acceptPeaceTreaty(data map[string]interface{}) {
    target := getFaction(uint16(data["id"].(float64)))
    relationship := f.getRelationWith(target)
    // If the peace has already been settled by the other faction
    if relationship.State != RelationHostile {
        return
    }
    relationship.State = RelationNeutral
    targetRelationship := target.getRelationWith(f)
    targetRelationship.State = RelationNeutral

    f.endWar(target)

    relationship.update()
    targetRelationship.update()
}

func (f *Faction) getRelationWith(target *Faction) *FactionRelation {
    for _, r := range f.Relations {
        if r.OtherFactionId == target.Id {
            return r
        }
    }
    panic(NewException(fmt.Sprintf("Faction %d relation not found with faction %d", f.Id, target.Id), nil))
}

func (f *Faction) isInWarWith(target *Faction) bool {
    return f.getRelationWith(target).State == RelationHostile
}

func (fr *FactionRelation) update() {
    if err := Database.Update(fr); err != nil {
        panic(NewException("Could not update faction relation", err))
    }
}