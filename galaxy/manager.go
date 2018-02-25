package galaxy

import(
    "kalaxia-game-api/base"
    "kalaxia-game-api/database"
    "kalaxia-game-api/server"
    "kalaxia-game-api/faction"
    "kalaxia-game-api/utils"
)

func GenerateMap(server *server.Server, factions []*faction.Faction, size uint16) *Map {
    gameMap := &Map{
        Server: server,
        ServerId: server.Id,
        Size: size,
    }
    if err := database.Connection.Insert(gameMap); err != nil {
        panic(err)
    }
    utils.GenerateMapSystems(gameMap, factions)
    return gameMap
}

func GetMapByServerId(serverId uint16) *Map {
    gameMap := &Map{ServerId: serverId}
    if err := database.Connection.Model(gameMap).Where("server_id = ?", serverId).Select(); err != nil {
        return nil
    }
    gameMap.Systems = GetMapSystems(gameMap.Id)
    return gameMap
}

func GetMapSystems(mapId uint16) []System {
    var systems []System
    if err := database.Connection.Model(&systems).Where("map_id = ?", mapId).Select(); err != nil {
        panic(err)
    }
    return systems
}

func GetSystem(id uint16) *System {
    system := System{Id: id}
    if err := database.Connection.Select(&system); err != nil {
        return nil
    }
    system.Planets = GetSystemPlanets(id)
    return &system
}

func GetSystemOrbits(id uint16) []SystemOrbit {
    var systemOrbits []SystemOrbit
    if err := database.Connection.Model(&systemOrbits).Where("system_id = ?", id).Select(); err != nil {
        panic(err)
    }
    return systemOrbits
}

func GetSystemPlanets(id uint16) []Planet {
    var planets []Planet
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

func GetPlayerPlanets(id uint16) []Planet {
    var planets []Planet
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

func GetPlanet(id uint16, playerId uint16) *Planet {
    var planet Planet
    if err := database.
        Connection.
        Model(&planet).
        Column("planet.*", "Player", "Resources", "System").
        Where("planet.id = ?", id).
        Select(); err != nil {
            return nil
    }
    relations := faction.GetPlanetRelations(planet.Id)
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

func getPlanetOwnerData(planet *Planet) {
    buildings, availableBuildings := base.GetPlanetBuildings(planet.Id)
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
