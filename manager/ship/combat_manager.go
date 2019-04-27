package shipManager

import(
	"kalaxia-game-api/database"
	"kalaxia-game-api/exception"
	"kalaxia-game-api/model"
	"math"
	"math/rand"
	"time"
)

func GetCombatReport(id uint16) *model.FleetCombat {
	report := &model.FleetCombat{}

	if err := database.
		Connection.
		Model(report).
		Column("Attacker", "Attacker.Player", "Attacker.Player.Faction", "Defender", "Defender.Player", "Defender.Player.Faction").
		Where("fleet_combat.id = ?", id).
		Select(); err != nil {
		panic(exception.NewHttpException(404, "Report not found", err))
	}
	return report
}

func GetCombatReports(player *model.Player) []model.FleetCombat {
	reports := make([]model.FleetCombat, 0)

	if err := database.
		Connection.
		Model(&reports).
		Column("Attacker", "Attacker.Player", "Attacker.Player.Faction", "Defender", "Defender.Player", "Defender.Player.Faction").
		Where("attacker__player.id = ?", player.Id).
		WhereOr("defender__player.id = ?", player.Id).
		Select(); err != nil {
		panic(exception.NewHttpException(500, "Could not retrieve combat reports", err))
	}
	return reports
}

func Engage(attacker *model.Fleet, defender *model.Fleet) *model.FleetCombat {
	attackerShips := GetFleetShips(attacker)
	defenderShips := GetFleetShips(defender)
	destroyedShips := make([]uint32, 0)

	nbAttackerShips := len(attackerShips)
	nbDefenderShips := len(defenderShips)

	combat := &model.FleetCombat {
		Attacker: attacker,
		AttackerId: attacker.Id,
		Defender: defender,
		DefenderId: defender.Id,

		AttackerShips: formatCombatShips(attackerShips),
		DefenderShips: formatCombatShips(defenderShips),

		AttackerLosses: make(map[string]uint16, 0),
		DefenderLosses: make(map[string]uint16, 0),

		ShipModels: make(map[uint][]*model.ShipSlot, 0),

		BeginAt: time.Now(),
	}

	for nbAttackerShips > 0 && nbDefenderShips > 0 {
		attackerShips, defenderShips, destroyedShips = fightRound(combat, attackerShips, defenderShips, destroyedShips)

		nbAttackerShips = len(attackerShips)
		nbDefenderShips = len(defenderShips)
	}

	combat.IsVictory = nbAttackerShips > 0

	if err := database.Connection.Insert(combat); err != nil {
		panic(exception.NewException("Could not create combat report", err))
	}
	RemoveShipsByIds(destroyedShips)

	return combat
}

func fightRound(combat *model.FleetCombat, attackerShips []model.Ship, defenderShips []model.Ship, destroyedShips []uint32) ([]model.Ship, []model.Ship, []uint32) {
	for _, ship := range defenderShips {
		if _, ok := combat.ShipModels[ship.Model.Id]; !ok {
			combat.ShipModels[ship.Model.Id] = gatherShipModelData(ship.Model)
		}
		shipData := combat.ShipModels[ship.Model.Id]
		index, target := pickRandomTarget(attackerShips)

		openFire(shipData, target)

		if float64(target.Damage) >= math.Ceil(float64(target.Model.Stats["armor"]) / 10) {
			attackerShips = append(attackerShips[:index], attackerShips[index+1:]...)
			destroyedShips = append(destroyedShips, target.Id)
			addLoss(combat, "attacker", target)

			if (len(attackerShips) == 0) {
				break
			}
		}
	}
	for _, ship := range attackerShips {
		if _, ok := combat.ShipModels[ship.Model.Id]; !ok {
			combat.ShipModels[ship.Model.Id] = gatherShipModelData(ship.Model)
		}
		shipData := combat.ShipModels[ship.Model.Id]
		index, target := pickRandomTarget(defenderShips)

		openFire(shipData, target)

		if float64(target.Damage) >= math.Ceil(float64(target.Model.Stats["armor"]) / 10) {
			defenderShips = append(defenderShips[:index], defenderShips[index+1:]...)
			destroyedShips = append(destroyedShips, target.Id)
			addLoss(combat, "defender", target)

			if (len(defenderShips) == 0) {
				break
			}
		}
	}
	return attackerShips, defenderShips, destroyedShips
}

func pickRandomTarget(ships []model.Ship) (int, *model.Ship) {
	index := rand.Intn(len(ships))
	return index, &ships[index]
}

func openFire(attackerSlots []*model.ShipSlot, target *model.Ship) {
	armor := int8(target.Model.Stats["armor"])

	for _, slot := range attackerSlots {
		if slot.Module.Type != "weapon" {
			continue
		}
		for i := 0; uint16(i) < slot.Module.Stats["nb_shots"]; i++ {
			if uint16(rand.Intn(100)) > slot.Module.Stats["precision"] {
				continue
			}
			armor -= int8(slot.Module.Stats["damage"])
			if armor < 0 {
				target.Damage = uint8(math.Abs(float64(armor)))
			}
		}
	}
}

func gatherShipModelData(shipModel *model.ShipModel) []*model.ShipSlot {
	shipSlots := make([]*model.ShipSlot, 0)
	if err := database.Connection.Model(&shipSlots).Where("model_id = ?", shipModel.Id).Select(); err != nil {
		panic(exception.NewException("Could not retrieve ship model slots", err))
	}
	for _, slot := range shipSlots {
		module := modulesData[slot.ModuleSlug]
		slot.Module = &module
	}
	return shipSlots
}

func formatCombatShips(ships []model.Ship) map[string]uint16 {
	formatted := make(map[string]uint16, 0)

	for _, ship := range ships {
		if _, ok := formatted[ship.Model.Type]; !ok {
			formatted[ship.Model.Type] = 0
		}
		formatted[ship.Model.Type]++
	}
	return formatted
}

func addLoss(combat *model.FleetCombat, side string, ship *model.Ship) {
	if side == "defender" {
		if _, ok := combat.DefenderLosses[ship.Model.Type]; !ok {
			combat.DefenderLosses[ship.Model.Type] = 0
		}
		combat.DefenderLosses[ship.Model.Type]++
	} else {
		if _, ok := combat.AttackerLosses[ship.Model.Type]; !ok {
			combat.AttackerLosses[ship.Model.Type] = 0
		}
		combat.AttackerLosses[ship.Model.Type]++
	}
}