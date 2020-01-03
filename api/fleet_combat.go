package api

import(
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

const(
	CombatActionsPerTurn = 10
	CombatActionTypeFire = "fire"
)

type(
	FleetCombat struct {
		tableName struct{} `pg:"fleet__combats"`

		Id uint16 `json:"id"`
		Attacker *Fleet `json:"attacker"`
		AttackerId uint16 `json:"-"`
		Defender *Fleet `json:"defender"`
		DefenderId uint16 `json:"-"`
		Place *Place `json:"place"`
		PlaceId uint32 `json:"-"`
		IsVictory bool `json:"is_victory" pg:",notnull,use_zero"`

		Rounds []*FleetCombatRound `json:"rounds" pg:"fk:combat_id"`

		AttackerShips map[string]uint16 `json:"attacker_ships" pg:",notnull,use_zero"`
		DefenderShips map[string]uint16 `json:"defender_ships" pg:",notnull,use_zero"`

		AttackerLosses map[string]uint16 `json:"attacker_losses" pg:",notnull,use_zero"`
		DefenderLosses map[string]uint16 `json:"defender_losses" pg:",notnull,use_zero"`

		BeginAt time.Time `json:"begin_at"`
		EndAt time.Time `json:"end_at"`

		ShipModels map[uint][]*ShipSlot `json:"-" pg:"-"`
	}

	FleetCombatRound struct{
		tableName struct{} `pg:"fleet__combat_rounds"`

		Id uint32 `json:"id"`
		Combat *FleetCombat `json:"combat"`
		CombatId uint16 `json:"-"`
		Squadrons []*FleetCombatSquadron `json:"squadrons" pg:"fk:round_id"`
	}

	FleetCombatSquadron struct{
		tableName struct{} `pg:"fleet__combat_squadrons"`

		Id uint64 `json:"id"`
		FleetId uint16 `json:"-"`
		Fleet *Fleet `json:"fleet"`
		Initiative uint16 `json:"-"`
		ShipModelId uint `json:"-"`
		ShipModel *ShipModel `json:"ship_model"`
		Squadron *FleetSquadron `json:"-" pg:"-"`
		Action *FleetSquadronAction `json:"action" pg:"fk:squadron_id"`
		Round *FleetCombatRound `json:"round"`
		RoundId uint32 `json:"-"`
		Quantity uint8 `json:"quantity"`
		Position *FleetGridPosition `json:"position" pg:"type:jsonb"`
	}

	FleetSquadronAction struct{
		tableName struct{} `pg:"fleet__combat_squadron_actions"`

		Id uint64 `json:"id"`
		Squadron *FleetCombatSquadron `json:"squadron"`
		SquadronId uint64 `json:"-"`
		Type string `json:"type"`
		Target *FleetCombatSquadron `json:"target"`
		TargetId uint64 `json:"-"`
		Loss uint8 `json:"loss"`
	}
)

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

func GetCombatRound(w http.ResponseWriter, r *http.Request) {
	player := context.Get(r, "player").(*Player)
    id, _ := strconv.ParseUint(mux.Vars(r)["combatId"], 10, 32)
	roundId, _ := strconv.ParseUint(mux.Vars(r)["roundId"], 10, 32)

	report := getCombatReport(uint16(id))

	if report.Attacker.Player.Id != player.Id && report.Defender.Player.Id != player.Id {
		panic(NewHttpException(403, "You do not own this combat report", nil))
	}
	SendJsonResponse(w, 200, report.getRound(uint32(roundId)))
}

func getCombatReport(id uint16) *FleetCombat {
	report := &FleetCombat{}
	if err := Database.
		Model(report).
		Relation("Rounds").
		Relation("Place.Planet.Player.Faction").
		Relation("Place.Planet.System").
		Relation("Attacker.Player.Faction").
		Relation("Defender.Player.Faction").
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
		Relation("Rounds").
		Relation("Attacker.Player.Faction").
		Relation("Defender.Player.Faction").
		Where("attacker__player.id = ?", p.Id).
		WhereOr("defender__player.id = ?", p.Id).
		Order("end_at DESC").
		Select(); err != nil {
		panic(NewHttpException(500, "Could not retrieve combat reports", err))
	}
	return reports
}

func (attacker *Fleet) engage(defender *Fleet) *FleetCombat {
	attacker.Squadrons = attacker.getSquadrons()
	defender.Squadrons = defender.getSquadrons()

	if !defender.hasSquadrons() {
		return nil
	}
	combat := attacker.newCombat(defender)
	wsMessage := &WsMessage{
		Action: "combat_start",
		Data: combat,
	}
	WsHub.sendTo(attacker.Player, wsMessage)
	WsHub.sendTo(defender.Player, wsMessage)
	for attacker.hasSquadrons() && defender.hasSquadrons() {
		combat.fightRound(attacker, defender)
	}
	combat.IsVictory = attacker.hasSquadrons()
	combat.AttackerLosses = attacker.calculateLosses(combat.AttackerShips)
	combat.DefenderLosses = defender.calculateLosses(combat.DefenderShips)
	combat.EndAt = time.Now()

	combat.update()
	combat.notifyEnding()
	combat.processConsequences()

	return combat
}

func (c *FleetCombat) notifyEnding() {
	// We reload the report in a serializable format: the object used to process the combat has too many circular dependencies
	report := getCombatReport(c.Id)
	wsMessage := &WsMessage{
		Action: "combat_end",
		Data: report,
	}
	WsHub.sendTo(report.Attacker.Player, wsMessage)
	WsHub.sendTo(report.Defender.Player, wsMessage)

	if report.IsVictory {
		report.Attacker.notifyCombatEnding(report, report.Defender, "victory")
		report.Defender.notifyCombatEnding(report, report.Attacker, "defeat")
	} else {
		report.Defender.notifyCombatEnding(report, report.Attacker, "victory")
		report.Attacker.notifyCombatEnding(report, report.Defender, "defeat")
	}
}

func (c *FleetCombat) processConsequences() {
	attackerFaction := getFaction(c.Attacker.Player.FactionId)
	defenderFaction := getFaction(c.Defender.Player.FactionId)

	if !attackerFaction.isInWarWith(defenderFaction) {
		attackerFaction.createCasusBelli(defenderFaction, c.Attacker.Player, CasusBelliTypeCombat, map[string]interface{}{
			"combat": getCombatReport(c.Id),
		})
	}
}

func (f *Fleet) newCombat(opponent *Fleet) *FleetCombat {
	combat := &FleetCombat {
		Attacker: f,
		AttackerId: f.Id,
		Defender: opponent,
		DefenderId: opponent.Id,

		Place: opponent.Place,
		PlaceId: opponent.PlaceId,

		AttackerShips: f.formatCombatShips(),
		DefenderShips: opponent.formatCombatShips(),

		AttackerLosses: make(map[string]uint16, 0),
		DefenderLosses: make(map[string]uint16, 0),

		ShipModels: make(map[uint][]*ShipSlot, 0),

		BeginAt: time.Now(),
	}
	if err := Database.Insert(combat); err != nil {
		panic(NewException("Could not create combat report", err))
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

func (c *FleetCombat) fightRound(attacker *Fleet, defender *Fleet) {
	round := c.newRound(attacker, defender)
	round.initSquadrons(attacker)
	round.initSquadrons(defender)
	round.processInitiative()

	nbActions := len(round.Squadrons)
	if nbActions > CombatActionsPerTurn {
		nbActions = CombatActionsPerTurn
	}

	for i := 0; i < nbActions; i++ {
		round.Squadrons[i].act()
	}
	// We reload the round data to avoid circular dependencies while serializing the data
	wsMessage := &WsMessage{
		Action: "combat_round",
		Data: struct{
			Combat *FleetCombat `json:"combat"`
			Round *FleetCombatRound `json:"round"`
		}{
			Combat: getCombatReport(c.Id),
			Round: c.getRound(round.Id),
		},
	}
	WsHub.sendTo(attacker.Player, wsMessage)
	WsHub.sendTo(defender.Player, wsMessage)
	time.Sleep(time.Duration(1000 * nbActions) * time.Millisecond)
}

func (c *FleetCombat) newRound(attacker *Fleet, defender *Fleet) *FleetCombatRound {
	round := &FleetCombatRound{
		Combat: c,
		CombatId: c.Id,
	}
	if err := Database.Insert(round); err != nil {
		panic(NewException("Could not create fleet combat round", err))
	}
	return round
}

func (c *FleetCombat) getRound(roundId uint32) *FleetCombatRound {
	round := &FleetCombatRound{}
	if err := Database.Model(round).Relation("Squadrons.Action.Target").Relation("Squadrons.Fleet").Relation("Squadrons.ShipModel").Where("combat_id = ?", c.Id).Where("id = ?", roundId).Select(); err != nil {
		panic(NewHttpException(404, "Could not retrieve combat round", err))
	}
	return round
}

func (r *FleetCombatRound) initSquadrons(f *Fleet) {
	for _, s := range f.Squadrons {
		s.Fleet = f
		s.FleetId = f.Id
		r.Squadrons = append(r.Squadrons, s.createCombatCopy(r))
	}
}

func (s *FleetSquadron) createCombatCopy(r *FleetCombatRound) *FleetCombatSquadron {
	squadron := &FleetCombatSquadron{
		Fleet: s.Fleet,
		FleetId: s.FleetId,
		Initiative: s.CombatInitiative,
		Quantity: s.Quantity,
		Round: r,
		RoundId: r.Id,
		ShipModel: s.ShipModel,
		ShipModelId: s.ShipModelId,
		Squadron: s,
	}
	if s.CombatPosition != nil {
		squadron.Position = s.CombatPosition
	} else {
		squadron.Position = s.Position
	}
	if err := Database.Insert(squadron); err != nil {
		panic(NewException("Could not create combat squadron", err))
	}
	return squadron
}

func (r *FleetCombatRound) processInitiative() {
	for _, s := range r.Squadrons {
		s.Initiative += s.calculateInitiative()
	}
	sort.Slice(r.Squadrons, func(i, j int) bool {
		return r.Squadrons[i].Initiative > r.Squadrons[j].Initiative
	})
}

func (s *FleetCombatSquadron) calculateInitiative() uint16 {
	speed := s.ShipModel.Stats[ShipStatSpeed]

	return speed + uint16(rand.Intn(100))
}

func (s *FleetCombatSquadron) act() {
	action := &FleetSquadronAction{
		Type: CombatActionTypeFire,
		Squadron: s,
		SquadronId: s.Id,
	}
	if hasTarget := action.pickTarget(); !hasTarget {
		return
	}
	action.openFire()

	if err := Database.Insert(action); err != nil {
		panic(NewException("Could not create fleet squadron action", err))
	}
	s.Squadron.CombatInitiative = 0
}

func (action *FleetSquadronAction) pickTarget() bool {
	possibleTargets := make([]*FleetCombatSquadron, 0)

	for _, squadron := range action.Squadron.Round.Squadrons {
		if squadron.FleetId != action.Squadron.FleetId && squadron.Quantity > 0 {
			possibleTargets = append(possibleTargets, squadron)
		}
	}
	if len(possibleTargets) > 1 {
		action.Target = possibleTargets[rand.Intn(len(possibleTargets) - 1)]
	} else if len(possibleTargets) == 1 {
		action.Target = possibleTargets[0]
	} else {
		return false
	}
	action.TargetId = action.Target.Id
	return true
}

func (action *FleetSquadronAction) openFire() {
	damage := uint16(0)

	for _, slot := range action.Squadron.getSlots() {
		if slot.canShoot() {
			damage += slot.shoot(action.Target)
		}
	}
	damage = damage * uint16(action.Squadron.Quantity)
	action.Loss = action.Target.receiveDamage(damage)
	action.Target.Quantity -= action.Loss

	if action.Target.Quantity == 0 {
		action.Target.delete()
	} else {
		action.Target.Squadron.Quantity = action.Target.Quantity
	}
}

func (s *FleetCombatSquadron) getSlots() []*ShipSlot {
	if _, ok := s.Round.Combat.ShipModels[s.ShipModel.Id]; !ok {
		s.Round.Combat.ShipModels[s.ShipModel.Id] = s.ShipModel.getSlots()
	}
	return s.Round.Combat.ShipModels[s.ShipModel.Id]
}

func (s *ShipSlot) shoot(target *FleetCombatSquadron) (damage uint16) {
	for i := 0; uint16(i) < s.Module.Stats["nb_shots"]; i++ {
		if s.doesHit(target) {
			damage += uint16(float64(s.Module.Stats["damage"]) * 0.15)
		}
	}
	return
}

func (s *ShipSlot) canShoot() bool {
	return s.Module != nil && s.Module.Type == "weapon"
}

func (s *ShipSlot) doesHit(target *FleetCombatSquadron) bool {
	return uint16(rand.Intn(100)) <= s.Module.Stats["precision"]
}

func (s *FleetCombatSquadron) receiveDamage(damage uint16) uint8 {
	loss := uint8(0)
	armor := s.ShipModel.Stats["armor"]
	hitPoints := s.ShipModel.getHitPoints()

	for damage > 0 {
		if armor >= damage {
			break
		}
		takenDamage := int32(damage) - int32(armor)
		if takenDamage <= 0 {
			break
		}
		if float64(takenDamage) > hitPoints {
			takenDamage = int32(hitPoints)
			loss++
		}
		if s.Quantity == loss {
			break
		}
		damage -= uint16(takenDamage)
	}
	return loss
}

func (sm *ShipModel) getHitPoints() float64 {
	return math.Ceil(float64(sm.Stats["armor"]) / 10)
}

func (sm *ShipModel) getSlots() []*ShipSlot {
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

func (f *Fleet) calculateLosses(initialShips map[string]uint16) map[string]uint16 {
	losses := make(map[string]uint16, 0)
	currentShips := f.formatCombatShips()

	for model, initialQuantity := range initialShips {
		if _, ok := currentShips[model]; !ok {
			losses[model] = initialQuantity
		} else if loss := initialQuantity - currentShips[model]; loss > 0  {
			losses[model] = loss
		}
	}
	return losses
}

func (f *Fleet) formatCombatShips() map[string]uint16 {
	formatted := make(map[string]uint16)

	for _, squadron := range f.Squadrons {
		if _, ok := formatted[squadron.ShipModel.Type]; !ok {
			formatted[squadron.ShipModel.Type] = 0
		}
		formatted[squadron.ShipModel.Type] += uint16(squadron.Quantity)
	}
	return formatted
}

func (c *FleetCombat) update() {
	if err := Database.Update(c); err != nil {
		panic(NewException("Could not update fleet combat", err))
	}
}

func (s *FleetCombatSquadron) delete() {
	s.Squadron.Fleet.deleteSquadron(s.Squadron)
}