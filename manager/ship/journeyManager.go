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

func GetRange () model.RangeContainer{
    return journeyRangeData;
}

func GetTimeLaws() model.TimeLawsContainer{
    return journeyTimeData;
}

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
    
    journeySteps := GetAllJourneyStepsInternal();
    now := time.Now();
    for _,step := range journeySteps { //< hic sunt dracones
        if step.Journey.CurrentStep.Id == step.Id { //< This if is the reason we defer and then go (see later comment)
            // we only treat the step if it is the current step
            // Indeed in finishStep we sheldur ( or do if the time has passed) finishStep for step.NextStep so we don't reshedur it
            // (It also resolve some proble with step deletion and concurency problem )
            if now.Before(step.TimeArrival){
                utils.Scheduler.AddTask(uint(now.Sub(step.TimeArrival).Seconds()), func () { // TODO review this uint
                    finishStep(step);
                }); //< Do I need to defer that to be safe ? (see comment below)
            } else {// if the time is already passed we directly execute it
                defer func () { go finishStep(step); }(); //< finishStep delete step in DB
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
    
    nextStep := GetStepInternal(step.NextStep.Id); //< Here I refresh the data because step.NextStep does not have step.NextStep.NextStep.*
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
    
    fleet := GetFleetOnJourneyInternal(journey.Id);
    
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
    
    UpdateFleetInternal(fleet);
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

func GetFleetOnJourneyInternal (idJourney uint16) *model.Fleet{
    var fleet model.Fleet
    if err := database.
        Connection.
        Model(&fleet).
        Column("fleet.*", "Player", "Journey").
        Where("fleet.journey_id = ?", idJourney).
        Select(); err != nil {
            panic(exception.NewException("Fleet not found", err));
        }
    return &fleet;
}

func GetAllJourneySteps() []*model.FleetJourneyStep {
    var steps []*model.FleetJourneyStep;
    if err := database.
        Connection.
        Model(&steps).
        Column("*", "Journey", "NextStep", "Journey.CurrentStep", "NextStep.PlanetStart", "NextStep.PlanetFinal").
        Select(); err != nil {
            panic(exception.NewHttpException(404, "Fleets not found", err));
    }
    return steps;
}

func GetAllJourneyStepsInternal() []*model.FleetJourneyStep {
    var steps []*model.FleetJourneyStep;
    if err := database.
        Connection.
        Model(&steps).
        Column("*", "Journey", "NextStep", "Journey.CurrentStep", "NextStep.PlanetStart", "NextStep.PlanetFinal").
        Select(); err != nil {
            panic(exception.NewException("steps not found", err));
    }
    return steps;
}

func GetStep(stepId uint16) *model.FleetJourneyStep {
    var step model.FleetJourneyStep;
    if err := database.
        Connection.
        Model(&step).
        Column("step.*", "Journey", "NextStep", "Journey.CurrentStep", "NextStep.PlanetStart", "NextStep.PlanetFinal.System", "NextStep.PlanetStart.System", "NextStep.PlanetFinal.System","PlanetFinal", "PlanetStart","PlanetFinal.System", "PlanetStart.System").
        Where("step.id = ?",stepId).
        Select(); err != nil {
            panic(exception.NewHttpException(404, "Fleets not found", err));
    }
    return &step;
}

func GetStepInternal(stepId uint16) *model.FleetJourneyStep {
    var step model.FleetJourneyStep;
    if err := database.
        Connection.
        Model(&step).
        Column("step.*", "Journey", "NextStep", "Journey.CurrentStep", "NextStep.PlanetStart", "NextStep.PlanetFinal.System", "NextStep.PlanetStart.System", "NextStep.PlanetFinal.System","PlanetFinal", "PlanetStart","PlanetFinal.System", "PlanetStart.System").
        Where("step.id = ?",stepId).
        Select(); err != nil {
            panic(exception.NewException("step not found", err));
    }
    return &step;
}

func GetStepsById(ids []uint16) []*model.FleetJourneyStep {
    var steps []*model.FleetJourneyStep;
    if err := database.
        Connection.
        Model(&steps).
        Column("step.*", "Journey", "NextStep", "Journey.CurrentStep", "NextStep.PlanetStart", "NextStep.PlanetFinal.System", "NextStep.PlanetStart.System", "NextStep.PlanetFinal.System","PlanetFinal", "PlanetStart","PlanetFinal.System", "PlanetStart.System").
        WhereIn("step.id IN ?", ids).
        Select(); err != nil {
            panic(exception.NewHttpException(404, "Fleets not found", err));
    }
    return steps;
}

func GetStepsByJourneyId(journeyId uint16) []*model.FleetJourneyStep {
    var steps []*model.FleetJourneyStep;
    if err := database.
        Connection.
        Model(&steps).
        Column("step.*", "Journey", "NextStep", "Journey.CurrentStep", "NextStep.PlanetStart", "NextStep.PlanetFinal.System", "NextStep.PlanetStart.System", "NextStep.PlanetFinal.System","PlanetFinal", "PlanetStart","PlanetFinal.System", "PlanetStart.System").
        Where("step.journey_id = ?",journeyId).
        Select(); err != nil {
            panic(exception.NewHttpException(404, "Fleets not found", err));
    }
    return steps;
}

func SendFleetOnJourney (planetIds []uint16, x []float64, y []float64, fleet *model.Fleet) []*model.FleetJourneyStep{
    var steps []*model.FleetJourneyStep;
    journey := &model.FleetJourney{ // a new Journey is created
        CreatedAt : time.Now(),
    };
    
    steps = createStepStruct (fleet.Location, fleet.MapPosX, fleet.MapPosY, time.Now(),planetIds, x, y,0); //< we create the structurs
    
    journey.EndedAt =  steps[len(steps)-1].TimeArrival;
    if err := database.Connection.Insert(journey); err != nil {
		panic(exception.NewHttpException(500, "Journey could not be created", err))
    }
    
    insetFollowingStepInDBWithLinkCreation(steps,journey);
    
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
    }
    
    return steps;
}

func AddStepsToJourney (journey *model.FleetJourney, planetIds []uint16, x []float64, y []float64 ) []*model.FleetJourneyStep {
    if journey == nil || journey.CurrentStep == nil {
        panic(exception.NewHttpException(400, "Journey is already finished or does not exist", nil))
    }
    
    previousSteps := GetStepsByJourneyId(journey.Id);
    var oldLastStep *model.FleetJourneyStep = previousSteps[0];
    
    for _,step := range previousSteps {
        if step.StepNumber > oldLastStep.StepNumber { // we search for the higher step number; this should be the last step
            oldLastStep = step;
        }
    }
    if oldLastStep.NextStep != nil {
        panic(exception.NewHttpException(400, "The step with the higher step number is not the last step. Potential loop : abording", nil))
    }
    var steps []*model.FleetJourneyStep = createStepStruct (oldLastStep.PlanetFinal, oldLastStep.MapPosXFinal, oldLastStep.MapPosYFinal, oldLastStep.TimeArrival,planetIds, x, y,oldLastStep.StepNumber);
    
    journey.EndedAt =  steps[len(steps)-1].TimeArrival;
    UpdateJourney(journey);
    insetFollowingStepInDBWithLinkCreation(steps,journey);
    
    oldLastStep.NextStep = steps[0];
    oldLastStep.NextStepId = steps[0].Id;
    UpdateJourneyStep(oldLastStep);
    
    return steps;
}

func createStepStruct (firstPlanet *model.Planet, firstPosX float64, firstPosY float64, timeStart time.Time,planetIds []uint16, x []float64, y []float64, stepNumberOffset uint32) []*model.FleetJourneyStep {
    // Create the step structurs without the linking
    // The step is return in the correct order of execution
    // if the first location is not a planet, input firstPlanet as nil
    // The stepNumberOffset is the should be the StepNumber of pervious steps if there is otherwise input 0 ( infact you can input any number but prefer using 0)
    // if the stepNumberOffset is too low you may have error in the future
    if len(planetIds) != len(x) || len(x) != len(y)  {// the len(planetIds) != len(y) is unessesary
        panic(exception.NewHttpException(400, "Invalid input : planetIds, x and y should be of the same length", nil))
    }
    
    var steps []*model.FleetJourneyStep;
    var planetPrevId uint16 ;
    if firstPlanet != nil {
        planetPrevId = firstPlanet.Id;
    } else {
        planetPrevId = 0;
    }
    
    var PosPrevX float64 = firstPosX;
    var PosPrevY float64 = firstPosY;
    time_Last := timeStart;
    
    for i,_ := range planetIds {
        var step *model.FleetJourneyStep;
        if planetIds[i] != 0 {
            if planetPrevId !=0 {
                step = &model.FleetJourneyStep{
                    PlanetStartId : planetPrevId,
                    PlanetFinalId : planetIds[i],
                    StepNumber : uint32(i+1)+stepNumberOffset,
                };
            } else{
                step = &model.FleetJourneyStep{
                    MapPosXStart : PosPrevX,
                    MapPosYStart : PosPrevY,
                    PlanetFinalId : planetIds[i],
                    StepNumber : uint32(i+1)+stepNumberOffset,
                };
            }
        } else {
            if planetPrevId !=0 {
                step = &model.FleetJourneyStep{
                    PlanetStartId : planetPrevId,
                    MapPosXFinal : x[i],
                    MapPosYFinal : y[i],
                    StepNumber : uint32(i+1)+stepNumberOffset,
                };
            } else{
                step = &model.FleetJourneyStep{
                    MapPosXStart : PosPrevX,
                    MapPosYStart : PosPrevY,
                    MapPosXFinal : x[i],
                    MapPosYFinal : y[i],
                    StepNumber : uint32(i+1)+stepNumberOffset,
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
        time_Last = step.TimeArrival;
    }
    
    return steps;
}


func insetFollowingStepInDBWithLinkCreation (steps []*model.FleetJourneyStep,journey *model.FleetJourney) {
    for i := len(steps)-1; i >= 0; i-- { //< we read the table in reverse
        stepC := steps[i];
        stepC.Journey = journey;
        stepC.JourneyId = journey.Id;
        if i < len(steps) -1{ //< normaly this steps is already in the DB
            stepC.NextStep =  steps[i+1]
            stepC.NextStepId =  steps[i+1].Id
        } else {
            stepC.NextStep = nil;
            stepC.NextStepId = 0;
        }
        
        if err := database.Connection.Insert(stepC); err != nil {
    		panic(exception.NewHttpException(500, "Step could not be created", err));
        }
    }
}

func RemoveStepsAndFollowingFromJourney (journey *model.FleetJourney, step *model.FleetJourneyStep ) {
    
    if journey.Id != step.JourneyId{
        panic(exception.NewHttpException(400, "This step is not linked to the fleet journey", nil));
    }
    if journey.CurrentStepId == step.Id{
        panic(exception.NewHttpException(400, "current step cannot be canceled", nil));
    }
    if journey.CurrentStep.StepNumber >= step.StepNumber{
        panic(exception.NewHttpException(400, "cannot remove step with smaler step number that the current one ", nil));
    }
    
    steps := GetStepsByJourneyId(journey.Id);
    
    for _,stepR := range steps{
        if stepR.NextStepId == step.Id{
            stepR.NextStepId =0;
            stepR.NextStep = nil;
            journey.EndedAt = stepR.TimeArrival;
            UpdateJourneyStep(stepR);
            UpdateJourney(journey);
            break;
        }
    }
    
    deleteStepsAndFollowingRecursiveUtils(step);
}

func deleteStepsAndFollowingRecursiveUtils(step *model.FleetJourneyStep) {
    // To call this function properly step must be deletable meaning the previous link must be already brocken
    nextStepId := step.NextStepId;
    step.NextStepId =0;
    step.NextStep = nil;
    
    if err := database.Connection.Delete(step); err != nil {
        panic(exception.NewException("Journey step could not be deleted", err));
    }
    
    if nextStepId != 0 {
        deleteStepsAndFollowingRecursiveUtils(GetStep(nextStepId) );
    }
}
