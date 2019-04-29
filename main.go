package main

import (
	"net/http"
	"log"
	_ "kalaxia-game-api/security"
	"kalaxia-game-api/route"
	"kalaxia-game-api/websocket"
)

func main() {
	hub := websocket.NewHub()
	go hub.Run()

	router := route.NewRouter()
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.Serve(hub, w, r)
	})
  	log.Fatal(http.ListenAndServe(":80", router))
}
