package player

import(
  "time"
  "errors"
  "kalaxia-game-api/faction"
  "kalaxia-game-api/galaxy"
  "kalaxia-game-api/database"
  "kalaxia-game-api/server"
)

func GetPlayer(id uint16) *Player {
    var player Player
    if err := database.Connection.Model(&player).Column("player.*", "Faction").Where("player.id = ?", id).Select(); err != nil {
        return nil
    }
    return &player
}

func GetPlayerByUsername(username string, server *server.Server) *Player {
    player := Player{Username: username}
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

func CreatePlayer(username string, server *server.Server) *Player {
    player := player.Player{
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

func RegisterPlayer(player *Player, factionId uint16, planetId uint16) {
    faction := faction.GetFaction(factionId)
    if faction == nil {
        panic(errors.New("faction not found"))
    }
    planet := galaxy.GetPlanet(planetId, uint16(player.Id))
    if planet == nil {
        panic(errors.New("planet not found"))
    }
    planet.PlayerId = player.Id
    planet.Player = player
    player.FactionId = faction.Id
    player.Faction = faction
    player.IsActive = true
    faction.IncreasePlayerRelation(planet, player, 150)
    if err := database.Connection.Update(player); err != nil {
        panic(err)
    }
    if err := database.Connection.Update(planet); err != nil {
        panic(err)
    }
}
