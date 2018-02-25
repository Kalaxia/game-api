package player

import(
    "time"
)

type(
    Player struct {
        Id uint16 `json:"id"`
        Username string `json:"-" sql:"type:varchar(180);not null;unique"`
        Pseudo string `json:"pseudo" sql:"type:varchar(180);not null;unique"`
        ServerId uint16 `json:"-"`
        Server *interface{} `json:"-"`
        FactionId uint16 `json:"-"`
        Faction *interface{} `json:"faction"`
        IsActive bool `json:"is_active"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
    }
    Players []Player
)
