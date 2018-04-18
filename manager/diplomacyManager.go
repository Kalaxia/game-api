package manager

import(
  "kalaxia-game-api/database"
  "kalaxia-game-api/exception"
  "kalaxia-game-api/model"
)

func GetPlanetRelations(planetId uint16) []model.DiplomaticRelation {
  relations := make([]model.DiplomaticRelation, 0)
  if err := database.
    Connection.
    Model(&relations).
    Column("diplomatic_relation.*", "Faction", "Player", "Player.Faction").
    Where("diplomatic_relation.planet_id = ?", planetId).
    Select(); err != nil {
    panic(exception.NewException("Planet not found", err))
  }
  return relations
}

func IncreasePlayerRelation(planet *model.Planet, player *model.Player, score int) {
    var relation model.DiplomaticRelation
    if err := database.
        Connection.
        Model(&relation).
        Where("planet_id = ?", planet.Id).
        Where("player_id = ?", player.Id).
        Select(); err != nil {
            createPlayerRelation(planet, player, score)
            return
    }
    relation.Score += score
    if _, err := database.
        Connection.
        Model(&relation).
        Set("score = ?score").
        Where("planet_id = ?", planet.Id).
        Where("player_id = ?", player.Id).
        Update(); err != nil {
            panic(exception.NewException("Planet relation could not be updated", err))
    }
}

func createPlayerRelation(planet *model.Planet, player *model.Player, score int) {
    relation := &model.DiplomaticRelation{
        Planet: planet,
        PlanetId: planet.Id,
        Player: player,
        PlayerId: player.Id,
        Score: score,
    }
    if err := database.Connection.Insert(relation); err != nil {
        panic(exception.NewException("Planet relation could not be created", err))
    }
}
