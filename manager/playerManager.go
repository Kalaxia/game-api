package manager

import(
  "fmt"
  "time"
  playerModel "kalaxia-game-api/model/player"
  serverModel "kalaxia-game-api/model/server"
)

func GetPlayer(id int16) *playerModel.Player {
  player := playerModel.Player{Id: id}
  if err := db.Select(&player); err != nil {
    return nil
  }
  return &player
}

func GetPlayerByUsername(username string) *playerModel.Player {
  player := playerModel.Player{Username: username}
  if err := db.Model(&player).Column("player.*", "Server").Where("username = ?", username).Select(); err != nil {
    return nil
  }
  return &player
}

func CreatePlayer(username string, server *serverModel.Server) *playerModel.Player {
  fmt.Println(server)

  player := playerModel.Player{
    Username: username,
    Pseudo: username,
    ServerId: server.Id,
    Server: server,
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
  }
  fmt.Println(player)
  if err := db.Insert(&player); err != nil {
    panic(err)
  }
  return &player
}
