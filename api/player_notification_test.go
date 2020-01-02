package api

import(
	"testing"
)

func TestNotify(t *testing.T) {
	InitDatabaseMock()
	InitWsHubMock()

	p := &Player{ Id: 1 }
	n := p.notify(NotificationTypeBuilding, "planet.buildings.construction_success", map[string]interface{}{
		"planet": &Planet{ Id: 1, },
		"building": &Building{ Id: 7 },
	})
	if n.Type != NotificationTypeBuilding {
		t.Errorf("Notification type should be 'building', got %s", n.Type)
	}
	if n.Player.Id != 1 {
		t.Errorf("Notification player id should equal 1, got %d", n.Player.Id)
	}
	if n.Content != "planet.buildings.construction_success" {
		t.Errorf("Notification content should equal 'planet.buildings.construction_success', got %s", n.Content)
	}
	if n.Data["planet"].(*Planet).Id != 1 {
		t.Errorf("Notification data should contain planet id 1, got %d", n.Data["planet"].(*Planet).Id)
	}
}