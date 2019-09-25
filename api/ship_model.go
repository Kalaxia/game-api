package api

import(
    "net/http"
    "io/ioutil"
    "encoding/json"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "strconv"
)

const ShipTypeFighter = "fighter"
const ShipTypeBomber = "bomber"
const ShipTypeFreighter = "freighter"
const ShipTypeCorvette = "corvette"
const ShipTypeFrigate = "frigate"

const ModuleTypeWeapon = "weapon"
const ModuleTypeEngine = "engine"
const ModuleTypeShield = "shield"
const ModuleTypeCargo = "cargo"

var framesData map[string]ShipFrame
var modulesData map[string]ShipModule

func InitShipConfiguration() {
    defer CatchException(nil)
    framesDataJSON, err := ioutil.ReadFile("/go/src/kalaxia-game-api/resources/ship_frames.json")
    if err != nil {
        panic(NewException("Can't open ship frames configuration file", err))
    }
    if err := json.Unmarshal(framesDataJSON, &framesData); err != nil {
        panic(NewException("Can't read ship frames configuration file", err))
    }
    modulesDataJSON, err := ioutil.ReadFile("/go/src/kalaxia-game-api/resources/ship_modules.json")
    if err != nil {
        panic(NewException("Can't open ship modules configuration file", err))
    }
    if err := json.Unmarshal(modulesDataJSON, &modulesData); err != nil {
        panic(NewException("Can't read ship modules configuration file", err))
    }
}

type(
	ShipFrame struct {
		Slug string `json:"slug"`
		Picture string `json:"picture"`
		Picto string `json:"picto"`
		Slots []ShipSlotPlan `json:"slots"`
		Stats map[string]uint16 `json:"stats"`
		Price []Price `json:"price"`
	}
	ShipModel struct {
		TableName struct{} `json:"-" sql:"ship__models"`

		Id uint `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
		PlayerId uint16 `json:"-"`
		Player *Player `json:"player"`
		FrameSlug string `json:"frame"`
		Frame *ShipFrame `json:"-" sql:"-"`
		Slots []ShipSlot `json:"slots" sql:"-"`
		Stats map[string]uint16 `json:"stats"`
		Price []Price `json:"price"`
	}
	ShipModule struct {
		Picture string `json:"picture"`
		PictureFlipX bool `json:"picture_flip_x"`
		PictureFlipY bool `json:"picture_flip_y"`
		Slug string `json:"slug"`
		Type string `json:"type"`
		Shape string `json:"shape"`
		Size string `json:"size"`
		Stats map[string]uint16 `json:"stats"`
		ShipStats map[string]uint16 `json:"ship_stats"`
		Scores map[string]uint8 `json:"scores"`
		Price []Price `json:"price"`
	}
	ShipPlayerModel struct {
		TableName struct{} `json:"-" sql:"ship__player_models"`

		PlayerId uint16 `json:"-"`
		Player *Player `json:"player"`
		ModelId uint `json:"-"`
		Model *ShipModel `json:"model"`
	}
	ShipPlayerModule struct {
		TableName struct{} `json:"-" sql:"ship__player_modules"`

		PlayerId uint16 `json:"-"`
		Player *Player `json:"player"`
		ModuleSlug string `json:"-"`
		Module *ShipModule `json:"module"`
	}
    ShipSlot struct {
        TableName struct{} `json:"-" sql:"ship__slots"`

        Id uint16 `json:"id"`
        ModelId uint `json:"-"`
        Model *ShipModel `json:"model"`
        Position uint8 `json:"position"`
        ModuleSlug string `json:"module"`
        Module *ShipModule `json:"-" sql:"-"`
    }
    ShipSlotPlan struct {
        Position uint8 `json:"position"`
        Shape string `json:"shape"`
        Size string `json:"size"`
    }
)

func GetPlayerShipModels(w http.ResponseWriter, r *http.Request) {
    SendJsonResponse(w, 200, context.Get(r, "player").(*Player).getShipModels())
}

func GetShipModel(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)

    SendJsonResponse(w, 200, context.Get(r, "player").(*Player).getShipModel(
        uint32(id),
    ))
}

func CreateShipModel(w http.ResponseWriter, r *http.Request) {
    SendJsonResponse(w, 201, context.Get(r, "player").(*Player).createShipModel(
        DecodeJsonRequest(r),
    ))
}

func (p *Player) createShipModel(data map[string]interface{}) *ShipModel {
    frame := extractFrame(data)
    slots := frame.extractSlotsData(data)
    shipType, stats := frame.getShipModelInfo(slots)
    shipModel := &ShipModel{
        Name: data["name"].(string),
        Type: shipType,
        PlayerId: p.Id,
        Player: p,
        FrameSlug: frame.Slug,
        Frame: frame,
        Slots: slots,
        Stats: stats,
    }
    shipModel.calculatePrices()
    if err := Database.Insert(shipModel); err != nil {
      panic(NewHttpException(500, "Ship model could not be created", err))
    }
    shipModel.createShipModelSlots()
    p.addShipModel(shipModel)
    return shipModel
}

func extractFrame(data map[string]interface{}) *ShipFrame {
    if frame, ok := framesData[data["frame"].(string)]; ok {
        return &frame
    }
    panic(NewHttpException(400, "The given ship frame does not exists", nil))
}

func (p *Player) addShipModel(sm *ShipModel) *ShipPlayerModel {
    playerShipModel := &ShipPlayerModel{
        ModelId: sm.Id,
        Model: sm,
        PlayerId: p.Id,
        Player: p,
    }
    if err := Database.Insert(playerShipModel); err != nil {
      panic(NewHttpException(500, "Player ship model could not be created", err))
    }
    return playerShipModel
}

func (sm *ShipModel) createShipModelSlots() {
    for _, slot := range sm.Slots {
        slot.ModelId = sm.Id
        slot.Model = sm
        if err := Database.Insert(&slot); err != nil {
          panic(NewHttpException(500, "Ship model slot could not be created", err))
        }
    }
}

func (p *Player) getShipModels() []*ShipModel {
    shipPlayerModels := make([]ShipPlayerModel, 0)
    if err := Database.Model(&shipPlayerModels).Relation("Model").Where("Model.player_id = ?", p.Id).Select(); err != nil {
        panic(NewHttpException(500, "Could not retrieve player ship models", err))
    }
    models := make([]*ShipModel, len(shipPlayerModels))
    for i, spm := range shipPlayerModels {
        spm.Model.loadSlots()
        models[i] = spm.Model
    }
    return models
}

func (p *Player) getShipModel(modelId uint32) *ShipModel {
    shipPlayerModel := &ShipPlayerModel{}
    if err := Database.Model(shipPlayerModel).Relation("Model").Where("Model.player_id = ?", p.Id).Where("Model.id = ?", modelId).Select(); err != nil {
        panic(NewHttpException(404, "Player ship model not found", err))
    }
    shipPlayerModel.Model.loadSlots()
    return shipPlayerModel.Model
}

func (frame *ShipFrame) getShipModelInfo(slots []ShipSlot) (string, map[string]uint16) {
    scores := make(map[string]uint8, 0)
    stats := make(map[string]uint16, 0)
    
    frame.addStats(stats)
    for _, slot := range slots {
        if module, ok := modulesData[slot.ModuleSlug]; ok {
            module.addScores(scores)
            module.addStats(stats)
        }
    }
    return getShipModelType(scores), stats
}

func (f *ShipFrame) addStats(stats map[string]uint16) {
    for stat, value := range f.Stats {
        if storedStat, ok := stats[stat]; ok {
            stats[stat] = storedStat + value
        } else {
            stats[stat] = value
        }
    }
}

func (m *ShipModule) addStats(stats map[string]uint16) {
    for stat, value := range m.ShipStats {
        if storedStat, ok := stats[stat]; ok {
            stats[stat] = storedStat + value
        } else {
            stats[stat] = value
        }
    }
}

func (m *ShipModule) addScores(scores map[string]uint8) {
    for score, amount := range m.Scores {
        if storedScore, ok := scores[score]; ok {
            scores[score] = storedScore + amount
        } else {
            scores[score] = amount
        }
    }
}

func getShipModelType(scores map[string]uint8) (shipType string) {
    highestScore := uint8(0)
    for score, value := range scores {
        if value > highestScore {
            shipType = score
            highestScore = value
        }
    }
    return
}

func (f *ShipFrame) extractSlotsData(data map[string]interface{}) []ShipSlot {
    slotsData := data["slots"].([]interface{})
    slots := make([]ShipSlot, len(slotsData))
    for i, s := range slotsData {
        slots[i] = f.extractSlot(s.(map[string]interface{}))
    }
    return slots
}

func (f *ShipFrame) extractSlot(data map[string]interface{}) (slot ShipSlot) {
    slot.Position = uint8(data["position"].(float64))
    
    if data["module"] == nil {
        return
    }
    if module, ok := modulesData[data["module"].(string)]; ok {
        slot.Module = &module
        slot.ModuleSlug = data["module"].(string)
        if !f.isValidSlot(slot) {
            panic(NewHttpException(400, "ships.invalid_slot", nil))
        }
        return
    }
    panic(NewHttpException(400, "Invalid module", nil))
}

func (f *ShipFrame) isValidSlot(slot ShipSlot) bool {
    for _, s := range f.Slots {
        if s.Position == slot.Position {
            return s.Shape == slot.Module.Shape && s.Size == slot.Module.Size
        }
    }
    return false
}

func (sm *ShipModel) loadSlots() {
    slots := make([]ShipSlot, 0)
    if err := Database.Model(&slots).Where("model_id = ?", sm.Id).Select(); err != nil {
        panic(NewHttpException(500, "Could not retrieve ship slots", err))
    }
    sm.Slots = slots
}

func (sm *ShipModel) calculatePrices() {
    sm.Price = make([]Price, 0)
    pricesMap := make(map[string]Price, 0)
    addPrice := func(price Price) {
        var priceType string
        if price.Type != PriceTypeResources {
            priceType = price.Type
        } else {
            priceType = price.Resource
        }
        if p, ok := pricesMap[priceType]; ok {
            p.Amount += price.Amount
        } else {
            pricesMap[priceType] = price
        }
    }
    for _, price := range sm.Frame.Price {
        addPrice(price)
    }
    for _, slot := range sm.Slots {
        if slot.Module == nil {
            continue
        }
        for _, price := range slot.Module.Price {
            addPrice(price)
        }
    }
    for _, price := range pricesMap {
        sm.Price = append(sm.Price, price)
    }
}

