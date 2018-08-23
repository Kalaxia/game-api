package model

import( 
    "time"
    "math"
)

type(
    Fleet struct {
        TableName struct{} `json:"-" sql:"fleet__fleets" pg:",discard_unknown_columns"` //TEMP discard_unknown_columns

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
        //FistStep *FleetJourneyStep `json:"first_step"`
        //FistStepId uint16 `json:"-"` //TODO remove unsed code
        CurrentStep *FleetJourneyStep `json:"current_step"`
        CurrentStepId uint16 `json:"-"`
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
        //
        StepNumber uint32 `json:"step_number"`
    }
    
    TimeDistanceLaw struct {
        Constane float64 `json:"constante"`
        Linear float64 `json:"linear"`
        Quadratic float64 `json:"quadratic"`
    }
    
    TimeLawsConfig struct {
        SameSystem TimeDistanceLaw `json:"same_system"`
        PlanetToPlanet TimeDistanceLaw `json:"planet_to_planet"`
        PlanetToPosition TimeDistanceLaw `json:"planet_to_position"`
        PositionToPosition TimeDistanceLaw `json:"position_to_position"`
        PositionToPlanet TimeDistanceLaw `json:"position_to_planet"`
    }
    
    TimeLawsContainer struct{
        WarmTime TimeLawsConfig `json:"warm_time"`
        TravelTime TimeLawsConfig `json:"travel_time"`
        Cooldown TimeLawsConfig `json:"cooldown"` //< cooldown unimplemented
    }
    
    RangeContainer struct{
        SameSystem float64 `json:"same_system"`
        PlanetToPlanet float64 `json:"planet_to_planet"`
        PlanetToPosition float64 `json:"planet_to_position"`
        PositionToPosition float64 `json:"position_to_position"`
        PositionToPlanet float64 `json:"position_to_planet"`
    }
    
)
/********************************************/
// Methode 
/*------------------------------------------*/
//Fleets

func (fleet Fleet) IsOnPlanet () bool{
    return fleet.Location == nil;
}
func (fleet Fleet) IsOnJourney () bool{
    return fleet.Journey == nil;
}
func (fleet Fleet) IsOnMap () bool{
    booleanX := ! math.IsNaN(fleet.MapPosX) && fleet.MapPosX >= 0;
    booleanY := ! math.IsNaN(fleet.MapPosY) && fleet.MapPosY >= 0;
    return booleanX && booleanY;
}
func (fleet Fleet) GetPositionOnMap () ([2]float64 , bool) {
    if (fleet.IsOnMap()){
        return [2]float64{fleet.MapPosX,fleet.MapPosY} , true;
    } else {
        return [2]float64{math.NaN(),math.NaN()} , false;
    }
}

/*------------------------------------------*/
// FleetJourneyStep
func (journeyStep FleetJourneyStep) GetNextStep () FleetJourneyStep{
    return *(journeyStep.NextStep);
}

func (journeyStep FleetJourneyStep) GetDistance () float64{
    if journeyStep.PlanetStart != nil{
        if journeyStep.PlanetFinal != nil {
            if journeyStep.PlanetFinal.SystemId == journeyStep.PlanetStart.SystemId {
                return 0.;
            } else {
                distance := math.Pow(math.Pow(float64( float64(journeyStep.PlanetFinal.System.X) - float64(journeyStep.PlanetStart.System.X)),2.) + math.Pow(float64( float64(journeyStep.PlanetFinal.System.Y) - float64(journeyStep.PlanetStart.System.Y)),2.) , 0.5);
                return distance;
            }
        } else {
            distance := math.Pow(math.Pow(float64(journeyStep.MapPosXFinal - float64(journeyStep.PlanetStart.System.X)),2.) + math.Pow(float64(journeyStep.MapPosYFinal - float64(journeyStep.PlanetStart.System.Y)),2.) , 0.5);
            return distance;
        }
    } else{
        if journeyStep.PlanetFinal != nil {
            distance := math.Pow(math.Pow(float64(float64(journeyStep.PlanetFinal.System.X) - journeyStep.MapPosXStart),2.) + math.Pow(float64(float64(journeyStep.PlanetFinal.System.Y) - journeyStep.MapPosYStart),2.) , 0.5);
            return distance;
        } else {
            distance := math.Pow(math.Pow(float64(journeyStep.MapPosXFinal - journeyStep.MapPosXStart),2.) + math.Pow(float64(journeyStep.MapPosYFinal - journeyStep.MapPosYStart),2.) , 0.5);
            return distance;
        }
    }
}

/*------------------------------------------*/
// TimeDistanceLaw

func (law TimeDistanceLaw) GetTime(distance float64) float64{
    return law.Constane+ law.Linear*distance+ law.Quadratic * distance*distance;
}

/*------------------------------------------*/
// TimeLawsConfig

func (law TimeLawsConfig) GetTimeForStep(journeyStep *FleetJourneyStep) float64 {
    // I use the Euclidian metric; I am too lazy to implement other metric which will probaly not be used
    // NOTE this function will need to be modified if we add other "location"
    if journeyStep.PlanetStart != nil{
        if journeyStep.PlanetFinal != nil {
            if journeyStep.PlanetFinal.SystemId == journeyStep.PlanetStart.SystemId {
                return law.SameSystem.GetTime(0.);
            } else {
                distance := math.Pow(math.Pow(float64( float64(journeyStep.PlanetFinal.System.X) - float64(journeyStep.PlanetStart.System.X)),2.) + math.Pow(float64( float64(journeyStep.PlanetFinal.System.Y) - float64(journeyStep.PlanetStart.System.Y)),2.) , 0.5);
                return law.PlanetToPlanet.GetTime(distance);
            }
        } else {
            distance := math.Pow(math.Pow(float64(journeyStep.MapPosXFinal - float64(journeyStep.PlanetStart.System.X)),2.) + math.Pow(float64(journeyStep.MapPosYFinal - float64(journeyStep.PlanetStart.System.Y)),2.) , 0.5);
            return law.PlanetToPosition.GetTime(distance);
        }
    } else{
        if journeyStep.PlanetFinal != nil {
            distance := math.Pow(math.Pow(float64(float64(journeyStep.PlanetFinal.System.X) - journeyStep.MapPosXStart),2.) + math.Pow(float64(float64(journeyStep.PlanetFinal.System.Y) - journeyStep.MapPosYStart),2.) , 0.5);
            return law.PositionToPlanet.GetTime(distance);
        } else {
            distance := math.Pow(math.Pow(float64(journeyStep.MapPosXFinal - journeyStep.MapPosXStart),2.) + math.Pow(float64(journeyStep.MapPosYFinal - journeyStep.MapPosYStart),2.) , 0.5);
            return law.PositionToPosition.GetTime(distance);
        }
    }
}
/*------------------------------------------*/
// RangeContainer
func (rangeC RangeContainer) IsOnRange(journeyStep *FleetJourneyStep) bool {
    if journeyStep.PlanetStart != nil{
        if journeyStep.PlanetFinal != nil {
            if journeyStep.PlanetFinal.SystemId == journeyStep.PlanetStart.SystemId {
                return rangeC.SameSystem >= 0.;
            } else {
                distance := math.Pow(math.Pow(float64( float64(journeyStep.PlanetFinal.System.X) - float64(journeyStep.PlanetStart.System.X)),2.) + math.Pow(float64( float64(journeyStep.PlanetFinal.System.Y) - float64(journeyStep.PlanetStart.System.Y)),2.) , 0.5);
                return rangeC.PlanetToPlanet >= distance;
            }
        } else {
            distance := math.Pow(math.Pow(float64(journeyStep.MapPosXFinal - float64(journeyStep.PlanetStart.System.X)),2.) + math.Pow(float64(journeyStep.MapPosYFinal - float64(journeyStep.PlanetStart.System.Y)),2.) , 0.5);
            return rangeC.PlanetToPosition >= distance;
        }
    } else{
        if journeyStep.PlanetFinal != nil {
            distance := math.Pow(math.Pow(float64(float64(journeyStep.PlanetFinal.System.X) - journeyStep.MapPosXStart),2.) + math.Pow(float64(float64(journeyStep.PlanetFinal.System.Y) - journeyStep.MapPosYStart),2.) , 0.5);
            return rangeC.PositionToPlanet >= distance;
        } else {
            distance := math.Pow(math.Pow(float64(journeyStep.MapPosXFinal - journeyStep.MapPosXStart),2.) + math.Pow(float64(journeyStep.MapPosYFinal - journeyStep.MapPosYStart),2.) , 0.5);
            return rangeC.PositionToPosition >= distance;
        }
    }
}

/*------------------------------------------*/
// utils
