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
	s := &FleetJourneyStep{
		PlanetStart: &Planet{
			SystemId: 1,
			System: &System {
				X: 91,
				Y: 16,
			},
		},
		PlanetFinal: &Planet{
			SystemId: 2,
			System: &System {
				X: 75,
				Y: 20,
			},
		},
	}
	if distance := s.getDistance(); distance != 16.49242250247064234259 {
		t.Errorf("Journey step from planet to planet should be 16.49242250247064234259, not %.20f", distance)
	}
}

func TestGetDistanceBetweenSystemPlanets(t *testing.T) {
	s := &FleetJourneyStep{
		PlanetStart: &Planet{
			SystemId: 1,
		},
		PlanetFinal: &Planet{
			SystemId: 1,
		},
	}
	if distance := s.getDistance(); distance != 0. {
		t.Errorf("Journey step in same system should be 0, not %.20f", distance)
	}
}

func TestGetDistanceBetweenPlanetAndPoint(t *testing.T) {
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
	if distance := s.getDistance(); distance != 5.65685424949238058190 {
		t.Errorf("Distance should be 5.65685424949238058190, not %.20f", distance)
	}
}

func TestGetDistanceBetweenPointAndPlanet(t *testing.T) {
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
	if distance := s.getDistance(); distance != 7.07106781186547550533 {
		t.Errorf("Distance should be 7.07106781186547550533, not %.20f", distance)
	}
}

func TestGetDistanceBetweenPoints(t *testing.T) {
	s := &FleetJourneyStep{
		MapPosXStart: 20,
		MapPosYStart: 82,
		MapPosXFinal: 25,
		MapPosYFinal: 84,
	}
	if distance := s.getDistance(); distance != 5.38516480713450373941 {
		t.Errorf("Distance should be 5.38516480713450373941, not %.20f", distance)
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