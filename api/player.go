package api

import(
    "github.com/gorilla/mux"
    "github.com/gorilla/context"
    "net/http"
    "strconv"
    "time"
)

const(
    GenderMale = "male"
    GenderFemale = "female"
)

type(
    Player struct {
        Id uint16 `json:"id"`
        Username string `json:"-" sql:"type:varchar(180);not null;unique"`
        Pseudo string `json:"pseudo" sql:"type:varchar(180);not null;unique"`
        Gender string `json:"gender"`
        Avatar string `json:"avatar"`
        ServerId uint16 `json:"-"`
        Server *Server `json:"-"`
        FactionId uint16 `json:"-"`
        Faction *Faction `json:"faction"`
        IsActive bool `json:"is_active" sql:",notnull"`
        Wallet uint32 `json:"wallet"`
        Notifications Notifications `json:"notifications"`
        CurrentPlanet *Planet `json:"current_planet"`
        CurrentPlanetId uint16 `json:"-"`
        //Technologies []*PlayerTechnology `json:"technologies"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
    }
    Players []Player

    PlayerTechnology struct {
        TableName struct{} `json:"-" sql:"player__technologies"`

        Technology *Technology `json:"technology"`
        ResearchState *ResearchState `json:"research_state"`
    }
)

func GetCurrentPlayer(w http.ResponseWriter, r *http.Request) {
    SendJsonResponse(w, 200, context.Get(r, "player"))
}

func GetPlayer(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    currentPlayer := context.Get(r, "player").(*Player)

    SendJsonResponse(w, 200, getPlayer(uint16(id), currentPlayer.Id == uint16(id)))
}

func GetPlayerPlanets(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    player := getPlayer(uint16(id), false)

    SendJsonResponse(w, 200, player.getPlanets())
}

func UpdateCurrentPlanet(w http.ResponseWriter, r *http.Request) {
    data := DecodeJsonRequest(r)
    player := context.Get(r, "player").(*Player)

	player.changeCurrentPlanet(uint16(data["planet_id"].(float64)))
	player.update()
    
    w.WriteHeader(204)
    w.Write([]byte(""))
}

func RegisterPlayer(w http.ResponseWriter, r *http.Request) {
    data := DecodeJsonRequest(r)
    player := context.Get(r, "player").(*Player)
    if player.IsActive == true {
        panic(NewHttpException(http.StatusForbidden, "Player account is already active", nil))
    }
    player.register(
        data["pseudo"].(string),
        data["gender"].(string),
        data["avatar"].(string),
        uint16(data["faction_id"].(float64)),
        uint16(data["planet_id"].(float64)),
    )
    w.WriteHeader(204)
    w.Write([]byte(""))
}

func getPlayer(id uint16, isSelf bool) *Player {
    player := &Player{}
    if err := Database.Model(player).Column("player.*", "Faction", "Notifications", "CurrentPlanet").Where("player.id = ?", id).Select(); err != nil {
        panic(NewException("players.not_found", err))
    }
    if !isSelf {
        player.CurrentPlanet = nil
        player.Notifications = make(Notifications, 0)
    }
    return player
}

func (s *Server) getPlayerByUsername(username string) *Player {
    player := &Player{}
    if err := Database.
        Model(player).
        Column("player.*", "Server").
        Where("username = ?", username).
        Where("server_id = ?", s.Id).
        Select(); err != nil {
        return nil
    }
    return player
}

func (s *Server) createPlayer(username string) *Player {
    player := &Player{
        Username: username,
        Pseudo: username,
        ServerId: s.Id,
        Server: s,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    if err := Database.Insert(player); err != nil {
        panic(NewHttpException(500, "Player could not be created", err))
    }
    return player
}

func (p *Player) update() {
    if err := Database.Update(p); err != nil {
        panic(NewException("Player could not be updated ", err))
    }
}

func (p *Player) register(pseudo string, gender string, avatar string, factionId uint16, planetId uint16) {
    faction := getFaction(factionId)
    if faction == nil {
        panic(NewHttpException(404, "faction not found", nil))
    }
    planet := getPlayerPlanet(planetId, uint16(p.Id))
    if planet == nil {
        panic(NewHttpException(404, "planet not found", nil))
    }
    planet.PlayerId = p.Id
    planet.Player = p
    p.FactionId = faction.Id
    p.Faction = faction
    p.Pseudo = pseudo
    p.Avatar = avatar
    p.Gender = gender
    p.IsActive = true
    p.CurrentPlanet = planet
    p.CurrentPlanetId = planet.Id
    p.Wallet = 40000

    planet.increasePlayerRelation(p, 150)

    p.update()
    planet.update()
}

func (p *Player) updateWallet(amount int32) bool {
    if newAmount := int32(p.Wallet) + amount; newAmount >= 0 {
        p.Wallet = uint32(newAmount)
        return true
    }
    return false
}

func (p *Player) changeCurrentPlanet(planetId uint16) {
    planet := getPlanet(planetId)

    if planet.Player.Id != p.Id {
        panic(NewHttpException(403, "You do not own this planet", nil))
    }
    p.CurrentPlanet = planet
    p.CurrentPlanetId = planet.Id
    p.update()
}

func (p *Player) relocate() {
    if planets := p.getPlanets(); len(planets) > 0 {
        p.changeCurrentPlanet(planets[0].Id)
        return
    }
    planets := getFactionPlanetChoices(p.Faction.Id)
    planets[0].changeOwner(p)
    planets[0].update()
    p.notify(
        NotificationTypeDiplomacy,
        "notifications.diplomacy.relocated",
        map[string]interface{}{
            "planet": planets[0],
        },
    )
}