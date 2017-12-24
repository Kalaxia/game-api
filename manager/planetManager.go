package manager

import(
  "kalaxia-game-api/database"
  "kalaxia-game-api/model/map"
)

func GetSystemPlanets(id uint16) []model.Planet {
  var planets []model.Planet
  if err := database.Connection.Model(&planets).Where("system_id = ?", id).Select(); err != nil {
    panic(err)
  }
  return planets
}
