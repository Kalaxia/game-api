package api

import(
	"testing"
)

func TestCreateCombatCopy(t *testing.T) {
	squadron := &FleetSquadron{
		Id: 13,
		Fleet: &Fleet{
			Id: 1,
		},
		FleetId: 1,
		ShipModel: &ShipModel{
			Id: 10,
		},
		ShipModelId: 10,
		Quantity: 5,
		CombatInitiative: 200,
		Position: &FleetGridPosition{
			X: 10,
			Y: 15,
		},
	}
	round := &FleetCombatRound{
		Id: 4,
	}

	copy := squadron.createCombatCopy(round)

	if id := copy.Squadron.Id; id != 13 {
		t.Errorf("Id should equal 13, got %d", id)
	}
	if fleetId := copy.Fleet.Id; fleetId != 1 {
		t.Errorf("Fleet id should equal 1, got %d", fleetId)
	}
	if shipModelId := copy.ShipModel.Id; shipModelId != 10 {
		t.Errorf("Ship model id should equal 10, got %d", shipModelId)
	}
	if quantity := copy.Quantity; quantity != 5 {
		t.Errorf("Quantity should equal 5, got %d", quantity)
	}
	if initiative := copy.Initiative; initiative != 200 {
		t.Errorf("Initiative should equal 200, got %d", initiative)
	}
	if x := copy.Position.X; x != 10 {
		t.Errorf("X position should equal 10, got %d", x)
	}
	if y := copy.Position.Y; y != 15 {
		t.Errorf("Y position should equal 15, got %d", y)
	}

	squadron.CombatPosition = &FleetGridPosition{
		X: 20,
		Y: 30,
	}

	copy = squadron.createCombatCopy(round)

	if x := copy.Position.X; x != 20 {
		t.Errorf("X position should equal 20, got %d", x)
	}
	if y := copy.Position.Y; y != 30 {
		t.Errorf("Y position should equal 30, got %d", y)
	}
}

func TestPickRandomTarget(t *testing.T) {
	action := &FleetSquadronAction{
		Squadron: &FleetCombatSquadron{
			Id: 1,
			FleetId: 1,
			Round: &FleetCombatRound{
				Squadrons: []*FleetCombatSquadron{
					&FleetCombatSquadron{
						Id: 1,
						FleetId: 1,
						Quantity: 3,
					},
					&FleetCombatSquadron{
						Id: 2,
						FleetId: 2,
						Quantity: 4,
					},
					&FleetCombatSquadron{
						Id: 3,
						FleetId: 2,
						Quantity: 0,
					},
				},
			},
		},
	}
	action.pickTarget()

	if action.Target == nil {
		t.Errorf("Action target should be defined")
	}
	if targetId := action.Target.Id; targetId != 2 {
		t.Errorf("Target id should equal 2, got %d", targetId)
	}
	if action.TargetId != action.Target.Id {
		t.Errorf("Target id should be set in action field")
	}
}

func TestCanShoot(t *testing.T) {
	emptySlot := &ShipSlot{}
	weaponSlot := &ShipSlot{ Module: &ShipModule{ Type: "weapon" }}
	shieldSlot := &ShipSlot{ Module: &ShipModule{ Type: "shield" }}

	if emptySlot.canShoot() {
		t.Errorf("Empty slot can't shoot")
	}
	if !weaponSlot.canShoot() {
		t.Errorf("Weapon slot must shoot")
	}
	if shieldSlot.canShoot() {
		t.Errorf("shield slot can't shoot")
	}
}

func TestDoesHit(t *testing.T) {
	missingSlot := &ShipSlot{ Module: &ShipModule{ Stats: map[string]uint16{ "precision": uint16(0) }}}
	hittingSlot := &ShipSlot{ Module: &ShipModule{ Stats: map[string]uint16{ "precision": uint16(100) }}}

	target := &FleetCombatSquadron{}

	if missingSlot.doesHit(target) {
		t.Errorf("Missing slot can't hit target")
	}
	if !hittingSlot.doesHit(target) {
		t.Errorf("Hitting slot must hit target")
	}
}

func TestShoot(t *testing.T) {
	hittingSlot := &ShipSlot{ Module: &ShipModule{ Stats: map[string]uint16{
		"precision": uint16(100),
		"nb_shots": uint16(4),
		"damage": uint16(10),
	}}}
	missingSlot := &ShipSlot{ Module: &ShipModule{ Stats: map[string]uint16{
		"precision": uint16(0),
		"nb_shots": uint16(2),
		"damage": uint16(100),
	}}}

	target := &FleetCombatSquadron{}

	if damage := hittingSlot.shoot(target); damage != 4 {
		t.Errorf("Slot did not make 40 damages, made %d", damage)
	}
	if damage := missingSlot.shoot(target); damage != 0 {
		t.Errorf("Missing slot made %d damages", damage)
	}
}

func TestProcessInitiative(t *testing.T) {
	round := &FleetCombatRound{
		Squadrons: []*FleetCombatSquadron{
			&FleetCombatSquadron{
				Id: 1,
				Initiative: 0,
				ShipModel: &ShipModel{
					Stats: map[string]uint16{
						"speed": 255,
					},
				},
			},
			&FleetCombatSquadron{
				Id: 2,
				Initiative: 0,
				ShipModel: &ShipModel{
					Stats: map[string]uint16{
						"speed": 100,
					},
				},
			},
			&FleetCombatSquadron{
				Id: 3,
				Initiative: 0,
				ShipModel: &ShipModel{
					Stats: map[string]uint16{
						"speed": 450,
					},
				},
			},
			&FleetCombatSquadron{
				Id: 4,
				Initiative: 0,
				ShipModel: &ShipModel{
					Stats: map[string]uint16{
						"speed": 225,
					},
				},
			},
		},
	}

	round.processInitiative()

	if id := round.Squadrons[0].Id; id != 3 {
		t.Errorf("First squadron should be squadron 3, not %d", id)
	}
	if id := round.Squadrons[3].Id; id != 2 {
		t.Errorf("Last squadron should be squadron 2, not %d", id)
	}
}

func TestCalculateInitiative(t *testing.T) {
	squadron := FleetCombatSquadron{
		ShipModel: &ShipModel{
			Stats: map[string]uint16{
				"speed": 225,
			},
		},
	}
	if initiative := squadron.calculateInitiative(); initiative < 225 || initiative > 325 {
		t.Errorf("Initiative should be between 225 and 325, got %d", initiative)
	}
}