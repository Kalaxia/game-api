package api

import(
	"testing"
)

func TestCreateShipModel(t *testing.T) {
	InitShipConfiguration()

	player := getPlayerMock(getFactionMock())
	shipModel := player.createShipModel(map[string]interface{}{
		"name": "XFB-1",
		"frame": "battle-runner",
		"slots": []interface{} {
			map[string]interface{}{
				"position": float64(1),
				"module": "laser-turret-meirrion",
			},
			map[string]interface{}{
				"position": float64(2),
				"module": "laser-turret-meirrion",
			},
			map[string]interface{}{
				"position": float64(3),
				"module": "ion-reactor",
			},
		},
	})
	if name := shipModel.Name; name != "XFB-1" {
		t.Errorf("Ship model name should be XFB-1, got %s", name)
	}
	if slug := shipModel.Frame.Slug; slug != "battle-runner" {
		t.Errorf("Ship frame slug should be a battle-runner, got a %s", slug)
	}
	if length := len(shipModel.Slots); length != 3 {
		t.Errorf("Ship model should have 3 slots, got %d", length)
	}
	if armor := shipModel.Stats["armor"]; armor != 60 {
		t.Errorf("Ship model armor should equal 60, got %d", armor)
	}
	if speed := shipModel.Stats[ShipStatSpeed]; speed != 225 {
		t.Errorf("Ship model speed should equal 115, got %d", speed)
	}
	if category := shipModel.Type; category != ShipTypeFighter {
		t.Errorf("Ship model type should be fighter, got %s", category)
	}
}

func TestCreateShipModelWithoutPropulsor(t *testing.T) {
	InitShipConfiguration()
	defer func() {
        if r := recover(); r != nil {
			exception := r.(*HttpException)

			if exception.Message != "ships.missing_propulsion" {
				t.Errorf("Exception message should tell about the missing propulsor, not %s", exception.Message)
			}
        } else {
			t.Errorf("Ship model creation should have thrown an exception")
		}
    }()

	player := getPlayerMock(getFactionMock())
	player.createShipModel(map[string]interface{}{
		"name": "XFB-1",
		"frame": "battle-runner",
		"slots": []interface{} {
			map[string]interface{}{
				"position": float64(1),
				"module": "laser-turret-meirrion",
			},
			map[string]interface{}{
				"position": float64(2),
				"module": "laser-turret-meirrion",
			},
			map[string]interface{}{
				"position": float64(3),
				"module": nil,
			},
		},
	})
}

func TestExtractSlotsData(t *testing.T) {
	InitShipConfiguration()
	frame := getShipFrameMock()

	slots := frame.extractSlotsData(map[string]interface{}{
		"slots": []interface{} {
			map[string]interface{}{
				"position": float64(1),
				"module": "mark-i-reactor",
			},
			map[string]interface{}{
				"position": float64(2),
				"module": "light-container",
			},
			map[string]interface{}{
				"position": float64(3),
				"module": "laser-turret-meirrion",
			},
		},
	})
	if length := len(slots); length != 3 {
		t.Errorf("There should be 3 slots, got %d", length)
	}
	if position := slots[1].Position; position != 2 {
		t.Errorf("The second slot should be in second position, got %d", position)
	}
	if slug := slots[0].ModuleSlug; slug != "mark-i-reactor" {
		t.Errorf("The first slot module should be mark-i-reactor, got %s", slug)
	}
	if shape := slots[1].Module.Shape; shape != "square" {
		t.Errorf("The second slot module shape should be a square, got a %s", shape)
	}
	if size := slots[2].Module.Size; size != "small" {
		t.Errorf("The third slot module size should be small, not %s", size)
	}
}

func TestExtractSlotData(t *testing.T) {
	InitShipConfiguration()
	frame := getShipFrameMock()

	slot := frame.extractSlot(map[string]interface{}{
		"position": float64(2),
		"module": "light-container",
	})

	if position := slot.Position; position != 2 {
		t.Errorf("Slot position should be %d", position)
	}
	if slug := slot.ModuleSlug; slug != "light-container" {
		t.Errorf("Module slug should be light container, not %s", slug)
	}
	if shape := slot.Module.Shape; shape != "square" {
		t.Errorf("Module shape should be a square, not a %s", shape)
	}
	if size := slot.Module.Size; size != "medium" {
		t.Errorf("Module size should be medium, not %s", size)
	}
}

func TestAddFrameStats(t *testing.T) {
	frame := getShipFrameMock()
	stats := map[string]uint16{
		ShipStatArmor: 20,
		ShipStatShield: 75,
	}
	frame.addStats(stats)

	if armor := stats[ShipStatArmor]; armor != 60 {
		t.Errorf("Ship frame armor should equal 60, got %d", armor)
	}
	if power := stats[ShipStatShield]; power != 75 {
		t.Errorf("Ship frame power should equal 75, got %d", power)
	}
	if speed := stats[ShipStatSpeed]; speed != 50 {
		t.Errorf("Ship frame speed should equal 50, got %d", speed)
	}
}

func TestAddModuleScores(t *testing.T) {
	module := getShipModuleMock()
	scores := map[string]uint8{
		"bomber": 30,
		"fighter": 20,
	}
	module.addScores(scores)

	if bomberScore := scores["bomber"]; bomberScore != 30 {
		t.Errorf("Bomber score should equal 30, got %d", bomberScore)
	}
	if freighterScore := scores["freighter"]; freighterScore != 20 {
		t.Errorf("Freighter score should equal 30, got %d", freighterScore)
	}
	if fighterScore := scores["fighter"]; fighterScore != 50 {
		t.Errorf("Fighter score should equal 30, got %d", fighterScore)
	}
}

func TestAddModuleStats(t *testing.T) {
	module := getShipModuleMock()
	stats := map[string]uint16{
		ShipStatArmor: 20,
		ShipStatShield: 20,
	}

	module.addStats(stats)

	if armor := stats[ShipStatArmor]; armor != 20 {
		t.Errorf("Armor should equal 20, not %d", armor)
	}
	if power := stats[ShipStatShield]; power != 70 {
		t.Errorf("Power should equal 70, not %d", power)
	}
	if cargo := stats[ShipStatCargo]; cargo != 500 {
		t.Errorf("Cargo should equal 500, not %d", cargo)
	}
}

func TestIsValidSlot(t *testing.T) {
	frame := getShipFrameMock()

	if frame.isValidSlot(ShipSlot{
		Position: 1,
		Module: &ShipModule{
			Shape: "square",
			Size: "large",
		},
	}) {
		t.Errorf("Ship slot should not be valid")
	}
}

func TestGetShipModelType(t *testing.T) {
	if category := getShipModelType(map[string]uint8{
		ShipTypeFighter: 40,
		ShipTypeBomber: 50,
		ShipTypeFreighter: 20,
	}); category != ShipTypeBomber {
		t.Errorf("Ship type should be bomber, got %s", category)
	}
}

func getShipFrameMock() *ShipFrame {
	return &ShipFrame{
		Slug: "light-hunter",
		Slots: []ShipSlotPlan{
			ShipSlotPlan{
				Position: 1,
				Shape: "triangle",
				Size: "small",
			},
			ShipSlotPlan{
				Position: 2,
				Shape: "square",
				Size: "medium",
			},
			ShipSlotPlan{
				Position: 3,
				Shape: "circle",
				Size: "small",
			},
		},
		Stats: map[string]uint16{
			"armor": 40,
			"speed": 50,
		},
		Price: []Price{
			Price{
				Type: "credits",
				Amount: 500,
			},
			Price{
				Type: "points",
				Amount: 10,
			},
		},
	}
}

func getShipModuleMock() *ShipModule {
	return &ShipModule{
		Shape: "square",
		Size: "small",
		ShipStats: map[string]uint16{
			ShipStatCargo: 500,
			ShipStatShield: 50,
		},
		Scores: map[string]uint8{
			ShipTypeFreighter: 20,
			ShipTypeFighter: 30,
		},
	}
}