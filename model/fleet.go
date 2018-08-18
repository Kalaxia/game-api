package model

import "time"

type(
    Fleet struct {
        TableName struct{} `json:"-" sql:"fleet__fleets"`

        Id uint16 `json:"id"`
        Player *Player `json:"player"`
        PlayerId uint16 `json:"-"`
        Location *Planet `json:"location"`
        LocationId uint16 `json:"-"`
        Journey *FleetJourney `json:"journey"`
        JourneyId uint16 `json:"-"`
        MapPosX float64 `json:"map_pos_x"`
        MapPosY float64 `json:"map_pos_y"`
    }
    FleetJourney struct {
        TableName struct{} `json:"-" sql:"fleet__journeys"`

        Id uint16 `json:"id"`
        CreatedAt time.Time `json:"created_at"`
        EndedAt time.Time `json:"ended_at"`
        FistStep *FleetJourneyStep `json:"first_step"`
        FistStepId uint16 `json:"-"`
    }
    FleetJourneyStep struct {
        TableName struct{} `json:"-" sql:"fleet__journeys_steps"`
        
        Id uint16 `json:"id"`
        Journey *FleetJourney `json:"journey"`
        JourneyId uint16 `json:"-"`
        NextStep *FleetJourneyStep `json:"next_step"`
        NextStepId uint16 `json:"-"`
        //
        PlanetStart *Planet `json:"planet_start"`
        PlanetStartId uint16 `json:"-"`
        MapPosXStart float64 `json:"map_pos_x_start"`
        MapPosYStart float64 `json:"map_pos_y_Start"`
        //
        PlanetFinal *Planet `json:"planet_final"`
        PlanetFinalId uint16 `json:"-"`
        MapPosXFinal float64 `json:"map_pos_x_final"`
        MapPosYFinal float64 `json:"map_pos_y_final"`
        //
        TimeStart time.Time `json:"time_start"`
        TimeJump time.Time `json:"time_jump"`
        TimeArrival time.Time `json:"time_arrival"`
    }
    
)
