package manager

import(
    "time"
    "kalaxia-game-api/database"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/model"
)

func GetPlayer(id uint16) *model.Player {
    var player model.Player
    if err := database.Connection.Model(&player).Column("player.*", "Faction").Where("player.id = ?", id).Select(); err != nil {
        return nil
    }
    return &player
}

func GetPlayerByUsername(username string, server *model.Server) *model.Player {
    player := model.Player{Username: username}
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

func CreatePlayer(username string, server *model.Server) *model.Player {
    player := model.Player{
        Username: username,
        Pseudo: username,
        ServerId: server.Id,
        Server: server,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    if err := database.Connection.Insert(&player); err != nil {
        panic(exception.NewHttpException(500, "Player could not be created", err))
    }
    return &player
}

func RegisterPlayer(player *model.Player, factionId uint16, planetId uint16) {
    faction := GetFaction(factionId)
    if faction == nil {
        panic(exception.NewHttpException(404, "faction not found", nil))
    }
    planet := GetPlanet(planetId, uint16(player.Id))
    if planet == nil {
        panic(exception.NewHttpException(404, "planet not found", nil))
    }
    planet.PlayerId = player.Id
    planet.Player = player
    player.FactionId = faction.Id
    player.Faction = faction
    player.IsActive = true
    player.Money = 0
    IncreasePlayerMoney(player, 40000)
    IncreasePlayerRelation(planet, player, 150)
    if err := database.Connection.Update(player); err != nil {
        panic(exception.NewHttpException(500, "Player could not be updated", err))
    }
    if err := database.Connection.Update(planet); err != nil {
        panic(exception.NewHttpException(500, "Planet could not be updated", err))
    }
}

func IncreasePlayerMoney(player *model.Player, amount  uint64){
    int newAmount =player.Money + amount
    if newAmount := player.Money + amount; newAmount>=0 {
      player.Money = newAmount
    }
    else {
      player.Money = 0;
    }
    if err := database.Connection.Update(player); err != nil {
        panic(exception.NewHttpException(500, "Player could not be updated", err))
    }
}
