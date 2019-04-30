package main

import (
	"net/http"
	"log"
	"kalaxia-game-api/api"
)

func main() {
	initConfigurations()
	initScheduledTasks()

	hub := api.NewWsHub()
	go hub.Run()

	router := NewRouter()
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		api.ServeWs(hub, w, r)
	})
  	log.Fatal(http.ListenAndServe(":80", router))
}

func initConfigurations() {
	api.InitDatabase()
	api.InitRsaVault()
	api.InitShipConfiguration()
}

func initScheduledTasks() {
	api.InitScheduler()

    api.Scheduler.AddHourlyTask(func () { api.CalculatePlayersWage() })
    api.Scheduler.AddHourlyTask(func () { api.CalculatePlanetsProductions() })
	api.Scheduler.AddHourlyTask(func () { api.CheckShipsBuildingState() })
	
	api.InitFleetJourneys()
}