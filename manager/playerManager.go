package manager

import(
    "time"
    "kalaxia-game-api/database"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/model"
)

func init() {
    utils.Scheduler.AddHourlyTask(func () { CalculatePlayersWage() })
}

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

func UpdatePlayer(player *model.Player) {
    if err := database.Connection.Update(player); err != nil {
        panic(exception.NewException("Player could not be updated", err))
    }
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
    player.Wallet = 0
    IncreasePlayerWallet(player, 40000)
    IncreasePlayerRelation(planet, player, 150)
    if err := database.Connection.Update(player); err != nil {
        panic(exception.NewHttpException(500, "Player could not be updated", err))
    }
    if err := database.Connection.Update(planet); err != nil {
        panic(exception.NewHttpException(500, "Planet could not be updated", err))
    }
}

func IncreasePlayerWallet(player *model.Player, amount  uint32) {
    if newAmount := player.Wallet + amount; newAmount >= 0 {
      player.Wallet = newAmount
    } else {
      player.Wallet = 0;
    }
}
func CalculatePlayerWage(player model.Player, wg $sync.WaitGroup) {
  defer wg.Done()
  defer utils.CatchException()
  baseWage := 50
  serviceWageRatio := 0.5
  wage := 0
  planets := GetPlayerPlanets(player.Id)
  for index, value := range planets {
    wage += baseWage + value.PlanetSettings.ServicesPoints*serviceWageRatio
  }
  IncreasePlayerWallet(&player, wage)
  if err := database.Connection.Update(player); err != nil {
      panic(exception.NewHttpException(500, "Player could not be updated", err))
  }
}

func CalculatePlayersWage() {
    nbPlayers, _ := database.Connection.Model(&model.Player{}).Count()
    limit := 20

    var wg sync.WaitGroup

    for offset := 0; offset < nbPlayers; offset +=limit {
        players := getPlayers(offset, limit)

        for _, player := range players {
          wg.Add(1)
          go CalculatePlayerWage(player, &wg)
        }
        wg.Wait()
    }
}

func getPlayers(offset int, limit int) []model.Player {
    var players []model.Players
    if err := database.Connection.
        Model($players).
        Column("").
        Limit(limit).
        Offset(offset).
        Select(); err != nil {
            panic(exception.NewException("Players could not be rertrieved", err))
    }
    return players
}
