package api

import(
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
	"math"
	"math/rand"
	"strconv"
	"time"
)

type FleetCombat struct {
	TableName struct{} `json:"-" sql:"fleet__combats"`

	Id uint16 `json:"id"`
	Attacker *Fleet `json:"attacker"`
	AttackerId uint16 `json:"-"`
	Defender *Fleet `json:"defender"`
	DefenderId uint16 `json:"-"`
	IsVictory bool `json:"is_victory" sql:",notnull"`

	AttackerShips map[string]uint16 `json:"attacker_ships" sql:",notnull"`
	DefenderShips map[string]uint16 `json:"defender_ships" sql:",notnull"`

	AttackerLosses map[string]uint16 `json:"attacker_losses" sql:",notnull"`
	DefenderLosses map[string]uint16 `json:"defender_losses" sql:",notnull"`

	BeginAt time.Time `json:"begin_at"`
	EndAt time.Time `json:"end_at"`

	ShipModels map[uint][]*ShipSlot `json:"-" sql:"-"`
}

func GetCombatReport(w http.ResponseWriter, r *http.Request) {
	player := context.Get(r, "player").(*Player)
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)

	report := getCombatReport(uint16(id))

	if report.Attacker.Player.Id != player.Id && report.Defender.Player.Id != player.Id {
		panic(NewHttpException(403, "You do not own this combat report", nil))
	}
	SendJsonResponse(w, 200, report)
}

func GetCombatReports(w http.ResponseWriter, r *http.Request) {
	player := context.Get(r, "player").(*Player)

	SendJsonResponse(w, 200, player.getCombatReports())
}

func getCombatReport(id uint16) *FleetCombat {
	report := &FleetCombat{}
	if err := Database.
		Model(report).
		Column("Attacker", "Attacker.Player", "Attacker.Player.Faction", "Defender", "Defender.Player", "Defender.Player.Faction").
		Where("fleet_combat.id = ?", id).
		Select(); err != nil {
		panic(NewHttpException(404, "Report not found", err))
	}
	return report
}

func (p *Player) getCombatReports() []FleetCombat {
	reports := make([]FleetCombat, 0)

	if err := Database.
		Model(&reports).
		Column("Attacker", "Attacker.Player", "Attacker.Player.Faction", "Defender", "Defender.Player", "Defender.Player.Faction").
		Where("attacker__player.id = ?", p.Id).
		WhereOr("defender__player.id = ?", p.Id).
		Order("end_at DESC").
		Select(); err != nil {
		panic(NewHttpException(500, "Could not retrieve combat reports", err))
	}
	return reports
}

func (attacker *Fleet) engage(defender *Fleet) *FleetCombat {
	attackerShips := attacker.getShips()
	defenderShips := defender.getShips()

	if len(defenderShips) == 0 {
		return nil
	}
	destroyedShips := make([]uint32, 0)

	nbAttackerShips := len(attackerShips)
	nbDefenderShips := len(defenderShips)

	combat := &FleetCombat {
		Attacker: attacker,
		AttackerId: attacker.Id,
		Defender: defender,
		DefenderId: defender.Id,

		AttackerShips: formatCombatShips(attackerShips),
		DefenderShips: formatCombatShips(defenderShips),

		AttackerLosses: make(map[string]uint16, 0),
		DefenderLosses: make(map[string]uint16, 0),

		ShipModels: make(map[uint][]*ShipSlot, 0),

		BeginAt: time.Now(),
	}

	for nbAttackerShips > 0 && nbDefenderShips > 0 {
		attackerShips, defenderShips, destroyedShips = combat.fightRound(attackerShips, defenderShips, destroyedShips)

		nbAttackerShips = len(attackerShips)
		nbDefenderShips = len(defenderShips)
	}
	combat.IsVictory = nbAttackerShips > 0

	if err := Database.Insert(combat); err != nil {
		panic(NewException("Could not create combat report", err))
	}
	removeShipsByIds(destroyedShips)

	if combat.IsVictory {
		attacker.notifyCombatEnding(combat, defender, "victory")
		defender.notifyCombatEnding(combat, attacker, "defeat")
	} else {
		defender.notifyCombatEnding(combat, attacker, "victory")
		attacker.notifyCombatEnding(combat, defender, "defeat")
	}
	return combat
}

func (f *Fleet) notifyCombatEnding(report *FleetCombat, opponent *Fleet, state string) {
	f.Player.notify(
		NotificationTypeMilitary,
		"notifications.military.fleet_" + state,
		map[string]interface{}{
			"fleet": f,
			"report": report,
			"opponent": opponent,
		},
	)
}

func (c *FleetCombat) fightRound(attackerShips []Ship, defenderShips []Ship, destroyedShips []uint32) ([]Ship, []Ship, []uint32) {
	for _, ship := range defenderShips {
		if _, ok := c.ShipModels[ship.Model.Id]; !ok {
			c.ShipModels[ship.Model.Id] = ship.Model.gatherData()
		}
		shipData := c.ShipModels[ship.Model.Id]
		index, target := pickRandomTarget(attackerShips)

		openFire(shipData, target)

		if float64(target.Damage) >= math.Ceil(float64(target.Model.Stats["armor"]) / 10) {
			attackerShips = append(attackerShips[:index], attackerShips[index+1:]...)
			destroyedShips = append(destroyedShips, target.Id)
			c.addLoss("attacker", target)

			if (len(attackerShips) == 0) {
				break
			}
		}
	}
	for _, ship := range attackerShips {
		if _, ok := c.ShipModels[ship.Model.Id]; !ok {
			c.ShipModels[ship.Model.Id] = ship.Model.gatherData()
		}
		shipData := c.ShipModels[ship.Model.Id]
		index, target := pickRandomTarget(defenderShips)

		openFire(shipData, target)

		if float64(target.Damage) >= math.Ceil(float64(target.Model.Stats["armor"]) / 10) {
			defenderShips = append(defenderShips[:index], defenderShips[index+1:]...)
			destroyedShips = append(destroyedShips, target.Id)
			c.addLoss("defender", target)

			if (len(defenderShips) == 0) {
				break
			}
		}
	}
	return attackerShips, defenderShips, destroyedShips
}

func pickRandomTarget(ships []Ship) (int, *Ship) {
	index := rand.Intn(len(ships))
	return index, &ships[index]
}

func openFire(attackerSlots []*ShipSlot, target *Ship) {
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

func (sm *ShipModel) gatherData() []*ShipSlot {
	shipSlots := make([]*ShipSlot, 0)
	if err := Database.Model(&shipSlots).Where("model_id = ?", sm.Id).Select(); err != nil {
		panic(NewException("Could not retrieve ship model slots", err))
	}
	for _, slot := range shipSlots {
		module := modulesData[slot.ModuleSlug]
		slot.Module = &module
	}
	return shipSlots
}

func formatCombatShips(ships []Ship) map[string]uint16 {
	formatted := make(map[string]uint16)

	for _, ship := range ships {
		if _, ok := formatted[ship.Model.Type]; !ok {
			formatted[ship.Model.Type] = 0
		}
		formatted[ship.Model.Type]++
	}
	return formatted
}

func (c *FleetCombat) addLoss(side string, ship *Ship) {
	if side == "defender" {
		if _, ok := c.DefenderLosses[ship.Model.Type]; !ok {
			c.DefenderLosses[ship.Model.Type] = 0
		}
		c.DefenderLosses[ship.Model.Type]++
	} else {
		if _, ok := c.AttackerLosses[ship.Model.Type]; !ok {
			c.AttackerLosses[ship.Model.Type] = 0
		}
		c.AttackerLosses[ship.Model.Type]++
	}
}