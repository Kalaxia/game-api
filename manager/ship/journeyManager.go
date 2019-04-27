package shipManager

import(
    "kalaxia-game-api/database"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/model"
    "kalaxia-game-api/utils"
    "kalaxia-game-api/manager"
    "encoding/json"
    "io/ioutil"
    "time"
)

/********************************************/
//global variable

var journeyTimeData model.TimeLawsContainer
var journeyRangeData model.RangeContainer

/********************************************/
// init

func init() {
    defer utils.CatchException() // defer is a "stack" (first in last out) so we call utils.CatchException() after the other defer of this function
    
    journeyTimeDataJSON, err := ioutil.ReadFile("../kalaxia-game-api/resources/journey_times.json")
    if err != nil {
        panic(exception.NewException("Can't open journey time configuration file", err))
    }
    if err := json.Unmarshal(journeyTimeDataJSON, &journeyTimeData); err != nil {
        panic(exception.NewException("Can't read journey time configuration file", err))
    }
    
    journeyRangeDataJSON, err := ioutil.ReadFile("../kalaxia-game-api/resources/journey_range.json")
    if err != nil {
        panic(exception.NewException("Can't open journey time configuration file", err))
    }
    if err := json.Unmarshal(journeyRangeDataJSON, &journeyRangeData); err != nil {
        panic(exception.NewException("Can't read journey time configuration file", err))
    }
    
    journeySteps := GetAllJourneySteps()
    now := time.Now()
    for _,step := range journeySteps { //< hic sunt dracones
        if step.Journey.CurrentStep.Id == step.Id { //< This if is the reason we defer and then go (see later comment)
            // we only treat the step if it is the current step
            // Indeed in finishStep we sheldur ( or do if the time has passed) finishStep for step.NextStep so we don't reshedur it
            // (It also resolve some proble with step deletion and concurency problem )
            if now.Before(step.TimeArrival){
                utils.Scheduler.AddTask(uint(now.Sub(step.TimeArrival).Seconds()), func () { // TODO review this uint
                    finishStep(step)
                }) //< Do I need to defer that to be safe ? (see comment below)
            } else {// if the time is already passed we directly execute it
                defer func () { go finishStep(step) }() //< finishStep delete step in DB
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

/********************************************/
// reponse to controller

/*------------------------------------------*/
// getter

func GetRange () model.RangeContainer{
    return journeyRangeData
}

func GetTimeLaws() model.TimeLawsContainer{
    return journeyTimeData
}

func SendFleetOnJourney (fleet *model.Fleet, data []interface{}) []*model.FleetJourneyStep {
    fleet.Journey = &model.FleetJourney {
        CreatedAt : time.Now(),
    }
    steps := createSteps(fleet, data, 0)
    fleet.Journey.EndedAt = steps[len(steps)-1].TimeArrival
    if err := database.Connection.Insert(fleet.Journey); err != nil {
		panic(exception.NewHttpException(500, "Journey could not be created", err))
    }
    fleet.JourneyId = fleet.Journey.Id
    insertSteps(steps)
    fleet.Journey.CurrentStep = steps[0]
    fleet.Journey.CurrentStepId = steps[0].Id
    UpdateJourney(fleet.Journey)
    fleet.Location = nil
    fleet.LocationId = 0
    UpdateFleet(fleet)
    
    now := time.Now()
    if now.Before(fleet.Journey.CurrentStep.TimeArrival){
        utils.Scheduler.AddTask(uint(now.Sub(fleet.Journey.CurrentStep.TimeArrival).Seconds()), func () {
            finishStep(fleet.Journey.CurrentStep)
        })
    } else {
        defer finishStep(fleet.Journey.CurrentStep) // if the time is already passed we directly execute it
    }
    return steps
}

func createSteps(fleet *model.Fleet, data []interface{}, firstNumber uint8) []*model.FleetJourneyStep {
    steps := make([]*model.FleetJourneyStep, len(data))

    previousPlanet := fleet.Location
    previousX := float64(fleet.Location.System.X)
    previousY := float64(fleet.Location.System.Y)
    previousTime := time.Now()

    for i, s := range data {
        stepData := s.(map[string]interface{})
        step := &model.FleetJourneyStep {
            Journey: fleet.Journey,

            PlanetStart: previousPlanet,
            MapPosXStart: previousX,
            MapPosYStart: previousY,

            StepNumber: uint32(firstNumber + 1 + uint8(i)),
        }
        if previousPlanet != nil {
            step.PlanetStartId = previousPlanet.Id
        }
    
        planetId := uint16(stepData["planetId"].(float64))

        if planetId != 0 {
            planet := manager.GetPlanet(planetId)

            step.PlanetFinal = planet
            step.PlanetFinalId = planet.Id
            step.MapPosXFinal = float64(planet.System.X)
            step.MapPosYFinal = float64(planet.System.Y)

            previousPlanet = planet
        } else {
            previousPlanet = nil

            step.MapPosXFinal = stepData["x"].(float64)
            step.MapPosYFinal = stepData["y"].(float64)
        }
        previousX = step.MapPosXFinal
        previousY = step.MapPosYFinal
        //TODO implement cooldown ?
        step.TimeStart = previousTime
        step.Order = stepData["order"].(string)
        step.TimeJump = step.TimeStart.Add( time.Duration(journeyTimeData.WarmTime.GetTimeForStep(step)) * time.Second )
        step.TimeArrival = step.TimeJump.Add( time.Duration(journeyTimeData.TravelTime.GetTimeForStep(step)) * time.Second)

        if !journeyRangeData.IsOnRange(step) {
            panic(exception.NewHttpException(400, "Step not in range", nil))
        }
        steps[i] = step
        previousTime = step.TimeArrival
    }
    return steps;
}

func AddStepsToJourney (fleet *model.Fleet, data []interface{}) []*model.FleetJourneyStep {
    journey := fleet.Journey
    if journey == nil || journey.CurrentStep == nil {
        panic(exception.NewHttpException(400, "Journey is already finished or does not exist", nil))
    }
    
    previousSteps := GetStepsByJourneyId(journey.Id)
    var oldLastStep *model.FleetJourneyStep = previousSteps[0]
    
    for _, step := range previousSteps {
        if step.StepNumber > oldLastStep.StepNumber { // we search for the higher step number this should be the last step
            oldLastStep = step
        }
    }
    if oldLastStep.NextStep != nil {
        panic(exception.NewHttpException(400, "The step with the higher step number is not the last step. Potential loop : abording", nil))
    }
    steps := createSteps(fleet, data, uint8(oldLastStep.StepNumber))
    
    journey.EndedAt =  steps[len(steps)-1].TimeArrival
    UpdateJourney(journey)
    insertSteps(steps)
    
    oldLastStep.NextStep = steps[0]
    oldLastStep.NextStepId = steps[0].Id
    UpdateJourneyStep(oldLastStep)
    
    return steps
}

func RemoveStepsAndFollowingFromJourney (journey *model.FleetJourney, step *model.FleetJourneyStep ) {
    
    if journey.Id != step.JourneyId{
        panic(exception.NewHttpException(400, "This step is not linked to the fleet journey", nil))
    }
    if journey.CurrentStepId == step.Id{
        panic(exception.NewHttpException(400, "current step cannot be canceled", nil))
    }
    if journey.CurrentStep.StepNumber >= step.StepNumber{
        panic(exception.NewHttpException(400, "cannot remove step with smaler step number that the current one ", nil))
    }
    
    steps := GetStepsByJourneyId(journey.Id)
    
    for _,stepR := range steps{
        if stepR.NextStepId == step.Id{
            stepR.NextStepId =0
            stepR.NextStep = nil
            journey.EndedAt = stepR.TimeArrival
            UpdateJourneyStep(stepR)
            UpdateJourney(journey)
            break
        }
    }
    
    deleteStepsRecursive(step)
}


/********************************************/
// internal logic

func finishStep(step *model.FleetJourneyStep) {
    var hasToDeleteJourney bool = false
    var journey *model.FleetJourney
    if step.Journey.CurrentStep != nil {
        if step.Journey.CurrentStep.Id == step.Id {
            processStepOrder(step)
            if step.NextStep != nil {
                if step.NextStep.StepNumber > step.StepNumber {
                    beginNextStep(step)
                } else {
                    panic(exception.NewException("Potential loop in Steps, abording",nil))
                }
            } else {
                finishJourney(step.Journey)
                hasToDeleteJourney = true
            }
        } else {
            // we do nothing this is an old step already finished
            // we could delete
        }
        
    } else {
        finishJourney(step.Journey)
        hasToDeleteJourney = true
    }
    
    // delete old step
    journey = step.Journey
    if err := database.Connection.Delete(step); err != nil {
        panic(exception.NewException("Journey step could not be deleted", err))
    }
    
    if hasToDeleteJourney {
        if err := database.Connection.Delete(journey); err != nil {
            panic(exception.NewException("Journey could not be deleted", err))
        }
    }
    
}

func processStepOrder(step *model.FleetJourneyStep) {
    switch (step.Order) {
        case model.FleetOrderPass:
            return
        case model.FleetOrderConquer:
            fleet := GetFleetByJourney(step.Journey)
            ConquerPlanet(fleet, manager.GetPlanet(step.PlanetFinalId))
            return
    }
}

func beginNextStep(step *model.FleetJourneyStep){
    step.Journey.CurrentStep = step.NextStep
    step.Journey.CurrentStepId = step.NextStepId
    UpdateJourney(step.Journey)
    
    var needToUpdateNextStep = false
    var defaultTime time.Time // This varaible is not set so it give the default time for unset varaible
    if step.NextStep.TimeStart == defaultTime { // TODO think of when we do this assignation
        //step.NextStep.TimeStart = step.TimeArrival.Add( time.Duration(journeyTimeData.Cooldown.GetTimeForStep(step.NextStep)) * time.Second )
        //TODO implement cooldown ?
        step.NextStep.TimeStart = step.TimeArrival
        needToUpdateNextStep = true
    }
    if step.NextStep.TimeJump == defaultTime { // TODO think of when we do this assignation
        step.NextStep.TimeJump = step.NextStep.TimeStart.Add( time.Duration(journeyTimeData.WarmTime.GetTimeForStep(step.NextStep)) * time.Second )
        needToUpdateNextStep = true
    }
    if step.NextStep.TimeArrival == defaultTime { // TODO think of when we do this assignation
        step.NextStep.TimeArrival = step.NextStep.TimeJump.Add( time.Duration(journeyTimeData.TravelTime.GetTimeForStep(step.NextStep)) * time.Second)
        needToUpdateNextStep = true
    }
    if needToUpdateNextStep{
        UpdateJourneyStep(step.NextStep)
    }
    
    nextStep := GetStep(step.NextStep.Id) //< Here I refresh the data because step.NextStep does not have step.NextStep.NextStep.*
    now := time.Now()
    if now.Before(nextStep.TimeArrival){
        utils.Scheduler.AddTask(uint(now.Sub(nextStep.TimeArrival).Seconds()), func () {
            finishStep(nextStep)
        })
    } else {
        defer finishStep(nextStep) // if the time is already passed we directly execute it
        // here I defer It (and not use go ) because I need that the current step is deleted before this one is
    }
    
}

func finishJourney(journey *model.FleetJourney){
    
    fleet := GetFleetOnJourney(journey.Id)
    journey = GetJourneyById(journey.Id) // we refresh the journey to get all the data
    // NOTE ^ this is nessesary
    
    fleet.Journey = nil
    fleet.JourneyId = 0
    if journey.CurrentStep.PlanetFinal != nil {
        fleet.Location = journey.CurrentStep.PlanetFinal
        fleet.LocationId = journey.CurrentStep.PlanetFinalId
    } else {
        fleet.MapPosX = journey.CurrentStep.MapPosXFinal
        fleet.MapPosY = journey.CurrentStep.MapPosYFinal
        fleet.Location = nil
        fleet.LocationId =0
    }
    journey.CurrentStep = nil
    journey.CurrentStepId = 0
    UpdateJourney(journey) //< I need to update it to breack the link to the current step so that I can remove the current step
    // And the I will be able to delete the journey (the step points to the journey)
    
    UpdateFleet(fleet)
}

func insertSteps (steps []*model.FleetJourneyStep) {
    // We must insert the last steps first, to reference them later in NextStep fields
    for i := len(steps)-1; i >= 0; i-- {
        step := steps[i]
        step.JourneyId = step.Journey.Id
        if i < len(steps) -1 {
            step.NextStep = steps[i+1]
            step.NextStepId = steps[i+1].Id
        }
        if err := database.Connection.Insert(step); err != nil {
    		panic(exception.NewHttpException(500, "Step could not be created", err))
        }
    }
}

func deleteStepsRecursive(step *model.FleetJourneyStep) {
    // delete Steps And Following Recursivly Utils
    // To call this function properly step must be deletable meaning the previous link must be already brocken
    nextStepId := step.NextStepId
    step.NextStepId =0
    step.NextStep = nil
    
    if err := database.Connection.Delete(step); err != nil {
        panic(exception.NewException("Journey step could not be deleted", err))
    }
    
    if nextStepId != 0 {
        deleteStepsRecursive(GetStep(nextStepId) )
    }
}

/********************************************/
// Data Base Connection
/*------------------------------------------*/
// getter

func GetJourneyById(idJourney uint16) *model.FleetJourney {
    journey := &model.FleetJourney{
        Id: idJourney,
    }
    if err := database.
        Connection.
        Model(journey).
        Column("CurrentStep.PlanetFinal").
        Where("fleet_journey.id = ?", idJourney  ).
        Select(); err != nil {
            panic(exception.NewException("journey not found on GetJourneyById", err))
        }
    return journey
}

func GetFleetOnJourney (idJourney uint16) *model.Fleet{
    var fleet model.Fleet
    if err := database.
        Connection.
        Model(&fleet).
        Column("Player", "Journey").
        Where("journey_id = ?", idJourney).
        Select(); err != nil {
            panic(exception.NewException("Fleet not found on GetFleetOnJourney", err))
        }
    return &fleet
}


func GetAllJourneySteps() []*model.FleetJourneyStep {
    var steps []*model.FleetJourneyStep
    if err := database.
        Connection.
        Model(&steps).
        Column("Journey.CurrentStep", "NextStep.PlanetStart", "NextStep.PlanetFinal").
        Select(); err != nil {
            panic(exception.NewException("steps not found on GetAllJourneySteps", err))
    }
    return steps
}


func GetStep(stepId uint16) *model.FleetJourneyStep {
    step := &model.FleetJourneyStep {
        Id : stepId,
    }
    
    if err := database.
        Connection.
        Model(step).
        Column("Journey.CurrentStep", "NextStep.Journey", "NextStep.PlanetFinal.System", "NextStep.PlanetStart.System","PlanetFinal.System", "PlanetStart.System").
        Where( "fleet_journey_step.id = ?",stepId).
        Select(); err != nil {
            panic(exception.NewException("step not found in GetStep", err))
    }
    return step
}

func GetStepsById(ids []uint16) []*model.FleetJourneyStep {
    var steps []*model.FleetJourneyStep
    if err := database.
        Connection.
        Model(&steps).
        Column("Journey.CurrentStep", "NextStep.PlanetFinal.System", "NextStep.PlanetStart.System", "PlanetFinal.System", "PlanetStart.System").
        WhereIn("fleet_journey_step.id IN ?", ids).
        Select(); err != nil {
            panic(exception.NewHttpException(404, "Fleets not found", err))
    }
    return steps
}

func GetStepsByJourneyId(journeyId uint16) []*model.FleetJourneyStep {
    var steps []*model.FleetJourneyStep
    if err := database.
        Connection.
        Model(&steps).
        Column("Journey.CurrentStep", "NextStep.PlanetFinal.System", "NextStep.PlanetStart.System","PlanetFinal.System", "PlanetStart.System").
        Where("fleet_journey_step.journey_id = ?",journeyId).
        Select(); err != nil {
            panic(exception.NewHttpException(404, "Fleets not found on GetStepsById", err))
    }
    return steps
}

/*------------------------------------------*/
// updater
func UpdateJourney(journey *model.FleetJourney){
    
    if err := database.Connection.Update(journey); err != nil {
        panic(exception.NewException("journey could not be updated on UpdateJourney", err))
    }
    
}

func UpdateJourneyStep(step *model.FleetJourneyStep){
    
    if err := database.Connection.Update(step); err != nil {
        panic(exception.NewException("step could not be updated on UpdateJourneyStep", err))
    }
    
}
