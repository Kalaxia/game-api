package manager

import(
  "kalaxia-game-api/database"
  "kalaxia-game-api/model/diplomacy"
)

func GetPlanetRelations(planetId uint16) []model.DiplomaticRelation {
  relations := make([]model.DiplomaticRelation, 0)
  if err := database.
    Connection.
    Model(&relations).
    Column("diplomatic_relation.*", "Faction", "Player").
    Where("diplomatic_relation.planet_id = ?", planetId).
    Select(); err != nil {
    return nil
  }
  return relations
}
