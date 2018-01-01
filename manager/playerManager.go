package manager

import(
  "time"
  "errors"
  "kalaxia-game-api/database"
  playerModel "kalaxia-game-api/model/player"
  serverModel "kalaxia-game-api/model/server"
)

func GetPlayer(id uint16) *playerModel.Player {
  var player playerModel.Player
  if err := database.Connection.Model(&player).Column("player.*", "Faction").Where("player.id = ?", id).Select(); err != nil {
    return nil
  }
  return &player
}

func GetPlayerByUsername(username string, server *serverModel.Server) *playerModel.Player {
  player := playerModel.Player{Username: username}
  err := database.
    Connection.
    Model(&player).
    Column("player.*", "Server").
    Where("username = ?", username).
    Where("server_id = ?", server.Id).
    Select()
  if err != nil {
    return nil
  }
  return &player
}

func CreatePlayer(username string, server *serverModel.Server) *playerModel.Player {
  player := playerModel.Player{
    Username: username,
    Pseudo: username,
    ServerId: server.Id,
    Server: server,
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
  }
  if err := database.Connection.Insert(&player); err != nil {
    panic(err)
  }
  return &player
}

func RegisterPlayer(player *playerModel.Player, factionId uint16, planetId uint16) {
  faction := GetFaction(factionId)
  if faction == nil {
    panic(errors.New("faction not found"))
  }
  planet := GetPlanet(planetId)
  if planet == nil {
    panic(errors.New("planet not found"))
  }
  planet.PlayerId = player.Id
  planet.Player = player
  player.FactionId = faction.Id
  player.Faction = faction
  player.IsActive = true
  if err := database.Connection.Update(player); err != nil {
    panic(err)
  }
  if err := database.Connection.Update(planet); err != nil {
    panic(err)
  }
}
