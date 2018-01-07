package manager

import(
  "kalaxia-game-api/database"
  "kalaxia-game-api/model/diplomacy"
  mapModel "kalaxia-game-api/model/map"
  playerModel "kalaxia-game-api/model/player"
)

func GetPlanetRelations(planetId uint16) []model.DiplomaticRelation {
  relations := make([]model.DiplomaticRelation, 0)
  if err := database.
    Connection.
    Model(&relations).
    Column("diplomatic_relation.*", "Faction", "Player", "Player.Faction").
    Where("diplomatic_relation.planet_id = ?", planetId).
    Select(); err != nil {
    return nil
  }
  return relations
}

func IncreasePlayerRelation(planet *mapModel.Planet, player *playerModel.Player, score int) {
  var relation *model.DiplomaticRelation
  if err := database.
    Connection.
    Model(relation).
    Where("diplomatic_relation.planet_id = ?", planet.Id).
    Where("diplomatic_relation.player_id = ?", player.Id).
    Select(); err != nil {
    createPlayerRelation(planet, player, score)
    return
  }
  relation.Score += score
  if err := database.Connection.Update(relation); err != nil {
    panic(err)
  }
}

func createPlayerRelation(planet *mapModel.Planet, player *playerModel.Player, score int) {
  relation := &model.DiplomaticRelation{
    Planet: planet,
    PlanetId: planet.Id,
    Player: player,
    PlayerId: player.Id,
    Score: score,
  }
  if err := database.Connection.Insert(relation); err != nil {
    panic(err)
  }
}
