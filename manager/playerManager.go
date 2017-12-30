package manager

import(
  "time"
  "kalaxia-game-api/database"
  playerModel "kalaxia-game-api/model/player"
  serverModel "kalaxia-game-api/model/server"
)

func GetPlayer(id uint16) *playerModel.Player {
  player := playerModel.Player{Id: id}
  if err := database.Connection.Select(&player); err != nil {
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
