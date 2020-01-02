package api

import(
	"testing"
)

func TestGetNewTerritoryStatus(t *testing.T) {
	system := &System{
		Territories: []*SystemTerritory{
			&SystemTerritory{
				Status: TerritoryStatusContest,
			},
		},
		Planets: []Planet{
			Planet{ Player: &Player{ Faction: &Faction{ Id: 1 }}},
			Planet{},
		},
	}
	hostileTerritory := &Territory{
		Planet: &Planet{ Player: &Player{ Faction: &Faction{ Id: 2 }}},
	}
	allyTerritory := &Territory {
		Planet: &Planet{ Player: &Player{ Faction: &Faction{ Id: 1 }}},
	}

	if status := system.getNewTerritoryStatus(hostileTerritory); status != TerritoryStatusContest {
		t.Errorf("New territory status should equal 'contest', not %s", status)
	}
	if status := system.getNewTerritoryStatus(allyTerritory); status != TerritoryStatusPledge {
		t.Errorf("New territory status should equal 'pledge', not %s", status)
	}
}

func TestHasSystem(t *testing.T) {
	s := &System{
		Territories: []*SystemTerritory{
			&SystemTerritory{
				TerritoryId: 1,
				Territory: &Territory{ Id: 1 },
			},
			&SystemTerritory{
				TerritoryId: 2,
				Territory: &Territory{ Id: 2 },
			},
		},
	}
	st1 := &SystemTerritory{
		TerritoryId: 1,
		Territory: &Territory{ Id: 1 },
	}
	st2 := &SystemTerritory{
		TerritoryId: 3,
		Territory: &Territory{ Id: 3 },
	}

	if !st1.hasSystem(s) {
		t.Errorf("System should be considered to be in territory")
	}
	if st2.hasSystem(s) {
		t.Errorf("System should not be considered to be in territory")
	}
}