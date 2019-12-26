package api

func (f *Fleet) conquerPlanet(p *Planet) bool {
	if isVictorious := f.attack(p); !isVictorious {
		return false
	}
	p.System.checkTerritories()
	p.notifyConquest(f)
	p.changeOwner(f.Player)
	return true
}

func (p *Planet) notifyConquest(f *Fleet) {
	f.Player.notify(
		NotificationTypeMilitary,
		"notifications.military.planet_conquest",
		map[string]interface{}{
			"planet": p,
			"opponent": p.Player,
		},
	)
	if p.Player != nil {
		p.Player.notify(
			NotificationTypeMilitary,
			"notifications.military.planet_conquerred",
			map[string]interface{}{
				"planet": p,
				"opponent": f.Player,
			},
		)
	}
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
	return p.Player.Faction.Id == f.Player.Faction.Id && f.Player.Faction.Id != fleet.Player.Faction.Id
}