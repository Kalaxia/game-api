package model

import "time"

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
        FrameSlug string `json:"-"`
        Frame *ShipFrame `json:"frame" sql:"-"`
        Slots []ShipSlot `json:"slots" sql:"-"`
        Stats map[string]uint16 `json:"stats"`
    }
    ShipModule struct {
        Slug string `json:"slug"`
        Type string `json:"type"`
        Shape string `json:"shape"`
        Size string `json:"size"`
        Stats map[string]uint16 `json:"stats"`
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

        HangarId uint16 `json:"-"`
        Hangar *Planet `json:"hangar"`
        FleetId uint16 `json:"-"`
        Fleet *Fleet `json:"fleet"`
        ModelId uint `json:"-"`
        Model *ShipModel `json:"model"`
        CreatedAt time.Time `json:"created_at"`
        BuiltAt time.Time `json:"updated_at"`
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
