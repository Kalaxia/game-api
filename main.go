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
	api.InitRedisClient()
	api.InitRsaVault()
	api.InitShipConfiguration()
	api.InitPlanetConstructions()
	api.InitFactionMotions()
}

func initScheduledTasks() {
	api.InitScheduler()

    api.Scheduler.AddHourlyTask(func () { api.CalculatePlayersWage() })
	api.Scheduler.AddHourlyTask(func () { api.CalculatePlanetsProductions() })

	api.Scheduler.AddDailyTask(func () { api.CalculateFactionsWages() })
	
	api.InitFleetJourneys()
}

func initWebsocketHub() {
	api.WsHub = api.NewWsHub()
	go api.WsHub.Run()
}