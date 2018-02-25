package faction

type(
    Faction struct {
        TableName struct{} `json:"-" sql:"faction__factions"`

        Id uint16 `json:"id"`
        Name string `json:"name"`
        Description string `json:"description"`
        Color string `json:"color"`
        ServerId uint16 `json:"-"`
        Server *interface{} `json:"-"`
    }
    DiplomaticRelation struct {
        TableName struct{} `json:"-" sql:"diplomacy__relations"`

        Planet *interface{} `json:"planet"`
        PlanetId uint16 `json:"-"`
        Faction *Faction `json:"faction"`
        FactionId uint16 `json:"-"`
        Player *interface{} `json:"player"`
        PlayerId uint16 `json:"-"`
        Score int `json:"score"`
    }
)
