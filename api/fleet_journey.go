package api

import(
    "encoding/json"
    "net/http"
    "io/ioutil"
    "math"
    "time"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "strconv"
)

var journeyTimeData TimeLawsContainer
var journeyRangeData RangeContainer

const FleetOrderPass = "pass"
const FleetOrderConquer = "conquer"

type(
    FleetJourney struct {
        TableName struct{} `json:"-" sql:"fleet__journeys"`

        Id uint16 `json:"id"`
        CreatedAt time.Time `json:"created_at"`
        EndedAt time.Time `json:"ended_at"`
        CurrentStep *FleetJourneyStep `json:"current_step"`
        CurrentStepId uint16 `json:"-"`
        Steps []*FleetJourneyStep `json:"steps" sql:"-"`
    }
    FleetJourneyStep struct {
        TableName struct{} `json:"-" sql:"fleet__journeys_steps"`
        
        Id uint16 `json:"id"`
        Journey *FleetJourney `json:"-"`
        JourneyId uint16 `json:"-"`
        NextStep *FleetJourneyStep `json:"-"`
        NextStepId uint16 `json:"-"`
        //
        PlanetStart *Planet `json:"planet_start"`
        PlanetStartId uint16 `json:"-"`
        MapPosXStart float64 `json:"map_pos_x_start"`
        MapPosYStart float64 `json:"map_pos_y_start"`
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
        Order string `json:"order" sql:"mission_order"`
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

func InitFleetJourneys() {
    defer CatchException()
    journeyTimeDataJSON, err := ioutil.ReadFile("../kalaxia-game-api/resources/journey_times.json")
    if err != nil {
        panic(NewException("Can't open journey time configuration file", err))
    }
    if err := json.Unmarshal(journeyTimeDataJSON, &journeyTimeData); err != nil {
        panic(NewException("Can't read journey time configuration file", err))
    }
    
    journeyRangeDataJSON, err := ioutil.ReadFile("../kalaxia-game-api/resources/journey_range.json")
    if err != nil {
        panic(NewException("Can't open journey time configuration file", err))
    }
    if err := json.Unmarshal(journeyRangeDataJSON, &journeyRangeData); err != nil {
        panic(NewException("Can't read journey time configuration file", err))
    }
    
    journeys := getAllJourneys()
    now := time.Now()
    for _, journey := range journeys { //< hic sunt dracones
        if journey.CurrentStep == nil {
            journey.end()
            continue
        }
        // We do not retrieve the current step journey data in the SQL query, to avoid more JOIN statements for data we already have
        journey.CurrentStep.JourneyId = journey.Id
        journey.CurrentStep.Journey = journey
        if journey.CurrentStep.TimeArrival.After(now) {
            Scheduler.AddTask(uint(time.Until(journey.CurrentStep.TimeArrival).Seconds()), func () { // TODO review this uint
                journey.CurrentStep.end()
            }) //< Do I need to defer that to be safe ? (see comment below)
        } else {// if the time is already passed we directly execute it
            defer func () { go journey.CurrentStep.end() }() //< finishStep delete step in DB
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

func GetJourney (w http.ResponseWriter, r *http.Request){
	player := context.Get(r, "player").(*Player)
	fleetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	fleet := getFleet(uint16(fleetId))
	
	if player.Id != fleet.Player.Id { // the player does not own the planet
		panic(NewHttpException(http.StatusForbidden, "", nil))
	}
	if !fleet.isOnJourney() {
		panic(NewHttpException(400, "This journey has ended", nil))
	}
	
	SendJsonResponse(w, 200, fleet.Journey)
}

func GetFleetSteps(w http.ResponseWriter, r *http.Request){
    player := context.Get(r, "player").(*Player)
    fleetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	fleet := getFleet(uint16(fleetId))
    
    if player.Id != fleet.Player.Id { // the player does not own the planet
		panic(NewHttpException(http.StatusForbidden, "", nil))
	}
	if !fleet.isOnJourney() {
		panic(NewHttpException(400, "This journey has ended", nil))
	}
    
    SendJsonResponse(w, 200, fleet.Journey.getSteps())
}

func SendFleetOnJourney(w http.ResponseWriter, r *http.Request){
    player := context.Get(r, "player").(*Player)
	fleetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    fleet := getFleet(uint16(fleetId))
    
    if player.Id != fleet.Player.Id { // the player does not own the planet
		panic(NewHttpException(http.StatusForbidden, "", nil))
	}
    
    if fleet.isOnJourney() {
		panic(NewHttpException(400, "Fleet already on journey", nil))
    }
    
    if ships := fleet.getShips(); len(ships) == 0 {
        panic(NewHttpException(400, "journey.errors.empty_fleet", nil))
    }
    
    data := DecodeJsonRequest(r)
    
    SendJsonResponse(w, 201, fleet.travel(data["steps"].([]interface{})))
}

func AddStepsToJourney(w http.ResponseWriter, r *http.Request){
    player := context.Get(r, "player").(*Player)
	fleetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    fleet := getFleet(uint16(fleetId))
    
    if player.Id != fleet.Player.Id { // the player does not own the planet
		panic(NewHttpException(http.StatusForbidden, "", nil))
	}
    
    if !fleet.isOnJourney() {
		panic(NewHttpException(400, "Fleet is not on journey", nil))
	}
    
    data := DecodeJsonRequest(r)
    
    SendJsonResponse(w, 202, fleet.addJourneySteps(data["steps"].([]interface{})))
}

func RemoveFleetJourneyStep(w http.ResponseWriter, r *http.Request){
    // Cancel a journey form a setp, it remove this step and evryone after this one
    player := context.Get(r, "player").(*Player)
    idFleet, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    fleet := getFleet(uint16(idFleet))
    
    if player.Id != fleet.Player.Id{
        panic(NewHttpException(http.StatusForbidden, "", nil))
    }
    if !fleet.isOnJourney(){
        panic(NewHttpException(400, "Fleet is not on journey", nil))
    }
    
    stepId, _ := strconv.ParseUint(mux.Vars(r)["stepId"], 10, 16)
    step := getStep(uint16(stepId))
    
    if fleet.JourneyId != step.JourneyId{
        panic(NewHttpException(400, "This step is not linked to the fleet journey", nil))
    }
    if fleet.Journey.CurrentStepId == step.Id{
        panic(NewHttpException(400, "current step cannot be canceled", nil))
    }
    if fleet.Journey.CurrentStep.StepNumber >= step.StepNumber{
        panic(NewHttpException(400, "cannot remove step with smaler step number that the current one ", nil))
    }
    
    fleet.Journey.removeStep(step)
    SendJsonResponse(w,204,"Deleted");
    
}

func (f *Fleet) travel(data []interface{}) []*FleetJourneyStep {
    f.Journey = &FleetJourney {
        CreatedAt : time.Now(),
    }
	steps := createSteps(f, data, 0)
	
    f.Journey.EndedAt = steps[len(steps)-1].TimeArrival
	
    if err := Database.Insert(f.Journey); err != nil {
		panic(NewHttpException(500, "Journey could not be created", err))
    }
	insertSteps(steps)
	
    f.Journey.CurrentStep = steps[0]
    f.Journey.CurrentStepId = steps[0].Id
	f.Journey.update()
	
    f.JourneyId = f.Journey.Id
    f.Location = nil
    f.LocationId = 0
    f.update()
    
    if f.Journey.CurrentStep.TimeArrival.After(time.Now()) {
        Scheduler.AddTask(uint(time.Until(f.Journey.CurrentStep.TimeArrival).Seconds()), func () {
            f.Journey.CurrentStep.end()
        })
    } else {
        defer f.Journey.CurrentStep.end() // if the time is already passed we directly execute it
    }
    return steps
}

func createSteps(fleet *Fleet, data []interface{}, firstNumber uint8) []*FleetJourneyStep {
    steps := make([]*FleetJourneyStep, len(data))

    previousPlanet := fleet.Location
    previousX := float64(fleet.Location.System.X)
    previousY := float64(fleet.Location.System.Y)
    previousTime := time.Now()

    for i, s := range data {
        stepData := s.(map[string]interface{})
        step := &FleetJourneyStep {
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
            planet := getPlanet(planetId)

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
        step.TimeJump = step.TimeStart.Add(time.Duration(journeyTimeData.WarmTime.getTimeForStep(step)) * time.Minute)
        step.TimeArrival = step.TimeJump.Add(time.Duration(journeyTimeData.TravelTime.getTimeForStep(step)) * time.Minute)

        if !journeyRangeData.isOnRange(step) {
            panic(NewHttpException(400, "Step not in range", nil))
        }
        steps[i] = step
        previousTime = step.TimeArrival
    }
    return steps;
}

func (f *Fleet) addJourneySteps(data []interface{}) []*FleetJourneyStep {
    journey := f.Journey
    if journey == nil || journey.CurrentStep == nil {
        panic(NewHttpException(400, "Journey is already finished or does not exist", nil))
    }
    
    previousSteps := journey.getSteps()
    oldLastStep := previousSteps[0]
    
    for _, step := range previousSteps {
        if step.StepNumber > oldLastStep.StepNumber { // we search for the higher step number this should be the last step
            oldLastStep = step
        }
    }
    if oldLastStep.NextStep != nil {
        panic(NewHttpException(400, "The step with the higher step number is not the last step. Potential loop : abording", nil))
    }
    steps := createSteps(f, data, uint8(oldLastStep.StepNumber))
    
    journey.EndedAt =  steps[len(steps)-1].TimeArrival
    journey.update()
    insertSteps(steps)
    
    oldLastStep.NextStep = steps[0]
    oldLastStep.NextStepId = steps[0].Id
    oldLastStep.update()
    
    return steps
}

func (j *FleetJourney) removeStep(step *FleetJourneyStep) {
    if j.Id != step.JourneyId{
        panic(NewHttpException(400, "This step is not linked to the fleet journey", nil))
    }
    if j.CurrentStepId == step.Id {
        panic(NewHttpException(400, "current step cannot be canceled", nil))
    }
    if j.CurrentStep.StepNumber >= step.StepNumber{
        panic(NewHttpException(400, "cannot remove step with smaler step number that the current one ", nil))
    }
    
    for _,stepR := range j.getSteps() {
        if stepR.NextStepId == step.Id{
            stepR.NextStepId = 0
            stepR.NextStep = nil
            j.EndedAt = stepR.TimeArrival
            stepR.update()
            j.update()
            break
        }
    }
    step.deleteStepsRecursive()
}


/********************************************/
// internal logic

func (step *FleetJourneyStep) end() {
    defer CatchException()

    if step.Id != step.Journey.CurrentStepId {
        return
    }
    step.processOrder()
    if step.NextStep == nil {
        step.Journey.end()
        return
    }
    if step.NextStep.StepNumber > step.StepNumber {
        step.beginNextStep()
    } else {
        panic(NewException("Potential loop in Steps, abording",nil))
    }
}

func (s *FleetJourneyStep) processOrder() {
    switch (s.Order) {
        case FleetOrderPass:
            return
        case FleetOrderConquer:
            fleet := s.Journey.getFleet()
            fleet.conquerPlanet(getPlanet(s.PlanetFinalId))
            return
    }
}

func (step *FleetJourneyStep) beginNextStep() {
    step.delete()
    
    var needToUpdateNextStep = false
    var defaultTime time.Time // This varaible is not set so it give the default time for unset varaible
    if step.NextStep.TimeStart == defaultTime { // TODO think of when we do this assignation
        //step.NextStep.TimeStart = step.TimeArrival.Add( time.Duration(journeyTimeData.Cooldown.getTimeForStep(step.NextStep)) * time.Second )
        //TODO implement cooldown ?
        step.NextStep.TimeStart = step.TimeArrival
        needToUpdateNextStep = true
    }
    if step.NextStep.TimeJump == defaultTime { // TODO think of when we do this assignation
        step.NextStep.TimeJump = step.NextStep.TimeStart.Add( time.Duration(journeyTimeData.WarmTime.getTimeForStep(step.NextStep)) * time.Second )
        needToUpdateNextStep = true
    }
    if step.NextStep.TimeArrival == defaultTime { // TODO think of when we do this assignation
        step.NextStep.TimeArrival = step.NextStep.TimeJump.Add( time.Duration(journeyTimeData.TravelTime.getTimeForStep(step.NextStep)) * time.Second)
        needToUpdateNextStep = true
    }
    if needToUpdateNextStep{
        step.NextStep.update()
    }
    
    nextStep := getStep(step.NextStep.Id) //< Here I refresh the data because step.NextStep does not have step.NextStep.NextStep.*
    if nextStep.TimeArrival.After(time.Now()) {
        step.Journey.CurrentStep = nextStep
        step.Journey.CurrentStepId = nextStep.Id
        step.Journey.update()
        Scheduler.AddTask(uint(time.Until(nextStep.TimeArrival).Seconds()), func () {
            nextStep.end()
        })
    } else {
        defer nextStep.end() // if the time is already passed we directly execute it
        // here I defer It (and not use go ) because I need that the current step is deleted before this one is
    }
    
}

func (j *FleetJourney) end() {
    fleet := j.getFleet()
    j = getJourney(j.Id) // we refresh the journey to get all the data
    // NOTE ^ this is necessary
    
    fleet.Journey = nil
    fleet.JourneyId = 0
    if j.CurrentStep.PlanetFinal != nil {
        fleet.Location = j.CurrentStep.PlanetFinal
        fleet.LocationId = j.CurrentStep.PlanetFinalId
    } else {
        fleet.MapPosX = j.CurrentStep.MapPosXFinal
        fleet.MapPosY = j.CurrentStep.MapPosYFinal
        fleet.Location = nil
        fleet.LocationId = 0
    }
    fleet.update()

    j.delete()
}

func insertSteps(steps []*FleetJourneyStep) {
    // We must insert the last steps first, to reference them later in NextStep fields
    for i := len(steps)-1; i >= 0; i-- {
        step := steps[i]
        step.JourneyId = step.Journey.Id
        if i < len(steps) -1 {
            step.NextStep = steps[i+1]
            step.NextStepId = steps[i+1].Id
        }
        if err := Database.Insert(step); err != nil {
    		panic(NewHttpException(500, "Step could not be created", err))
        }
    }
}

func (s *FleetJourneyStep) deleteStepsRecursive() {
    if err := Database.Delete(s); err != nil {
        panic(NewException("Journey step could not be deleted", err))
    }
    if s.NextStepId != 0 {
        getStep(s.NextStepId).deleteStepsRecursive()
    }
}

/********************************************/
// Data Base Connection
/*------------------------------------------*/
// getter

func getJourney(id uint16) *FleetJourney {
    journey := &FleetJourney{}
    if err := Database.
        Model(journey).
        Column("CurrentStep.PlanetFinal").
        Where("fleet_journey.id = ?", id).
        Select(); err != nil {
            panic(NewException("journey not found on GetJourneyById", err))
        }
    return journey
}

func (j *FleetJourney) getFleet() *Fleet{
    fleet := &Fleet{}
    if err := Database.
        Model(fleet).
        Column("Player.Faction", "Journey").
        Where("journey_id = ?", j.Id).
        Select(); err != nil {
            panic(NewException("Fleet not found on GetFleetOnJourney", err))
        }
    fleet.Journey = j
    return fleet
}

func getAllJourneys() []*FleetJourney {
    journeys := make([]*FleetJourney, 0)
    if err := Database.
        Model(&journeys).
        Column("CurrentStep.PlanetStart", "CurrentStep.PlanetFinal", "CurrentStep.NextStep.PlanetStart", "CurrentStep.NextStep.PlanetFinal").
        Select(); err != nil {
            panic(NewException("Journeys could not be retrieved", err))
    }
    return journeys
}

func getStep(stepId uint16) *FleetJourneyStep {
    step := &FleetJourneyStep {
        Id : stepId,
    }
    if err := Database.
        Model(step).
        Column("Journey.CurrentStep", "NextStep.Journey", "NextStep.PlanetFinal.System", "NextStep.PlanetStart.System","PlanetFinal.System", "PlanetStart.System").
        Where( "fleet_journey_step.id = ?",stepId).
        Select(); err != nil {
            panic(NewException("step not found in GetStep", err))
    }
    return step
}

func getStepsById(ids []uint16) []*FleetJourneyStep {
    steps := make([]*FleetJourneyStep, 0)
    if err := Database.
        Model(&steps).
        Column("Journey.CurrentStep", "NextStep.PlanetFinal.System", "NextStep.PlanetStart.System", "PlanetFinal.System", "PlanetStart.System").
        WhereIn("fleet_journey_step.id IN ?", ids).
        Select(); err != nil {
            panic(NewHttpException(404, "Fleets not found", err))
    }
    return steps
}

func (j *FleetJourney) getSteps() []*FleetJourneyStep {
    steps := make([]*FleetJourneyStep, 0)
    if err := Database.
        Model(&steps).
        Column("Journey.CurrentStep", "NextStep.PlanetFinal.System", "NextStep.PlanetStart.System","PlanetFinal.System", "PlanetStart.System").
        Where("fleet_journey_step.journey_id = ?", j.Id).
        Select(); err != nil {
            panic(NewHttpException(404, "Fleets not found on GetStepsById", err))
    }
    return steps
}

func (j *FleetJourney) update() {
    if err := Database.Update(j); err != nil {
        panic(NewException("journey could not be updated on UpdateJourney", err))
    }
}

func (j *FleetJourney) delete() {
    if err := Database.Delete(j); err != nil {
        panic(NewException("journey could not be deleted", err))
    }
}

func (s *FleetJourneyStep) update() {
    if err := Database.Update(s); err != nil {
        panic(NewException("step could not be updated on UpdateJourneyStep", err))
    }
}

func (s *FleetJourneyStep) delete() {
    if err := Database.Delete(s); err != nil {
        panic(NewException("Journey step could not be deleted", err))
    }
}

func (f *Fleet) isOnPlanet () bool{
    return f.Location == nil
}

func (f *Fleet) isOnJourney () bool{
    return f.Journey != nil
}

func (f *Fleet) isOnMap () bool{
    booleanX := !math.IsNaN(f.MapPosX) && f.MapPosX >= 0
	booleanY := !math.IsNaN(f.MapPosY) && f.MapPosY >= 0
	
    return booleanX && booleanY
}

func (f *Fleet) GetPositionOnMap () ([2]float64 , bool) {
    if f.isOnMap() {
        return [2]float64{f.MapPosX, f.MapPosY}, true
    }
	return [2]float64{math.NaN(),math.NaN()}, false
}

/*------------------------------------------*/
// FleetJourneyStep
func (journeyStep *FleetJourneyStep) getNextStep() *FleetJourneyStep{
    return journeyStep.NextStep
}

func (journeyStep FleetJourneyStep) getDistance() float64{
    if journeyStep.PlanetStart != nil{
        if journeyStep.PlanetFinal != nil {
            if journeyStep.PlanetFinal.SystemId == journeyStep.PlanetStart.SystemId {
                return 0.
            }
			return math.Pow(math.Pow(float64(journeyStep.PlanetFinal.System.X) - float64(journeyStep.PlanetStart.System.X),2.) + math.Pow( float64(journeyStep.PlanetFinal.System.Y) - float64(journeyStep.PlanetStart.System.Y),2.) , 0.5)
        }
		return math.Pow(math.Pow(journeyStep.MapPosXFinal - float64(journeyStep.PlanetStart.System.X),2.) + math.Pow(journeyStep.MapPosYFinal - float64(journeyStep.PlanetStart.System.Y),2.) , 0.5)
	}
	if journeyStep.PlanetFinal != nil {
		return math.Pow(math.Pow(float64(journeyStep.PlanetFinal.System.X) - journeyStep.MapPosXStart,2.) + math.Pow(float64(journeyStep.PlanetFinal.System.Y) - journeyStep.MapPosYStart,2.) , 0.5)
	}
	return math.Pow(math.Pow(journeyStep.MapPosXFinal - journeyStep.MapPosXStart,2.) + math.Pow(journeyStep.MapPosYFinal - journeyStep.MapPosYStart,2.) , 0.5)
}

/*------------------------------------------*/
// TimeDistanceLaw

func (law *TimeDistanceLaw) getTime(distance float64) float64{
    return law.Constane + (law.Linear * distance) + (law.Quadratic * distance * distance)
}

/*------------------------------------------*/
// TimeLawsConfig

func (law *TimeLawsConfig) getTimeForStep(journeyStep *FleetJourneyStep) float64 {
    // I use the Euclidian metric I am too lazy to implement other metric which will probaly not be used
    // NOTE this function will need to be modified if we add other "location"
    if journeyStep.PlanetStart != nil{
        if journeyStep.PlanetFinal != nil {
            if journeyStep.PlanetFinal.SystemId == journeyStep.PlanetStart.SystemId {
                return law.SameSystem.getTime(0.)
            }
			return law.PlanetToPlanet.getTime(journeyStep.getDistance())
        }
		return law.PlanetToPosition.getTime(journeyStep.getDistance())
    }
	if journeyStep.PlanetFinal != nil {
		return law.PositionToPlanet.getTime(journeyStep.getDistance())
	}
	return law.PositionToPosition.getTime(journeyStep.getDistance())
}

func (rc *RangeContainer) isOnRange(journeyStep *FleetJourneyStep) bool {
    if journeyStep.PlanetStart != nil{
        if journeyStep.PlanetFinal != nil {
            if journeyStep.PlanetFinal.SystemId == journeyStep.PlanetStart.SystemId {
                return rc.SameSystem >= 0.
            }
			return rc.PlanetToPlanet >= journeyStep.getDistance()
        }
		return rc.PlanetToPosition >= journeyStep.getDistance()
    }
	if journeyStep.PlanetFinal != nil {
		return rc.PositionToPlanet >= journeyStep.getDistance()
	}
	return rc.PositionToPosition >= journeyStep.getDistance()
}