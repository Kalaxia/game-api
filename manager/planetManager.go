package manager

import(
  "kalaxia-game-api/database"
  "kalaxia-game-api/model/map"
)

func GetSystemPlanets(id uint16) []model.Planet {
  var planets []model.Planet
  if err := database.
      Connection.
      Model(&planets).
      Column("planet.*", "Orbit").
      Where("planet.system_id = ?", id).
      Select(); err != nil {
    panic(err)
  }
  return planets
}

func GetPlanet(id uint16) *model.Planet {
  planet := model.Planet{Id: id}
  if err := database.Connection.Select(&planet); err != nil {
    return nil
  }
  return &planet
}
