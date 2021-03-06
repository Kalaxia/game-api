package api

func (f *Fleet) colonizePlanet(p *Planet, population int16) {
	p.notifyColonization(f)
	p.changeOwner(f.Player)
	p.System.checkTerritories()
	f.transferPassengers(p, -population)
	p.createPlayerRelation(f.Player, 400)
	p.createFactionRelation(f.Player.Faction, 400)
}

func (f *Fleet) conquerPlanet(p *Planet) bool {
	if isVictorious := f.attack(p); !isVictorious {
		return false
	}
	p.checkForCasusBelli(f.Player)
	p.notifyConquest(f)
	p.changeOwner(f.Player)
	p.System.checkTerritories()
	return true
}

func (p *Planet) checkForCasusBelli(attacker *Player) {
	if p.Player == nil || p.Player.FactionId == attacker.FactionId {
		return
	}
	attackerFaction := getFaction(attacker.FactionId)
	defenderFaction := getFaction(p.Player.FactionId)

	if !attackerFaction.isInWarWith(defenderFaction) {
		attackerFaction.createCasusBelli(defenderFaction, attacker, CasusBelliTypeConquest, map[string]interface{}{
			"planet": p,
		})
	}
}

func (p *Planet) notifyConquest(f *Fleet) {
	if p.Player == nil {
		f.Player.notify(
			NotificationTypeMilitary,
			"notifications.military.planet_conquest",
			map[string]interface{}{
				"planet_id": p.Id,
				"planet_name": p.Name,
			},
		)
		return
	}
	f.Player.notify(
		NotificationTypeMilitary,
		"notifications.military.player_planet_conquest",
		map[string]interface{}{
			"planet_id": p.Id,
			"planet_name": p.Name,
			"opponent_id": p.Player.Id,
			"opponent_pseudo": p.Player.Pseudo,
		},
	)
	p.Player.notify(
		NotificationTypeMilitary,
		"notifications.military.planet_conquerred",
		map[string]interface{}{
			"planet_id": p.Id,
			"planet_name": p.Name,
			"opponent_id": f.Player.Id,
			"opponent_pseudo": f.Player.Pseudo,
		},
	)
}


func (p *Planet) notifyColonization(f *Fleet) {
	f.Player.notify(
		NotificationTypeMilitary,
		"notifications.military.planet_colonization",
		map[string]interface{}{
			"planet_id": p.Id,
			"planet_name": p.Name,
		},
	)
}

func (fleet *Fleet) attack(p *Planet) bool {
	for _, f := range p.getOrbitingFleets() {
		if !f.willEngage(fleet, p) {
			continue
		}
		if combat := fleet.engage(f); combat != nil && !combat.IsVictory {
			return false
		}
	}
	return true
}

func (f *Fleet) willEngage(fleet *Fleet, p *Planet) bool {
	// The fleets will engage on a empty planet or to defend an ally planet (TODO ally relationship)
	return (p.Player == nil || p.Player.Faction.Id == f.Player.Faction.Id) && f.Player.Faction.Id != fleet.Player.Faction.Id
}