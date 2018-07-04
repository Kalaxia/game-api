package shipManager

import(
    "kalaxia-game-api/database"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/model"
)

func CreateShipModel(player *model.Player, data map[string]interface{}) *model.ShipModel {
    var frame model.ShipFrame
    var ok bool
    if frame, ok = framesData[data["frame"].(string)]; !ok {
        panic(exception.NewHttpException(400, "The given ship frame does not exists", nil))
    }
    slots := getSlotsData(data)
    shipType, stats := getShipModelInfo(frame, slots)
    shipModel := &model.ShipModel{
        Name: data["name"].(string),
        Type: shipType,
        FrameSlug: frame.Slug,
        Frame: &frame,
        Slots: slots,
        Stats: stats,
    }
    if err := database.Connection.Insert(shipModel); err != nil {
      panic(exception.NewHttpException(500, "Ship model could not be created", err))
    }
    playerShipModel := &model.ShipPlayerModel{
        ModelId: shipModel.Id,
        Model: shipModel,
        PlayerId: player.Id,
        Player: player,
    }
    if err := database.Connection.Insert(playerShipModel); err != nil {
      panic(exception.NewHttpException(500, "Player ship model could not be created", err))
    }
    createShipModelSlots(shipModel)
    return shipModel
}

func createShipModelSlots(shipModel *model.ShipModel) {
    for _, slot := range shipModel.Slots {
        slot.ModelId = shipModel.Id
        slot.Model = shipModel
        if err := database.Connection.Insert(&slot); err != nil {
          panic(exception.NewHttpException(500, "Ship model slot could not be created", err))
        }
    }
}

func GetShipPlayerModels(playerId uint16) []*model.ShipModel {
    var shipPlayerModels []model.ShipPlayerModel
    if err := database.Connection.Model(&shipPlayerModels).Column("Model").Where("player_id = ?", playerId).Select(); err != nil {
        panic(exception.NewHttpException(500, "Could not retrieve player ship models", err))
    }
    models := make([]*model.ShipModel, len(shipPlayerModels))
    for i, spm := range shipPlayerModels { models[i] = spm.Model }
    return models
}

func getShipModelInfo(frame model.ShipFrame, slots []model.ShipSlot) (string, map[string]uint16) {
    scores := make(map[string]uint8, 0)
    stats := make(map[string]uint16, 0)
    for stat, value := range frame.Stats {
        if storedStat, ok := stats[stat]; ok {
            storedStat += value
        } else {
            stats[stat] = value
        }
    }
    for _, slot := range slots {
        var module model.ShipModule
        var ok bool
        if module, ok = modulesData[slot.ModuleSlug]; !ok {
            continue
        }
        for score, amount := range module.Scores {
            if storedScore, ok := scores[score]; ok {
                storedScore += amount
            } else {
                scores[score] = amount
            }
        }
        for stat, value := range module.Stats {
            if storedStat, ok := stats[stat]; ok {
                storedStat += value
            } else {
                stats[stat] = value
            }
        }
    }
    return getShipModelType(scores), stats
}

func getShipModelType(scores map[string]uint8) string {
    shipType := ""
    var highestScore uint8
    highestScore = 0
    for score, value := range scores {
        if value > highestScore {
            shipType = score
            highestScore = value
        }
    }
    return shipType
}

func getSlotsData(data map[string]interface{}) []model.ShipSlot {
    slotsData := data["slots"].([]interface{})
    slots := make([]model.ShipSlot, len(slotsData))
    for i, slotData := range slotsData {
        slot := slotData.(map[string]interface{})
        slots[i] = model.ShipSlot{
            Position: uint8(slot["position"].(float64)),
            ModuleSlug: slot["module"].(string),
        }
    }
    return slots
}
