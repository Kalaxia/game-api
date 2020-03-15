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

const(
    FleetOrderPass = "pass"
    FleetOrderConquer = "conquer"
    FleetOrderColonize = "colonize"
    FleetOrderDeliver = "deliver"

    JourneyStepTypePlanetToPlanet = "planet_to_planet"
    JourneyStepTypeSameSystem = "same_system"
    JourneyStepTypeSamePlanet = "same_planet"
    JourneyStepTypePositionToPlanet = "position_to_planet"
    JourneyStepTypePlanetToPosition = "planet_to_position"
    JourneyStepTypePositionToPosition = "position_to_position"
)

type(
    FleetJourney struct {
        tableName struct{} `pg:"fleet__journeys"`

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
        tableName struct{} `pg:"fleet__journeys_steps"`
        
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
        Data map[string]interface{} `json:"data" pg:",use_zero,notnull"`
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
            Scheduler.AddTask(journey.CurrentStep.TimeArrival, func () {
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

func CalculateFleetTravelDuration(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)
	fleetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    fleet := getFleet(uint16(fleetId))
    
    if player.Id != fleet.Player.Id { // the player does not own the planet
		panic(NewHttpException(http.StatusForbidden, "", nil))
	}
    data := DecodeJsonRequest(r)
    
    SendJsonResponse(w, 200, fleet.calculateTravelDuration(data))
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
    f.validateJourney(steps)
    
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
        Scheduler.AddTask(f.Journey.CurrentStep.TimeArrival, func () {
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
        step := fleet.createStep(s.(map[string]interface{}), previousTime, previousPlace, uint32(firstNumber + 1 + uint8(i)))

        if !journeyRangeData.isOnRange(step) {
            panic(NewHttpException(400, "Step not in range", nil))
        }
        steps[i] = step
        previousPlace = step.EndPlace
        previousTime = step.TimeArrival
    }
    return steps;
}

func (f *Fleet) createStep(stepData map[string]interface{}, startTime time.Time, previousPlace *Place, stepNumber uint32) *FleetJourneyStep {
    step := &FleetJourneyStep {
        Journey: f.Journey,

        Order: stepData["order"].(string),

        StartPlace: previousPlace,
        StartPlaceId: previousPlace.Id,

        TimeStart: startTime,

        StepNumber: stepNumber,
    }
    if d, ok := stepData["data"]; ok {
        step.Data = d.(map[string]interface{})
    }

    endPlaceData := stepData["endPlace"].(map[string]interface{})

    if endPlaceData["planet"] != nil {
        planetId := uint16(endPlaceData["planet"].(map[string]interface{})["id"].(float64))
        planet := getPlanet(planetId)

        step.EndPlace = NewPlace(planet, float64(planet.System.X), float64(planet.System.Y))
    } else {
        coords := endPlaceData["coordinates"].(map[string]interface{})
        step.EndPlace = NewCoordinatesPlace(coords["x"].(float64), coords["y"].(float64))
    }
    warmTimeCoeff, travelTimeCoeff := f.getTimeCoefficients()
    step.EndPlaceId = step.EndPlace.Id
    step.TimeJump = step.TimeStart.Add(time.Duration(journeyTimeData.WarmTime.getTimeForStep(step, warmTimeCoeff)) * time.Minute)
    step.TimeArrival = step.TimeJump.Add(time.Duration(journeyTimeData.TravelTime.getTimeForStep(step, travelTimeCoeff)) * time.Minute)
    step.validate(f)
    return step
}

func (step *FleetJourneyStep) validate(f *Fleet) {
    if !step.hasValidOrder() {
        panic(NewHttpException(400, "fleet.journeys.invalid_order", nil))
    }
    if step.Order == FleetOrderDeliver && (step.Data["resources"] == nil || len(step.Data["resources"].([]interface{})) == 0) {
        panic(NewHttpException(400, "fleet.journeys.missing_resources_for_delivery", nil))
    }
    if step.Order == FleetOrderConquer && step.EndPlace.Planet.PlayerId == f.PlayerId {
        panic(NewHttpException(400, "fleet.journeys.cannot_conquer_your_own_planet", nil))
    }
    if step.Order == FleetOrderColonize && step.EndPlace.Planet.Population > 0 {
        panic(NewHttpException(400, "fleet.journeys.cannot_colonize_inhabited_planets", nil))
    }
}

func (f *Fleet) validateJourney(steps []*FleetJourneyStep) {
    f.validateDeliveryMissions(steps)
    f.validateColonizationMissions(steps)
}

// This method checks if the sum of the resources is really contained in the fleet cargo
func (f *Fleet) validateDeliveryMissions(steps []*FleetJourneyStep) {
    resources := make(map[string]uint16, 0)
    for _, s := range steps {
        if s.Order != FleetOrderDeliver {
            continue
        }
        for _, data := range s.Data["resources"].([]interface{}) {
            r := data.(map[string]interface{})
            resource := r["resource"].(string)
            quantity := uint16(r["quantity"].(float64))
            if _, ok := resources[resource]; !ok {
                resources[resource] = 0
            }
            resources[resource] += quantity
        }
    }
    for resource, quantity := range resources {
        if !f.hasResource(resource, quantity) {
            panic(NewHttpException(400, "fleet.delivery.not_enough_resources_for_all_steps", nil))
        }
    }
}

func (f *Fleet) validateColonizationMissions(steps []*FleetJourneyStep) {
    neededPopulation := uint16(0)
    for _, s := range steps {
        if s.Order != FleetOrderColonize {
            continue
        }
        passengers := uint16(s.Data["resources"].([]interface{})[0].(map[string]interface{})["quantity"].(float64))
        if passengers < 1000 {
            panic(NewHttpException(400, "fleet.colonization.need_more_colons", nil))
        }
        neededPopulation += passengers
    }
    if neededPopulation == 0 {
        return
    }
    if !f.hasResource("passengers", neededPopulation) {
        panic(NewHttpException(400, "fleet.colonization.not_enough_population", nil))
    }
}

func (step *FleetJourneyStep) hasValidOrder() bool {
    for _, o := range []string{ FleetOrderPass, FleetOrderConquer, FleetOrderColonize, FleetOrderDeliver } {
        if o == step.Order {
            return true
        }
    }
    return false
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
            continueJourney = fleet.conquerPlanet(getPlanet(s.EndPlace.PlanetId))
            if !continueJourney {
                fleet.delete()
            }
            break
        case FleetOrderColonize:
            continueJourney = true
            fleet := s.Journey.getFleet()
            fleet.colonizePlanet(getPlanet(s.EndPlace.PlanetId), int16(s.Data["resources"].([]interface{})[0].(map[string]interface{})["quantity"].(float64)))
            break
        case FleetOrderDeliver:
            fleet := s.Journey.getFleet()
            fleet.deliver(getPlanet(s.EndPlace.PlanetId), s.Data)
            continueJourney = true
            break
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
        Scheduler.AddTask(step.NextStep.TimeArrival, func () {
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

func (f *Fleet) calculateTravelDuration(data map[string]interface{}) map[string]time.Duration {
    var previousPlace *Place
    previousPlaceData := data["startPlace"].(map[string]interface{})

    if previousPlaceData["planet"] != nil {
        planetId := uint16(previousPlaceData["planet"].(map[string]interface{})["id"].(float64))

        planet := getPlanet(planetId)

        previousPlace = NewPlace(planet, float64(planet.System.X), float64(planet.System.Y))
    } else {
        coords := previousPlaceData["coordinates"].(map[string]interface{})
        previousPlace = NewCoordinatesPlace(coords["x"].(float64), coords["y"].(float64))
    }
    step := f.createStep(data, time.Now(), previousPlace, 1)
    warmTimeCoeff, travelTimeCoeff := f.getTimeCoefficients()
    return map[string]time.Duration{
        "warm": time.Duration(journeyTimeData.WarmTime.getTimeForStep(step, warmTimeCoeff)) * time.Minute,
        "travel": time.Duration(journeyTimeData.TravelTime.getTimeForStep(step, travelTimeCoeff)) * time.Minute,
    }   
}

func (f *Fleet) getTimeCoefficients() (warmTimeCoeff, travelTimeCoeff float64) {
    statToCoeff := func(coeff float64, stat uint16) float64 {
        c := float64(stat) / 100
        if coeff == 0 || coeff > c {
            return c
        }
        return coeff
    }
    for _, s := range f.getSquadrons() {
        warmTimeCoeff = statToCoeff(warmTimeCoeff, s.ShipModel.Stats[ShipStatCooldown]) * 10
        travelTimeCoeff = statToCoeff(travelTimeCoeff, s.ShipModel.Stats[ShipStatSpeed])
    }
    return
}

func (f *Fleet) isOnPlanet(p *Planet) bool {
    return f.Place != nil && f.Place.Planet != nil && (p == nil || p.Id == f.Place.PlanetId)
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

func (l *TimeDistanceLaw) getTime(distance, coeff float64) float64 {
    return (l.Constane + (l.Linear * distance) + (l.Quadratic * distance * distance)) / coeff
}

/*------------------------------------------*/
// TimeLawsConfig

func (l *TimeLawsConfig) getTimeForStep(s *FleetJourneyStep, coeff float64) float64 {
    // I use the Euclidian metric I am too lazy to implement other metric which will probaly not be used
    // NOTE this function will need to be modified if we add other "location"
    switch (s.getType()) {
        case JourneyStepTypeSamePlanet: return l.SamePlanet.getTime(s.getDistanceBetweenOrbitAndPlanet(), coeff)
        case JourneyStepTypeSameSystem: return l.SameSystem.getTime(s.getDistanceInsideSystem(), coeff)
        case JourneyStepTypePlanetToPlanet: return l.PlanetToPlanet.getTime(s.getDistanceBetweenPlanets(), coeff)
        case JourneyStepTypePlanetToPosition: return l.PlanetToPosition.getTime(s.getDistanceBetweenPlanetAndPosition(), coeff)
        case JourneyStepTypePositionToPlanet: return l.PositionToPlanet.getTime(s.getDistanceBetweenPositionAndPlanet(), coeff)
        case JourneyStepTypePositionToPosition: return l.PositionToPosition.getTime(s.getDistanceBetweenPositions(), coeff)
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