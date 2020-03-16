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
	planet := getPlayerPlanetMock(getPlayerMock(getFactionMock()))
	travellingFleet := &Fleet{}
	orbitingFleet := &Fleet{ Place: &Place{ PlanetId: 1, Planet: planet } }
	elsewhereFleet := &Fleet{ Place: &Place{ PlanetId: 50, Planet: &Planet{ Id: 50 } }}

	if travellingFleet.isOnPlanet(nil) {
		t.Errorf("Travelling fleet is not on planet")
	}
	if !orbitingFleet.isOnPlanet(nil) {
		t.Errorf("Orbiting fleet is on planet")
	}
	if !orbitingFleet.isOnPlanet(planet) {
		t.Errorf("Orbiting fleet is on planet")
	}
	if elsewhereFleet.isOnPlanet(planet) {
		t.Errorf("This fleet is not on this planet")
	}
}

func TestValidateStep(t *testing.T) {
	fleet := &Fleet{}
	step := &FleetJourneyStep{
		Order: FleetOrderPass,
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("journey step should be valid")
		}
	}()
	step.validate(fleet)
}

func TestValidateStepWithInvalidOrder(t *testing.T) {
	fleet := &Fleet{}
	step := &FleetJourneyStep{
		Order: "unexisting-order",
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("journey step should not be valid")
		}
	}()
	step.validate(fleet)
}

func TestValidateConquestStep(t *testing.T) {
	fleet := &Fleet{ PlayerId: 1 }
	step := &FleetJourneyStep{
		Order: FleetOrderConquer,
		EndPlace: &Place{ Planet: &Planet{ PlayerId: 2 }},
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("journey step should be valid")
		}
	}()
	step.validate(fleet)
}

func TestValidateConquestStepWithSamePlayer(t *testing.T) {
	fleet := &Fleet{ PlayerId: 1 }
	step := &FleetJourneyStep{
		Order: FleetOrderConquer,
		EndPlace: &Place{ Planet: &Planet{ PlayerId: 1 }},
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("journey step should not be valid")
		}
	}()
	step.validate(fleet)
}

func TestValidateColonizationStep(t *testing.T) {
	fleet := &Fleet{}
	step := &FleetJourneyStep{
		Order: FleetOrderColonize,
		EndPlace: &Place{ Planet: &Planet{ Population: 0 }},
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("journey step should be valid")
		}
	}()
	step.validate(fleet)
}

func TestValidateColonizationStepWithInhabitedPlanet(t *testing.T) {
	fleet := &Fleet{}
	step := &FleetJourneyStep{
		Order: FleetOrderColonize,
		EndPlace: &Place{ Planet: &Planet{ Population: 1000 }},
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("journey step should not be valid")
		}
	}()
	step.validate(fleet)
}

func TestValidateDeliveryStep(t *testing.T) {
	fleet := &Fleet{}
	step := &FleetJourneyStep{
		Order: FleetOrderDeliver,
		Data: map[string]interface{}{
			"resources": []interface{}{
				map[string]interface{}{
					"resource": "cristal",
					"quantity": 1000,
				},
			},
		},
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("journey step should be valid")
			t.Errorf(r.(error).Error())
		}
	}()
	step.validate(fleet)
}

func TestValidateDeliveryMissions(t *testing.T) {
	fleet := &Fleet{ Journey: getJourneyMock(1) }
	fleet.Cargo = map[string]uint16{
		"cristal": 2000,
		"red-ore": 1000,
	}
	fleet.Journey.Steps = append(fleet.Journey.Steps, &FleetJourneyStep{
		Order: FleetOrderDeliver,
		Data: map[string]interface{}{
			"resources": []interface{}{
				map[string]interface{}{
					"resource": "cristal",
					"quantity": float64(1000),
				},
			},
		},
	}, &FleetJourneyStep{
		Order: FleetOrderDeliver,
		Data: map[string]interface{}{
			"resources": []interface{}{
				map[string]interface{}{
					"resource": "cristal",
					"quantity": float64(500),
				},
				map[string]interface{}{
					"resource": "red-ore",
					"quantity": float64(1000),
				},
			},
		},
	})

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Fleet cargo alert should not have been raised")
		}
	}()

	fleet.validateDeliveryMissions(fleet.Journey.Steps)
}

func TestValidateDeliveryMissionsWithInsufficientCargo(t *testing.T) {
	fleet := &Fleet{ Journey: getJourneyMock(1) }
	fleet.Cargo = map[string]uint16{
		"cristal": 1000,
		"red-ore": 1000,
	}
	fleet.Journey.Steps = append(fleet.Journey.Steps, &FleetJourneyStep{
		Order: FleetOrderDeliver,
		Data: map[string]interface{}{
			"resources": []map[string]interface{}{
				map[string]interface{}{
					"resource": "cristal",
					"quantity": float64(1000),
				},
			},
		},
	}, &FleetJourneyStep{
		Order: FleetOrderDeliver,
		Data: map[string]interface{}{
			"resources": []map[string]interface{}{
				map[string]interface{}{
					"resource": "cristal",
					"quantity": float64(500),
				},
				map[string]interface{}{
					"resource": "red-ore",
					"quantity": float64(1000),
				},
			},
		},
	})

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Fleet cargo alert should have been raised")
		}
	}()

	fleet.validateDeliveryMissions(fleet.Journey.Steps)
}

func TestValidateColonizationMissions(t *testing.T) {

}

func TestGetDistanceBetweenPlanets(t *testing.T) {
	initJourneyData()
	s := &FleetJourneyStep{
		StartPlaceId: 1,
		StartPlace: &Place{
			PlanetId: 1,
			Planet: &Planet{
				Id: 1,
				SystemId: 1,
				System: &System {
					X: 91,
					Y: 16,
				},
			},
		},
		EndPlaceId: 2,
		EndPlace: &Place{
			PlanetId: 2,
			Planet: &Planet{
				Id: 2,
				SystemId: 2,
				System: &System {
					X: 75,
					Y: 20,
				},
			},
		},
	}
	if distance := s.getDistanceBetweenPlanets(); distance != 16.49242250247064234259 {
		t.Errorf("Journey step from planet to planet should be 16.49242250247064234259, not %.20f", distance)
	}
	if sType := s.getType(); sType != JourneyStepTypePlanetToPlanet {
		t.Errorf("Step is planet to planet, not %s", sType)
	}
	if time := journeyTimeData.TravelTime.getTimeForStep(s, 0.2); time != 164.92422502470640210959 {
		t.Errorf("Time should be 164.92422502470640210959, not %.20f", time)
	}
	if journeyRangeData.isOnRange(getFleetRangeMock(), s) {
		t.Errorf("Step should not be on range")
	}
}

func TestGetDistanceBetweenOrbitAndPlanet(t *testing.T) {
	initJourneyData()
	s := &FleetJourneyStep{
		StartPlaceId: 1,
		StartPlace: &Place {
			PlanetId: 1,
			Planet: &Planet{},
		},
		EndPlaceId: 2,
		EndPlace: &Place{
			PlanetId: 1,
			Planet: &Planet{},
		},
	}
	if distance := s.getDistanceBetweenOrbitAndPlanet(); distance != 0. {
		t.Errorf("Distance should be 0, not %.20f", distance)
	}
	if s.getType() != JourneyStepTypeSamePlanet {
		t.Errorf("Step is orbit to planet")
	}
	if time := journeyTimeData.TravelTime.getTimeForStep(s, 0.5); time != 4. {
		t.Errorf("Time should be 4., not %.20f", time)
	}
	if !journeyRangeData.isOnRange(getFleetRangeMock(), s) {
		t.Errorf("Step should be on range")
	}
}

func TestGetDistanceInsideSystem(t *testing.T) {
	initJourneyData()
	s := &FleetJourneyStep{
		StartPlace: &Place{
			PlanetId: 1,
			Planet: &Planet{
				Id: 1,
				SystemId: 1,
			},
		},
		EndPlace: &Place{
			PlanetId: 2,
			Planet: &Planet{
				Id: 2,
				SystemId: 1,
			},
		},
	}
	if distance := s.getDistanceInsideSystem(); distance != 0. {
		t.Errorf("Journey step in same system should be 0, not %.20f", distance)
	}
	if sType := s.getType(); sType != JourneyStepTypeSameSystem {
		t.Errorf("Step is same system, not %s", sType)
	}
	if time := journeyTimeData.TravelTime.getTimeForStep(s, 0.4); time != 12.5 {
		t.Errorf("Time should be 12.5, not %.20f", time)
	}
	if !journeyRangeData.isOnRange(getFleetRangeMock(), s) {
		t.Errorf("Step should be on range")
	}
}

func TestGetDistanceBetweenPlanetAndPosition(t *testing.T) {
	initJourneyData()
	s := &FleetJourneyStep{
		StartPlace: &Place{
			Planet: &Planet{
				System: &System {
					X: 4,
					Y: 16,
				},
			},
		},
		EndPlace: &Place{
			Coordinates: &Coordinates{
				X: 8,
				Y: 20,
			},
		},
	}
	if distance := s.getDistanceBetweenPlanetAndPosition(); distance != 5.65685424949238058190 {
		t.Errorf("Distance should be 5.65685424949238058190, not %.20f", distance)
	}
	if s.getType() != JourneyStepTypePlanetToPosition {
		t.Errorf("Step is position to planet")
	}
	if time := journeyTimeData.TravelTime.getTimeForStep(s, 0.8); time != 14.14213562373095101066 {
		t.Errorf("Time should be 14.14213562373095101066, not %.20f", time)
	}
	if !journeyRangeData.isOnRange(getFleetRangeMock(), s) {
		t.Errorf("Step should be on range")
	}
}

func TestGetDistanceBetweenPositionAndPlanet(t *testing.T) {
	initJourneyData()
	s := &FleetJourneyStep{
		StartPlace: &Place{
			Coordinates: &Coordinates{
				X: 30,
				Y: 40,
			},
		},
		EndPlace: &Place{
			Planet: &Planet{
				System: &System{
					X: 35,
					Y: 45,
				},
			},
		},
	}
	if distance := s.getDistanceBetweenPositionAndPlanet(); distance != 7.07106781186547550533 {
		t.Errorf("Distance should be 7.07106781186547550533, not %.20f", distance)
	}
	if sType := s.getType(); sType != JourneyStepTypePositionToPlanet {
		t.Errorf("Step is position to planet, not %s", sType)
	}
	if time := journeyTimeData.TravelTime.getTimeForStep(s, 0.5); time != 28.28427124746190202131 {
		t.Errorf("Time should be 28.28427124746190202131, not %.20f", time)
	}
	if !journeyRangeData.isOnRange(getFleetRangeMock(), s) {
		t.Errorf("Step should be on range")
	}
}

func TestGetDistanceBetweenPositions(t *testing.T) {
	initJourneyData()
	s := &FleetJourneyStep{
		StartPlace: &Place{
			Coordinates: &Coordinates{
				X: 20,
				Y: 82,
			},
		},
		EndPlace: &Place{
			Coordinates: &Coordinates{
				X: 25,
				Y: 84,
			},
		},
	}
	if distance := s.getDistanceBetweenPositions(); distance != 5.38516480713450373941 {
		t.Errorf("Distance should be 5.38516480713450373941, not %.20f", distance)
	}
	if s.getType() != JourneyStepTypePositionToPosition {
		t.Errorf("Step is position to position")
	}
	if time := journeyTimeData.TravelTime.getTimeForStep(s, 0.4); time != 26.92582403567251603249 {
		t.Errorf("Time should be 26.92582403567251603249, not %.20f", time)
	}
	if !journeyRangeData.isOnRange(getFleetRangeMock(), s) {
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

func getFleetRangeMock() map[string]float64 {
	return map[string]float64{
		JourneyStepTypeSamePlanet: 2.,
		JourneyStepTypeSameSystem: 10.,
		JourneyStepTypePlanetToPlanet: 10.,
		JourneyStepTypePlanetToPosition: 10.,
		JourneyStepTypePositionToPlanet: 10.,
		JourneyStepTypePositionToPosition: 10.,
	}
}