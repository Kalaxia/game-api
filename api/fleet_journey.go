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
        tableName struct{} `json:"-" pg:"fleet__journeys"`

        Id uint16 `json:"id"`
        CreatedAt time.Time `json:"created_at"`
        EndedAt time.Time `json:"ended_at"`

        StartPlace *Place `json:"start_place"`
        StartPlaceId uint32 `json:"-"`
        EndPlace *Place `json:"end_place"`
        EndPlaceId uint32 `json:"-"`

        CurrentStep *FleetJourneyStep `json:"current_step"`
        CurrentStepId uint16 `json:"-"`
        Steps []*FleetJourneyStep `json:"steps" pg:"-"`
    }
    FleetJourneyStep struct {
        tableName struct{} `json:"-" pg:"fleet__journeys_steps"`
        
        Id uint16 `json:"id"`
        Journey *FleetJourney `json:"-"`
        JourneyId uint16 `json:"-"`
        NextStep *FleetJourneyStep `json:"-"`
        NextStepId uint16 `json:"-"`
        //
        StartPlace *Place `json:"start_place"`
        StartPlaceId uint32 `json:"-"`
        EndPlace *Place `json:"end_place"`
        EndPlaceId uint32 `json:"-"`
        //
        TimeStart time.Time `json:"time_start"`
        TimeJump time.Time `json:"time_jump"`
        TimeArrival time.Time `json:"time_arrival"`
        //
        StepNumber uint32 `json:"step_number"`
        Order string `json:"order" pg:"mission_order"`
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
        CreatedAt: time.Now(),
    }
	steps := f.createSteps(data, 0)
    f.Journey.StartPlace = steps[0].StartPlace
    f.Journey.StartPlaceId = f.Journey.StartPlace.Id
    f.Journey.EndPlace = steps[len(steps) - 1].EndPlace
    f.Journey.EndPlaceId = f.Journey.EndPlace.Id
	
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
    f.Place = nil
    f.PlaceId = 0
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
    previousPlace := fleet.Place
    previousTime := time.Now()

    for i, s := range data {
        stepData := s.(map[string]interface{})
        step := &FleetJourneyStep {
            Journey: fleet.Journey,

            StartPlace: previousPlace,
            StartPlaceId: previousPlace.Id,

            StepNumber: uint32(firstNumber + 1 + uint8(i)),
        }
    
        planetId := uint16(stepData["planetId"].(float64))

        if planetId != 0 {
            planet := getPlanet(planetId)

            step.EndPlace = NewPlace(planet, float64(planet.System.X), float64(planet.System.Y))
        } else {
            step.EndPlace = NewCoordinatesPlace(stepData["x"].(float64), stepData["y"].(float64))
        }
        step.EndPlaceId = step.EndPlace.Id
        previousPlace = step.EndPlace
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

func (s *FleetJourneyStep) processOrder() (continueJourney bool) {
    switch (s.Order) {
        case FleetOrderPass:
            continueJourney = true
            break
        case FleetOrderConquer:
            fleet := s.Journey.getFleet()
            continueJourney = fleet.conquerPlanet(getPlanet(s.EndPlace.Planet.Id))
            if !continueJourney {
                fleet.delete()
            }
            break;
        default:
            continueJourney = true
            break
    }
    return
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
    fleet.Place = j.CurrentStep.EndPlace
    fleet.PlaceId = j.CurrentStep.EndPlaceId
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
        Relation("CurrentStep.EndPlace.Planet").
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
        Relation("CurrentStep.StartPlace.Planet").
        Relation("CurrentStep.EndPlace.Planet").
        Relation("CurrentStep.NextStep.StartPlace.Planet").
        Relation("CurrentStep.NextStep.EndPlace.Planet").
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
        Relation("NextStep.EndPlace.Planet.System").
        Relation("NextStep.StartPlace.Planet.System").
        Relation("EndPlace.Planet.System").
        Relation("StartPlace.Planet.System").
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
        Relation("NextStep.EndPlace.Planet.System").
        Relation("NextStep.StartPlace.Planet.System").
        Relation("EndPlace.Planet.System").
        Relation("StartPlace.Planet.System").
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
    return f.Place != nil && f.Place.Planet != nil
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
    return getDistance(float64(s.StartPlace.Planet.System.X), s.EndPlace.Coordinates.X, float64(s.StartPlace.Planet.System.Y), s.EndPlace.Coordinates.Y)
}

func (s *FleetJourneyStep) getDistanceBetweenPositionAndPlanet() float64 {
    return getDistance(s.StartPlace.Coordinates.X, float64(s.EndPlace.Planet.System.X), s.StartPlace.Coordinates.Y, float64(s.EndPlace.Planet.System.Y))
}

func (s *FleetJourneyStep) getDistanceBetweenPlanets() float64 {
    return getDistance(float64(s.StartPlace.Planet.System.X), float64(s.EndPlace.Planet.System.X), float64(s.StartPlace.Planet.System.Y), float64(s.EndPlace.Planet.System.Y))
}

func (s *FleetJourneyStep) getDistanceBetweenPositions() float64 {
	return getDistance(s.StartPlace.Coordinates.X, s.EndPlace.Coordinates.X, s.StartPlace.Coordinates.Y, s.EndPlace.Coordinates.Y)
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
    if s.StartPlace.Planet != nil{
        if s.EndPlace.Planet != nil {
            if s.EndPlace.Planet.Id == s.StartPlace.Planet.Id {
                return JourneyStepTypeSamePlanet
            }
            if s.EndPlace.Planet.SystemId == s.StartPlace.Planet.SystemId {
                return JourneyStepTypeSameSystem
            }
			return JourneyStepTypePlanetToPlanet
        }
		return JourneyStepTypePlanetToPosition
    }
	if s.EndPlace.Planet != nil {
		return JourneyStepTypePositionToPlanet
	}
	return JourneyStepTypePositionToPosition
}