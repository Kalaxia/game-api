package api

import(
	"testing"
)

func TestWillEngage(t *testing.T) {
	fleet1 := &Fleet{ Player: &Player{ Faction: &Faction{ Id: 1 }}}
	fleet2 := &Fleet{ Player: &Player{ Faction: &Faction{ Id: 2 }}}
	fleet3 := &Fleet{ Player: &Player{ Faction: &Faction{ Id: 1 }}}
	planet1 := &Planet{ Player: &Player{ Faction: &Faction{ Id: 1 }}}
	planet2 := &Planet{ Player: &Player{ Faction: &Faction{ Id: 3 }}}

	if !fleet1.willEngage(fleet2, planet1) {
		t.Errorf("Allied fleet must engage enemy fleet")
	}
	if fleet3.willEngage(fleet2, planet2) {
		t.Errorf("Neutral fleet must not engage attacking fleet")
	}
	if fleet3.willEngage(fleet1, planet2) {
		t.Errorf("Allied fleets must not engage each other")
	}
}