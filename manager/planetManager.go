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

func GetPlayerPlanets(id uint16) []model.Planet {
  var planets []model.Planet
  if err := database.
      Connection.
      Model(&planets).
      Column("planet.*", "Resources").
      Where("planet.player_id = ?", id).
      Select(); err != nil {
    panic(err)
  }
  return planets
}

func GetPlanet(id uint16, playerId uint16) *model.Planet {
  var planet model.Planet
  if err := database.
    Connection.
    Model(&planet).
    Column("planet.*", "Player", "Resources", "System").
    Where("planet.id = ?", id).
    Select(); err != nil {
    return nil
  }
  relations := GetPlanetRelations(planet.Id)
  r := make([]interface{}, len(relations))
  for i, v := range relations {
      r[i] = v
  }
  planet.Relations = r
  if planet.Player != nil && playerId == planet.Player.Id {
    getPlanetOwnerData(&planet)
  }
  return &planet
}

func getPlanetOwnerData(planet *model.Planet) {
  buildings, availableBuildings := GetPlanetBuildings(planet.Id)
  r := make([]interface{}, len(buildings))
  for i, v := range buildings {
      r[i] = v
  }
  planet.Buildings = r
  r = make([]interface{}, len(availableBuildings))
  for i, v := range availableBuildings {
      r[i] = v
  }
  planet.AvailableBuildings = r
  planet.NbBuildings = 3
}
