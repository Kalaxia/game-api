package main

import (
	"net/http"
	"log"
	_ "kalaxia-game-api/security"
	"kalaxia-game-api/route"
)

func main() {
  	log.Fatal(http.ListenAndServe(":80", route.NewRouter()))
}
