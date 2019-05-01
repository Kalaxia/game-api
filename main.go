package main

import (
	"net/http"
	"log"
	"kalaxia-game-api/api"
)

func main() {
	initConfigurations()
	initScheduledTasks()
	initWebsocketHub()

	router := NewRouter()
	router.HandleFunc("/ws", api.ServeWs)
  	log.Fatal(http.ListenAndServe(":80", router))
}

func initConfigurations() {
	api.InitDatabase()
	api.InitRsaVault()
	api.InitShipConfiguration()
	api.InitPlanetConstructions()
}

func initScheduledTasks() {
	api.InitScheduler()

    api.Scheduler.AddHourlyTask(func () { api.CalculatePlayersWage() })
    api.Scheduler.AddHourlyTask(func () { api.CalculatePlanetsProductions() })
	api.Scheduler.AddHourlyTask(func () { api.CheckShipsBuildingState() })
	
	api.InitFleetJourneys()
}

func initWebsocketHub() {
	api.WsHub = api.NewWsHub()
	go api.WsHub.Run()
}