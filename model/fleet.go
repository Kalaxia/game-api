package model

import "time"

type(
    Fleet struct {
        TableName struct{} `json:"-" sql:"fleet__fleets"`

        Id uint16 `json:"id"`
        Player *Player `json:"player"`
        Location *Planet `json:"location"`
        Journey *FleetJourney `json:"journey"`
    }
    FleetJourney struct {
        TableName struct{} `json:"-" sql:"fleet__journeys"`

        Id uint16 `json:"id"`
        CreatedAt time.Time `json:"created_at"`
        EndedAt time.Time `json:"ended_at"`
    }
)
