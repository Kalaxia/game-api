package shipManager


import(
    "kalaxia-game-api/database"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/model"
    "kalaxia-game-api/utils"
    "encoding/json"
    "io/ioutil"
    "time"
)

var journeyTimeData model.TimeLawsContainer
var journeyRangeData model.RangeContainer

func init() {
    defer utils.CatchException(); // defer is a "stack" (first in last out) so we call utils.CatchException() after the other defer of this function
    
    journeyTimeDataJSON, err := ioutil.ReadFile("../kalaxia-game-api/resources/journey_times.json")
    if err != nil {
        panic(exception.NewException("Can't open journey time configuration file", err));
    }
    if err := json.Unmarshal(journeyTimeDataJSON, &journeyTimeData); err != nil {
        panic(exception.NewException("Can't read journey time configuration file", err));
    }
    
    journeyRangeDataJSON, err := ioutil.ReadFile("../kalaxia-game-api/resources/journey_range.json")
    if err != nil {
        panic(exception.NewException("Can't open journey time configuration file", err));
    }
    if err := json.Unmarshal(journeyRangeDataJSON, &journeyRangeData); err != nil {
        panic(exception.NewException("Can't read journey time configuration file", err));
    }
    
    journeySteps := GetAllJourneySteps();
    now := time.Now();
    for _,step := range journeySteps {
        if step.Journey.CurrentStep.Id == step.Id { //< This if is the reason we defer and then go (see later comment) 
            // we only treat the step if it is the current step
            // Indeed in finishStep we sheldur ( or do if the time has passed) finishStep for step.NextStep so we don't reshedur it
            // (It also resolve some proble with step deletion and concurency problem )
            
            if now.Before(step.TimeArrival){
                utils.Scheduler.AddTask(uint(now.Sub(step.TimeArrival).Seconds()), func () { // TODO review this uint
                    finishStep(step);
                }); //< Do I need to defer that to be safe ? (see comment below)
            } else {// if the time is already passed we directly execute it
                defer func () {go finishStep(step);}();
                // defer for concurency reason 
                //
                // * In the case of go replacing defer 
                // * imagine that a journey has tow steps S0 and S1 which are both finished
                // * S0 is the current step
                // * finishStep(S0) is threated and call defer finishStep(S1) (see dunction finishStep and beginNextStep)
                // * but the treatment of  finishStep(S0) can be faster than for the loop to arrive to S1 
                // * So that step.Journey.CurrentStep.Id == step.Id is true and go finishStep(S1) can be called
                // * So S1 could be finished two times 
                //---------------------------------------------
                // But we do go after because all steps are in diffrent journey and can be finished simultaniously
            }
        }
    }
}

func finishStep(step *model.FleetJourneyStep) {
    var hasToDeleteJourney bool = false;
    if step.Journey.CurrentStep != nil {
        if step.Journey.CurrentStep.Id == step.Id {
            if step.NextStep != nil {
                if step.NextStep.StepNumber > step.StepNumber {
                    beginNextStep(step);
                } else {
                    panic(exception.NewException("Potential loop in Steps, abording",nil))
                }
            } else {
                finishJourney(step.Journey);
                hasToDeleteJourney = true;
                
            }
        } else {
            // we do nothing this is an old step already finished
            // we could delete 
            
        }
        
    } else {
        finishJourney(step.Journey);
        hasToDeleteJourney = true;
    }
    
    // delete old step
    journey := step.Journey;
    if err := database.Connection.Delete(step); err != nil {
        panic(exception.NewException("Journey step could not be deleted", err));
    }
    
    if hasToDeleteJourney {
        if err := database.Connection.Delete(journey); err != nil {
            panic(exception.NewException("Journey step could not be deleted", err));
        }
    }
}

func beginNextStep(step *model.FleetJourneyStep){
    step.Journey.CurrentStep = step.NextStep;
    step.Journey.CurrentStepId = step.NextStepId;
    
    UpdateJourney(step.Journey);
    
    var needToUpdateNextStep = false;
    
    var defaultTime time.Time // This varaible is not set so it give the default time for unset varaible
    if step.NextStep.TimeStart == defaultTime { // TODO think of when we do this assignation
        //step.NextStep.TimeStart = step.TimeArrival.Add( time.Duration(journeyTimeData.Cooldown.GetTimeForStep(step.NextStep)) * time.Second );
        //TODO implement cooldown ?
        step.NextStep.TimeStart = step.TimeArrival;
        needToUpdateNextStep = true;
    }
    if step.NextStep.TimeJump == defaultTime { // TODO think of when we do this assignation
        step.NextStep.TimeJump = step.NextStep.TimeStart.Add( time.Duration(journeyTimeData.WarmTime.GetTimeForStep(step.NextStep)) * time.Second );
        needToUpdateNextStep = true;
    }
    if step.NextStep.TimeArrival == defaultTime { // TODO think of when we do this assignation
        step.NextStep.TimeArrival = step.NextStep.TimeJump.Add( time.Duration(journeyTimeData.TravelTime.GetTimeForStep(step.NextStep)) * time.Second);
        needToUpdateNextStep = true;
    }
    if needToUpdateNextStep{
        UpdateJourneyStep(step.NextStep);
    }
    
    nextStep := GetJourneyStep(step.NextStep.Id); //< Here I refresh the data because step.NextStep does not have step.NextStep.NextStep.* 
    now := time.Now();
    if now.Before(nextStep.TimeArrival){
        utils.Scheduler.AddTask(uint(now.Sub(nextStep.TimeArrival).Seconds()), func () {
            finishStep(nextStep);
        });
    } else {
        defer finishStep(nextStep); // if the time is already passed we directly execute it
        // here I defer It (and not use go ) because I need that the current step is deleted before this one is 
    }
    
}

func finishJourney(journey *model.FleetJourney){
    
    fleet := GetFleetOnJourney(journey.Id);
    fleet.Journey = nil;
    fleet.JourneyId = 0;
    if journey.CurrentStep.PlanetFinal != nil {
        fleet.Location = journey.CurrentStep.PlanetFinal;
        fleet.LocationId = journey.CurrentStep.PlanetFinalId;
    } else {
        fleet.MapPosX = journey.CurrentStep.MapPosXFinal;
        fleet.MapPosY = journey.CurrentStep.MapPosYFinal;
        fleet.Location = nil;
        fleet.LocationId =0;
    }
    journey.CurrentStep = nil;
    journey.CurrentStepId = 0;
    UpdateJourney(journey); //< I need to update it to breack the link to the current step so that I can remove the current step
    // And the I will be able to delete the journey (the step points to the journey)
    
    UpdateFleet(fleet);
} 

func UpdateJourney(journey *model.FleetJourney){
    
    if err := database.Connection.Update(journey); err != nil {
        panic(exception.NewException("journey could not be updated", err));
    }
    
}
func UpdateJourneyStep(step *model.FleetJourneyStep){
    
    if err := database.Connection.Update(step); err != nil {
        panic(exception.NewException("step could not be updated", err));
    }
    
}

func GetFleetOnJourney (idJourney uint16) *model.Fleet{
    var fleet model.Fleet
    if err := database.
        Connection.
        Model(&fleet).
        Column("fleet.*", "Player", "Journey").
        Where("fleet.journey_id = ?", idJourney).
        Select(); err != nil {
            panic(exception.NewHttpException(404, "Fleet not found", err));
        }
    return &fleet;
}

func GetAllJourneySteps() []*model.FleetJourneyStep {
    var steps []*model.FleetJourneyStep;
    if err := database.
        Connection.
        Model(&steps).
        Column("step.*", "step.Journey", "step.NextStep", "step.Journey.CurrentStep", "step.NextStep.PlanetStart", "step.NextStep.PlanetFinal").
        Select(); err != nil {
            panic(exception.NewHttpException(404, "Fleets not found", err));
    }
    return steps;
}

func GetJourneyStep(id uint16) *model.FleetJourneyStep {
    var steps model.FleetJourneyStep;
    if err := database.
        Connection.
        Model(&steps).
        Column("step.*", "step.Journey", "step.NextStep", "step.Journey.CurrentStep", "step.NextStep.PlanetStart", "step.NextStep.PlanetFinal.System", "step.NextStep.PlanetStart.System", "step.NextStep.PlanetFinal.System").
        Select(); err != nil {
            panic(exception.NewHttpException(404, "Fleets not found", err));
    }
    return &steps;
}

func SendFleetOnJourney (planetIds []uint16, x []float64, y []float64, fleet *model.Fleet) []*model.FleetJourneyStep{
    var steps []*model.FleetJourneyStep
    
    journey := &model.FleetJourney{
        CreatedAt : time.Now(),
    };
    
    var planetPrevId uint16 = fleet.LocationId;
    var PosPrevX float64 = fleet.MapPosX;
    var PosPrevY float64 = fleet.MapPosY;
    time_Last := time.Now();
    for i,_ := range planetIds {
        var step *model.FleetJourneyStep;
        if planetIds[i] != 0 {
            if planetPrevId !=0 {
                step = &model.FleetJourneyStep{
                    PlanetStartId : planetPrevId,
                    PlanetFinalId : planetIds[i],
                    StepNumber : uint32(i+1),
                };
            } else{
                step = &model.FleetJourneyStep{
                    PlanetFinalId : planetIds[i],
                    MapPosXStart : PosPrevX,
                    MapPosYStart : PosPrevY,
                    StepNumber : uint32(i+1),
                };
            }
        } else{
            if planetPrevId !=0 {
                step = &model.FleetJourneyStep{
                    PlanetStartId : planetPrevId,
                    MapPosXFinal : x[i],
                    MapPosYFinal : y[i],
                    StepNumber : uint32(i+1),
                };
            } else{
                step = &model.FleetJourneyStep{
                    MapPosXStart : PosPrevX,
                    MapPosYStart : PosPrevY,
                    MapPosXFinal : x[i],
                    MapPosYFinal : y[i],
                    StepNumber : uint32(i+1),
                };
            }
        }
        
        //TODO implement cooldown ?
        step.TimeStart = time_Last;
    
        step.TimeJump = step.TimeStart.Add( time.Duration(journeyTimeData.WarmTime.GetTimeForStep(step)) * time.Second );
        
        step.TimeArrival = step.TimeJump.Add( time.Duration(journeyTimeData.TravelTime.GetTimeForStep(step)) * time.Second);
        if ! journeyRangeData.IsOnRange(step){
            panic(exception.NewHttpException(400, "Step not in range", nil));
        }
        steps = append(steps,step)
        planetPrevId = planetIds[i];
        PosPrevX = x[i];
        PosPrevY = y[i];
    }
    
    journey.EndedAt =  steps[len(steps)-1].TimeArrival;
    
    if err := database.Connection.Insert(journey); err != nil {
		panic(exception.NewHttpException(500, "Journey could not be created", err))
    }
    
    for i := len(steps)-1; i >= 0; i-- { //< we read the table in reverse
        stepC := steps[i];
        stepC.Journey = journey;
        stepC.JourneyId = journey.Id;
        if i < len(steps) -1{
            stepC.NextStep =  steps[i+1]
            stepC.NextStepId =  steps[i+1].Id
        } else {
            stepC.NextStep = nil;
            stepC.NextStepId = 0;
        }
        
        if err := database.Connection.Insert(stepC); err != nil {
    		panic(exception.NewHttpException(500, "Step could not be created", err))
        }
    }
    
    journey.CurrentStep = steps[0];
    journey.CurrentStepId = steps[0].Id;
    UpdateJourney(journey);
    
    fleet.Location = nil;
    fleet.LocationId =0;
    fleet.Journey = journey;
    fleet.JourneyId = journey.Id;
    
    UpdateFleet(fleet);
    
    now := time.Now();
    
    if now.Before(journey.CurrentStep.TimeArrival){
        utils.Scheduler.AddTask(uint(now.Sub(journey.CurrentStep.TimeArrival).Seconds()), func () {
            finishStep(journey.CurrentStep);
        });
    } else {
        defer finishStep(journey.CurrentStep); // if the time is already passed we directly execute it
        //
    }
    
    
    return steps;
}
