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
const JourneyStepTypePlanetToPlanet = "planet_to_planet"
const JourneyStepTypeSameSystem = "same_system"
const JourneyStepTypeSamePlanet = "same_planet"
const JourneyStepTypePositionToPlanet = "position_to_planet"
const JourneyStepTypePlanetToPosition = "planet_to_position"
const JourneyStepTypePositionToPosition = "position_to_position"

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
        SamePlanet TimeDistanceLaw `json:"same_planet"`
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
        SamePlanet float64 `json:"same_planet"`
        SameSystem float64 `json:"same_system"`
        PlanetToPlanet float64 `json:"planet_to_planet"`
        PlanetToPosition float64 `json:"planet_to_position"`
        PositionToPosition float64 `json:"position_to_position"`
        PositionToPlanet float64 `json:"position_to_planet"`
    }
)

func InitFleetJourneys() {
    defer CatchException(nil)
    initJourneyData()
    
    journeys := getAllJourneys()
    now := time.Now()
    for _, journey := range journeys { //< hic sunt dracones
        journey.Steps = journey.getSteps()
        if journey.CurrentStep == nil {
            journey.CurrentStep = journey.Steps[0]
            journey.CurrentStepId = journey.CurrentStep.Id
        }
        // We do not retrieve the current step journey data in the SQL query, to avoid more JOIN statements for data we already have
        journey.CurrentStep.JourneyId = journey.Id
        journey.CurrentStep.Journey = journey
        if journey.CurrentStep.TimeArrival.After(now) {
            Scheduler.AddTask(uint(time.Until(journey.CurrentStep.TimeArrival).Seconds()), func () {
                journey.CurrentStep.end()
            })
        } else {
            go journey.CurrentStep.end()
        }
    }
}

func initJourneyData() {
    journeyTimeDataJSON, err := ioutil.ReadFile("/go/src/kalaxia-game-api/resources/journey_times.json")
    if err != nil {
        panic(NewException("Can't open journey time configuration file", err))
    }
    if err := json.Unmarshal(journeyTimeDataJSON, &journeyTimeData); err != nil {
        panic(NewException("Can't read journey time configuration file", err))
    }
    
    journeyRangeDataJSON, err := ioutil.ReadFile("/go/src/kalaxia-game-api/resources/journey_range.json")
    if err != nil {
        panic(NewException("Can't open journey time configuration file", err))
    }
    if err := json.Unmarshal(journeyRangeDataJSON, &journeyRangeData); err != nil {
        panic(NewException("Can't read journey time configuration file", err))
    }
}

func GetJourney(w http.ResponseWriter, r *http.Request){
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
    
    if len(fleet.getSquadrons()) == 0 {
        panic(NewHttpException(400, "journey.errors.empty_fleet", nil))
    }
    
    data := DecodeJsonRequest(r)
    
    SendJsonResponse(w, 201, fleet.travel(data["steps"].([]interface{})))
}

func (f *Fleet) travel(data []interface{}) *FleetJourney {
    f.Journey = &FleetJourney {
        CreatedAt : time.Now(),
    }
	steps := f.createSteps(data, 0)
	
    f.Journey.EndedAt = steps[len(steps)-1].TimeArrival
	
    if err := Database.Insert(f.Journey); err != nil {
		panic(NewHttpException(500, "Journey could not be created", err))
    }
	insertSteps(steps)
    
    f.Journey.Steps = steps
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
    return f.Journey
}

func (fleet *Fleet) createSteps(data []interface{}, firstNumber uint8) []*FleetJourneyStep {
    steps := make([]*FleetJourneyStep, len(data))

    var previousPlanet *Planet
    var previousX, previousY float64

    if fleet.Location != nil {
        previousPlanet = fleet.Location
        previousX = float64(fleet.Location.System.X)
        previousY = float64(fleet.Location.System.Y)
    } else {
        previousPlanet = nil
        previousX = float64(fleet.MapPosX)
        previousY = float64(fleet.MapPosY)
    }
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

func (step *FleetJourneyStep) end() {
    defer CatchException(nil)
    step = getStep(step.Id)

    if step.Id != step.Journey.CurrentStepId {
        return
    }
    if success := step.processOrder(); !success {
        step.Journey.end()
        return
    }
    // The code duplication is on purpose because mission fail could have other consequences in the future
    if step.NextStep == nil {
        step.Journey.end()
        return
    }
    if step.NextStep.StepNumber > step.StepNumber {
        step.beginNextStep()
    } else {
        panic(NewException("Potential loop in Steps, aborting",nil))
    }
}

func (s *FleetJourneyStep) processOrder() bool {
    switch (s.Order) {
        case FleetOrderPass:
            return true
        case FleetOrderConquer:
            fleet := s.Journey.getFleet()
            return fleet.conquerPlanet(getPlanet(s.PlanetFinalId))
        default:
            return false
    }
}

func (step *FleetJourneyStep) beginNextStep() {
    step.Journey.CurrentStep = step.NextStep
    step.Journey.CurrentStepId = step.NextStep.Id
    step.Journey.update()
    step.delete()
    if step.NextStep.TimeArrival.After(time.Now()) {
        Scheduler.AddTask(uint(time.Until(step.NextStep.TimeArrival).Seconds()), func () {
            step.NextStep.end()
        })
    } else {
        step.NextStep.end()
    }
}

func (j *FleetJourney) end() {
    fleet := j.getFleet()
    fleet.Journey.CurrentStep = getStep(fleet.Journey.CurrentStepId)
    
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

/********************************************/
// Data Base Connection
/*------------------------------------------*/
// getter

func getJourney(id uint16) *FleetJourney {
    journey := &FleetJourney{}
    if err := Database.
        Model(journey).
        Relation("CurrentStep.PlanetFinal").
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
        Relation("Player.Faction").
        Relation("Journey").
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
        Relation("CurrentStep.PlanetStart").
        Relation("CurrentStep.PlanetFinal").
        Relation("CurrentStep.NextStep.PlanetStart").
        Relation("CurrentStep.NextStep.PlanetFinal").
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
        Relation("Journey.CurrentStep").
        Relation("NextStep.Journey").
        Relation("NextStep.PlanetFinal.System").
        Relation("NextStep.PlanetStart.System").
        Relation("PlanetFinal.System").
        Relation("PlanetStart.System").
        Where("fleet_journey_step.id = ?",stepId).
        Select(); err != nil {
            panic(NewException("step not found in GetStep", err))
    }
    return step
}

func (j *FleetJourney) getSteps() []*FleetJourneyStep {
    steps := make([]*FleetJourneyStep, 0)
    if err := Database.
        Model(&steps).
        Relation("Journey.CurrentStep").
        Relation("NextStep.PlanetFinal.System").
        Relation("NextStep.PlanetStart.System").
        Relation("PlanetFinal.System").
        Relation("PlanetStart.System").
        Where("fleet_journey_step.journey_id = ?", j.Id).
        Order("fleet_journey_step.step_number ASC").
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

func (f *Fleet) isOnPlanet() bool {
    return f.Location != nil
}

func (f *Fleet) isOnJourney() bool {
    return f.Journey != nil
}

func getDistance(x1, x2, y1, y2 float64) float64{
    return math.Pow(math.Pow((x2 - x1), 2.) + math.Pow(y2 - y1, 2.), 0.5)
}

func (s *FleetJourneyStep) getDistanceInsideSystem() float64 {
    return 0.
}

func (s *FleetJourneyStep) getDistanceBetweenOrbitAndPlanet() float64 {
    return 0.
}

func (s *FleetJourneyStep) getDistanceBetweenPlanetAndPosition() float64 {
    return getDistance(float64(s.PlanetStart.System.X), s.MapPosXFinal, float64(s.PlanetStart.System.Y), s.MapPosYFinal)
}

func (s *FleetJourneyStep) getDistanceBetweenPositionAndPlanet() float64 {
    return getDistance(s.MapPosXStart, float64(s.PlanetFinal.System.X), s.MapPosYStart, float64(s.PlanetFinal.System.Y))
}

func (s *FleetJourneyStep) getDistanceBetweenPlanets() float64 {
    return getDistance(float64(s.PlanetStart.System.X), float64(s.PlanetFinal.System.X), float64(s.PlanetStart.System.Y), float64(s.PlanetFinal.System.Y))
}

func (s *FleetJourneyStep) getDistanceBetweenPositions() float64 {
	return getDistance(s.MapPosXStart, s.MapPosXFinal, s.MapPosYStart, s.MapPosYFinal)
}

/*------------------------------------------*/
// TimeDistanceLaw

func (l *TimeDistanceLaw) getTime(distance float64) float64 {
    return l.Constane + (l.Linear * distance) + (l.Quadratic * distance * distance)
}

/*------------------------------------------*/
// TimeLawsConfig

func (l *TimeLawsConfig) getTimeForStep(s *FleetJourneyStep) float64 {
    // I use the Euclidian metric I am too lazy to implement other metric which will probaly not be used
    // NOTE this function will need to be modified if we add other "location"
    switch (s.getType()) {
        case JourneyStepTypeSamePlanet: return l.SamePlanet.getTime(s.getDistanceBetweenOrbitAndPlanet())
        case JourneyStepTypeSameSystem: return l.SameSystem.getTime(s.getDistanceInsideSystem())
        case JourneyStepTypePlanetToPlanet: return l.PlanetToPlanet.getTime(s.getDistanceBetweenPlanets())
        case JourneyStepTypePlanetToPosition: return l.PlanetToPosition.getTime(s.getDistanceBetweenPlanetAndPosition())
        case JourneyStepTypePositionToPlanet: return l.PositionToPlanet.getTime(s.getDistanceBetweenPositionAndPlanet())
        case JourneyStepTypePositionToPosition: return l.PositionToPosition.getTime(s.getDistanceBetweenPositions())
        default: panic(NewException("unknown step type", nil))
    }
}

func (rc *RangeContainer) isOnRange(s *FleetJourneyStep) bool {
    switch (s.getType()) {
        case JourneyStepTypeSamePlanet: return rc.SamePlanet >= s.getDistanceBetweenOrbitAndPlanet()
        case JourneyStepTypeSameSystem: return rc.SameSystem >= s.getDistanceInsideSystem()
        case JourneyStepTypePlanetToPlanet: return rc.PlanetToPlanet >= s.getDistanceBetweenPlanets()
        case JourneyStepTypePlanetToPosition: return rc.PlanetToPosition >= s.getDistanceBetweenPlanetAndPosition()
        case JourneyStepTypePositionToPlanet: return rc.PositionToPlanet >= s.getDistanceBetweenPositionAndPlanet()
        case JourneyStepTypePositionToPosition: return rc.PositionToPosition >= s.getDistanceBetweenPositions()
        default: panic(NewException("unknown step type", nil))
    }
}

func (s *FleetJourneyStep) getType() string {
    if s.PlanetStart != nil{
        if s.PlanetFinal != nil {
            if s.PlanetFinalId == s.PlanetStartId {
                return JourneyStepTypeSamePlanet
            }
            if s.PlanetFinal.SystemId == s.PlanetStart.SystemId {
                return JourneyStepTypeSameSystem
            }
			return JourneyStepTypePlanetToPlanet
        }
		return JourneyStepTypePlanetToPosition
    }
	if s.PlanetFinal != nil {
		return JourneyStepTypePositionToPlanet
	}
	return JourneyStepTypePositionToPosition
}