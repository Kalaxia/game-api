package api

import(
	"testing"
	"time"
)

func TestEndStep(t *testing.T) {
	InitDatabaseMock()

	player := getPlayerMock(getFactionMock())
	fleet := getFleetMock(1, player)

	fleet.Journey = getJourneyMock(1)
	fleet.JourneyId = fleet.Journey.Id

	// fleet.Journey.CurrentStep.end()

	// if fleet.Journey.CurrentStep.Id == 1 {
	// 	t.Errorf("Current step did not change")
	// }
	// if fleet.Journey.CurrentStep.Id == 2 {
	// 	t.Errorf("Current step is not the next step")
	// }
}

func TestIsOnJourney(t *testing.T) {
	travellingFleet := &Fleet{ Journey: getJourneyMock(1) }
	orbitingFleet := &Fleet{}

	if !travellingFleet.isOnJourney() {
		t.Errorf("Travelling fleet is on journey")
	}
	if orbitingFleet.isOnJourney() {
		t.Errorf("Orbiting fleet is not on journey")
	}
}

func TestIsOnPlanet(t *testing.T) {
	travellingFleet := &Fleet{}
	orbitingFleet := &Fleet{ Location: getPlayerPlanetMock(getPlayerMock(getFactionMock())) }

	if travellingFleet.isOnPlanet() {
		t.Errorf("Travelling fleet is not on planet")
	}
	if !orbitingFleet.isOnPlanet() {
		t.Errorf("Orbiting fleet is on planet")
	}
}

func TestGetDistanceBetweenPlanets(t *testing.T) {
	initJourneyData()
	s := &FleetJourneyStep{
		PlanetStartId: 1,
		PlanetStart: &Planet{
			SystemId: 1,
			System: &System {
				X: 91,
				Y: 16,
			},
		},
		PlanetFinalId: 2,
		PlanetFinal: &Planet{
			SystemId: 2,
			System: &System {
				X: 75,
				Y: 20,
			},
		},
	}
	if distance := s.getDistanceBetweenPlanets(); distance != 16.49242250247064234259 {
		t.Errorf("Journey step from planet to planet should be 16.49242250247064234259, not %.20f", distance)
	}
	if s.getType() != JourneyStepTypePlanetToPlanet {
		t.Errorf("Step is planet to planet")
	}
	if time := journeyTimeData.TravelTime.getTimeForStep(s); time != 32.98484500494128468517 {
		t.Errorf("Time should be 32.98484500494128468517, not %.20f", time)
	}
	if journeyRangeData.isOnRange(s) {
		t.Errorf("Step should not be on range")
	}
}

func TestGetDistanceBetweenOrbitAndPlanet(t *testing.T) {
	initJourneyData()
	s := &FleetJourneyStep{
		PlanetStartId: 1,
		PlanetStart: &Planet{},
		PlanetFinalId: 1,
		PlanetFinal: &Planet{},
	}
	if distance := s.getDistanceBetweenOrbitAndPlanet(); distance != 0. {
		t.Errorf("Distance should be 0, not %.20f", distance)
	}
	if s.getType() != JourneyStepTypeSamePlanet {
		t.Errorf("Step is orbit to planet")
	}
	if time := journeyTimeData.TravelTime.getTimeForStep(s); time != 2. {
		t.Errorf("Time should be 0., not %.20f", time)
	}
	if !journeyRangeData.isOnRange(s) {
		t.Errorf("Step should be on range")
	}
}

func TestGetDistanceInsideSystem(t *testing.T) {
	initJourneyData()
	s := &FleetJourneyStep{
		PlanetStartId: 1,
		PlanetStart: &Planet{
			SystemId: 1,
		},
		PlanetFinalId: 2,
		PlanetFinal: &Planet{
			SystemId: 1,
		},
	}
	if distance := s.getDistanceInsideSystem(); distance != 0. {
		t.Errorf("Journey step in same system should be 0, not %.20f", distance)
	}
	if sType := s.getType(); sType != JourneyStepTypeSameSystem {
		t.Errorf("Step is same system, not %s", sType)
	}
	if time := journeyTimeData.TravelTime.getTimeForStep(s); time != 5. {
		t.Errorf("Time should be 5., not %.20f", time)
	}
	if !journeyRangeData.isOnRange(s) {
		t.Errorf("Step should be on range")
	}
}

func TestGetDistanceBetweenPlanetAndPosition(t *testing.T) {
	initJourneyData()
	s := &FleetJourneyStep{
		PlanetStart: &Planet{
			System: &System {
				X: 4,
				Y: 16,
			},
		},
		MapPosXFinal: 8,
		MapPosYFinal: 20,
	}
	if distance := s.getDistanceBetweenPlanetAndPosition(); distance != 5.65685424949238058190 {
		t.Errorf("Distance should be 5.65685424949238058190, not %.20f", distance)
	}
	if s.getType() != JourneyStepTypePlanetToPosition {
		t.Errorf("Step is position to planet")
	}
	if time := journeyTimeData.TravelTime.getTimeForStep(s); time != 11.31370849898476116380 {
		t.Errorf("Time should be 11.31370849898476116380, not %.20f", time)
	}
	if !journeyRangeData.isOnRange(s) {
		t.Errorf("Step should be on range")
	}
}

func TestGetDistanceBetweenPositionAndPlanet(t *testing.T) {
	initJourneyData()
	s := &FleetJourneyStep{
		MapPosXStart: 30,
		MapPosYStart: 40,
		PlanetFinal: &Planet{
			System: &System{
				X: 35,
				Y: 45,
			},
		},
	}
	if distance := s.getDistanceBetweenPositionAndPlanet(); distance != 7.07106781186547550533 {
		t.Errorf("Distance should be 7.07106781186547550533, not %.20f", distance)
	}
	if sType := s.getType(); sType != JourneyStepTypePositionToPlanet {
		t.Errorf("Step is position to planet, not %s", sType)
	}
	if time := journeyTimeData.TravelTime.getTimeForStep(s); time != 14.14213562373095101066 {
		t.Errorf("Time should be 14.14213562373095101066, not %.20f", time)
	}
	if !journeyRangeData.isOnRange(s) {
		t.Errorf("Step should be on range")
	}
}

func TestGetDistanceBetweenPositions(t *testing.T) {
	initJourneyData()
	s := &FleetJourneyStep{
		MapPosXStart: 20,
		MapPosYStart: 82,
		MapPosXFinal: 25,
		MapPosYFinal: 84,
	}
	if distance := s.getDistanceBetweenPositions(); distance != 5.38516480713450373941 {
		t.Errorf("Distance should be 5.38516480713450373941, not %.20f", distance)
	}
	if s.getType() != JourneyStepTypePositionToPosition {
		t.Errorf("Step is position to position")
	}
	if time := journeyTimeData.TravelTime.getTimeForStep(s); time != 10.77032961426900747881 {
		t.Errorf("Time should be 10.77032961426900747881, not %.20f", time)
	}
	if !journeyRangeData.isOnRange(s) {
		t.Errorf("Step should be on range")
	}
}

func getFleetMock(id uint16, player *Player) *Fleet {
	return &Fleet{
		Id: id,
		Player: player,
		PlayerId: player.Id,
	}
}

func getJourneyMock(id uint16) *FleetJourney {
	journey := &FleetJourney{
		Id: id,
	}
	step1 := getJourneyStepMock(1, journey)
	step2 := getJourneyStepMock(2, journey)
	step3 := getJourneyStepMock(3, journey)

	step1.NextStep = step2
	step1.NextStepId = step2.Id
	step2.NextStep = step3
	step2.NextStepId = step3.Id

	journey.CurrentStep = step1
	journey.CurrentStepId = step1.Id
	journey.Steps = []*FleetJourneyStep{
		step1,
		step2,
		step3,
	}
	return journey
}

func getJourneyStepMock(id uint16, journey *FleetJourney) *FleetJourneyStep {
	return &FleetJourneyStep{
		Id: id,
		Journey: journey,
		JourneyId: journey.Id,
		TimeStart: time.Now().Local().Add(time.Minute * time.Duration((id - 1))),
		TimeArrival: time.Now().Local().Add(time.Minute * time.Duration(id)),
		Order: FleetOrderPass,
		StepNumber: uint32(id),
	}
}