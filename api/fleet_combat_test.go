package api

import(
	"reflect"
	"testing"
)

func TestPickRandomTarget(t *testing.T) {
	ships := []Ship{
		Ship{},
		Ship{},
		Ship{},
	}
	index, target := pickRandomTarget(ships)

	if index < 0 || index > 2 {
		t.Errorf("Target index is not between 0 and 2")
	}
	if structName := reflect.TypeOf(target).Elem().Name(); structName != "Ship" {
		t.Errorf("Target is a %s, not a ship", structName)
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

	target := &Ship{}

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

	target := &Ship{}

	if damage := hittingSlot.shoot(target); damage != 40 {
		t.Errorf("Slot did not make 40 damages, made %d", damage)
	}
	if damage := missingSlot.shoot(target); damage != 0 {
		t.Errorf("Missing slot made %d damages", damage)
	}
}

func TestIsDestroyed(t *testing.T) {
	destroyedShip := &Ship{ 
		Damage: uint8(5),
		Model: &ShipModel{ Stats: map[string]uint16{ "armor": uint16(50) }},
	}
	survivingShip := &Ship{ 
		Damage: uint8(3),
		Model: &ShipModel{ Stats: map[string]uint16{ "armor": uint16(50) }},
	}

	if !destroyedShip.isDestroyed() {
		t.Errorf("Destroyed ship is not destroyed")
	}
	if survivingShip.isDestroyed() {
		t.Errorf("Surviving ship is destroyed")
	}
}

func TestDestroyShip(t *testing.T) {
	combat := &FleetCombat{
		AttackerLosses: make(map[string]uint16, 0),
		DefenderLosses: make(map[string]uint16, 0),
	}
	ships := []Ship{
		Ship{ Id: 1, Model: &ShipModel{ Type: "cruiser" } },
		Ship{ Id: 2, Model: &ShipModel{ Type: "fighter" } },
		Ship{ Id: 3, Model: &ShipModel{ Type: "cargo" } },
	}
	destroyedShips := make([]uint32, 0)

	ships, destroyedShips = combat.destroyShip("defender", 1, ships, destroyedShips)

	if len(ships) != 2 {
		t.Errorf("There must be two remaining ships")
	}
	if len(destroyedShips) != 1 {
		t.Errorf("There must be one destroyed ship")
	}
	if len(combat.DefenderLosses) != 1 {
		t.Errorf("There must be one defender loss")
	}
	if v, ok := combat.DefenderLosses["fighter"]; !ok {
		t.Errorf("Defender loss must be a fighter")
		if v != 1 {
			t.Errorf("Defender fighter loss count must equals 1")
		}
	}
	if destroyedShips[0] != 2 {
		t.Errorf("The destroyed ship ID must equals 2")
	}
}