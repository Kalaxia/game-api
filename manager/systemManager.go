package manager

import(
    "math"
    "kalaxia-game-api/database"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/model"
)

func GetMapSystems(mapId uint16) []model.System {
    var systems []model.System
    if err := database.Connection.Model(&systems).Where("map_id = ?", mapId).Select(); err != nil {
        panic(exception.NewHttpException(404, "Map not found", err))
    }
    return systems
}

func GetSectorSystems(starmap *model.Map, sector uint16) []model.System {
    sectorsPerLine := starmap.Size / starmap.SectorSize
    lineNumber := uint16(math.Ceil(float64((sector - 1) / sectorsPerLine)))

    systems := make([]model.System, 0)
    if err := database.
        Connection.
        Model(&systems).
        Column("Planets", "Planets.Player", "Planets.Player.Faction").
        Where("map_id = ?", starmap.Id).
        Where("x >= ?", (sector - ((lineNumber * sectorsPerLine) + 1)) * starmap.SectorSize).
        Where("x <= ?", (sector - (lineNumber * sectorsPerLine)) * starmap.SectorSize).
        Where("y >= ?", lineNumber * sectorsPerLine).
        Where("y <= ?", (lineNumber + 1) * sectorsPerLine).
        Select(); err != nil {
        panic(exception.NewHttpException(404, "Map not found", err))
    }
    return systems
}

func GetSystem(id uint16) *model.System {
    system := model.System{Id: id}
    if err := database.Connection.Select(&system); err != nil {
        return nil
    }
    system.Planets = GetSystemPlanets(id)
    return &system
}
