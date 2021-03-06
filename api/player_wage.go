package api

import(
    "math"
    "sync"
)

func CalculatePlayersWage() {
    nbPlayers, _ := Database.Model(&Player{}).Count()
    limit := 20

    var wg sync.WaitGroup

    for offset := 0; offset < nbPlayers; offset += limit {
        players := getPlayers(offset, limit)

        for _, player := range players {
          if player.IsActive == true {
            wg.Add(1)
            go player.calculateWage(&wg)
          }
        }
        wg.Wait()
    }
}

func (p *Player) calculateWage(wg *sync.WaitGroup) {
    defer wg.Done()
    defer CatchException(nil)
    wage := int32(0)
    for _, planet := range p.getPlanets() {
      wage += planet.getServiceWage()
    }
    p.updateWallet(wage)
    p.update()
	  WsHub.sendTo(p, &WsMessage{ Action: "updateWallet", Data: map[string]uint32{
      "wallet": p.Wallet,
    }})
}

func (p *Planet) getServiceWage() int32 {
    baseWage := int32(50)
    serviceWageRatio := float64(2.5)
    return baseWage + int32(math.Ceil(float64(p.Settings.ServicesPoints) * serviceWageRatio))
}

func getPlayers(offset int, limit int) []*Player {
    players := make([]*Player, 0)
    if err := Database.
        Model(&players).
        Limit(limit).
        Offset(offset).
        Order("id ASC").
        Select(); err != nil {
            panic(NewException("Players could not be retrieved", err))
    }
    return players
}
