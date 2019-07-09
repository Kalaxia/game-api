package api

import(
	"testing"
	"time"
)

func TestEndStep( t *testing.T) {
	InitDatabaseMock()

	player := getPlayerMock(getFactionMock())
	fleet := getFleetMock(1, player)

	fleet.Journey = getJourneyMock(1)
	fleet.JourneyId = fleet.Journey.Id

	fleet.Journey.CurrentStep.end()

	if fleet.Journey.CurrentStep.Id == 1 {
		t.Errorf("Current step did not change")
	}
	if fleet.Journey.CurrentStep.Id == 2 {
		t.Errorf("Current step is not the next step")
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