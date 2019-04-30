
package api

import(
	"encoding/json"
)

var WsHub *Hub
// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	players map[string]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewWsHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		players:    make(map[string]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.players, client.player.Pseudo)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
					delete(h.players, client.player.Pseudo)
				}
			}
		}
	}
}

func (h *Hub) sendTo(player *Player, message *WsMessage) {
	if client, ok := h.players[player.Pseudo]; ok {
		jsonData, err := json.Marshal(message)
		if err != nil {
			panic(NewException("websocket.encoding_error", err))
		}
		client.send <- jsonData
	}
}