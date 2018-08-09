package model

import (
        "time"
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

type(
    ShipConstructionState struct {
        TableName struct{} `json:"-" sql:"ship__construction_states"`

        Id uint32 `json:"id"`
        CurrentPoints uint8 `json:"current_points" sql:",notnull"`
        Points uint8 `json:"points"`
    }
    ShipFrame struct {
        Slug string `json:"slug"`
        Slots []ShipSlotPlan `json:"slots"`
        Stats map[string]uint16 `json:"stats"`
        Price []Price `json:"price"`
    }
    ShipModel struct {
        TableName struct{} `json:"-" sql:"ship__models"`

        Id uint `json:"id"`
        Name string `json:"name"`
        Type string `json:"type"`
        FrameSlug string `json:"frame"`
        Frame *ShipFrame `json:"-" sql:"-"`
        Slots []ShipSlot `json:"slots" sql:"-"`
        Stats map[string]uint16 `json:"stats"`
        Price []Price `json:"price"`
    }
    ShipModule struct {
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
    Ship struct {
        TableName struct{} `json:"-" sql:"ship__ships"`

        Id uint32 `json:"id"`
        HangarId uint16 `json:"-"` //< if  the ship is not in a hangar the id wil be nil.
        Hangar *Planet `json:"hangar"` //< if  the ship is not in a hangar  the pointer will be nil
        FleetId uint16 `json:"-"` //< if  the ship is not in a fleet the id wil be nil
        Fleet *Fleet `json:"fleet"` //< if  the ship is not in a fleet the pointer will be nil
        //IsShipInFleet bool `json:"isShipInFleet"` //< Depreciated : this is used when a ship is try to be removed form the hangar in order to avoid teleporting ships
        ModelId uint `json:"-"` 
        Model *ShipModel `json:"model"`
        CreatedAt time.Time `json:"created_at"`
        ConstructionStateId uint32 `json:"-"`
        ConstructionState *ShipConstructionState `json:"construction_state"`
        
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
        Shape string `json:"shape"`
        Size string `json:"size"`
    }
)

func (sm *ShipModel) CalculatePrices() {
    sm.Price = make([]Price, 0)
    pricesMap := make(map[string]Price, 0)
    addPrice := func(price Price) {
        var priceType string
        if price.Type != PRICE_TYPE_RESOURCE {
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

