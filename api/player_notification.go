package api

import(
	"time"
	"net/http"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"strconv"
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
		tableName struct{} `json:"-" pg:"player__notifications"`

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

func DeleteNotification(w http.ResponseWriter, r *http.Request) {
	player := context.Get(r, "player").(*Player)
	id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	
	player.deleteNotification(uint32(id))

	w.WriteHeader(204)
	w.Write([]byte(""))
}

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

func (p *Player) deleteNotification(id uint32) {
	n := &Notification{}

	if err := Database.Model(n).Where("id = ?", id).Select(); err != nil {
		panic(NewException("Could not retrieve ", err))
	}
	if n.PlayerId != p.Id {
		panic(NewHttpException(403, "Forbidden", nil))
	}
	n.delete()
}

func (n *Notification) delete() {
	if err := Database.Delete(n); err != nil {
		panic(NewException("Could not delete notification", err))
	}
}