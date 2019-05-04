package api

import(
	"time"
)

const(
	NotificationTypeBuilding = "building"
	NotificationTypeResearch = "research"
	NotificationTypeShipyard = "shipyard"
	NotificationTypeTrade = "trade"
	NotificationTypeMilitary = "military"
	NotificationTypeDiplomacy = "diplomacy"
	NotificationTypeFaction = "faction"
)

type(
	Notification struct {
		TableName struct{} `json:"-" sql:"player__notifications"`

		Id uint32 `json:"id"`
		Player *Player `json:"player"`
		PlayerId uint16 `json:"-"`
		Type string `json:"type"`
		Content string `json:"content"`
		Data map[string]interface{} `json:"data"`
		CreatedAt time.Time `json:"created_at"`
		ReadAt time.Time `json:"read_at"`
	}
	Notifications []Notification
)

func (p *Player) notify(nType string, content string, data map[string]interface{}) *Notification {
	notification := &Notification {
		Player: p,
		PlayerId: p.Id,
		Type: nType,
		Content: content,
		Data: data,
		CreatedAt: time.Now(),
	}
	if err := Database.Insert(notification); err != nil {
		panic(NewException("notifications.creation_fail", err))
	}
	WsHub.sendTo(p, &WsMessage{ Action: "addNotification", Data: notification })
	return notification
}

func (p *Player) getNotifications() {
	notifications := make(Notifications, 0)

	if err := Database.Model(&notifications).Where("player_id = ?", p.Id).Order("created_at DESC").Select(); err != nil {
		panic(NewHttpException(404, "players.not_found", err))
	}
}