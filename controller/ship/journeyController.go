package shipController


import (
    "net/http"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "kalaxia-game-api/exception"
	"kalaxia-game-api/manager/ship"
    "kalaxia-game-api/model"
	"kalaxia-game-api/utils"
    "strconv"
    "math"
)

func GetJourney (w http.ResponseWriter, r *http.Request){
	
	player := context.Get(r, "player").(*model.Player)
	idFleet, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16);
	fleet := shipManager.GetFleetOnJourney (uint16(idFleet));
	
	
	if player.Id != fleet.Player.Id { // the player does not own the planet
		panic(exception.NewHttpException(http.StatusForbidden, "", nil));
	}
	if !fleet.IsOnJourney() {
		panic(exception.NewHttpException(400, "This journey has ended", nil));
	}
	
	utils.SendJsonResponse(w, 200,fleet.Journey);
}

func GetFleetSteps (w http.ResponseWriter, r *http.Request){
    player := context.Get(r, "player").(*model.Player)
    idFleet, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16);
	fleet := shipManager.GetFleetOnJourney (uint16(idFleet));
    
    if player.Id != fleet.Player.Id { // the player does not own the planet
		panic(exception.NewHttpException(http.StatusForbidden, "", nil));
	}
	if !fleet.IsOnJourney() {
		panic(exception.NewHttpException(400, "This journey has ended", nil));
	}
    
    utils.SendJsonResponse(w, 200,shipManager.GetStepsByJourneyId(fleet.Journey.Id));
}

func CreateJourney (w http.ResponseWriter, r *http.Request) {
	
	
	utils.SendJsonResponse(w, 202,"");
}

func SendFleetOnJourney (w http.ResponseWriter, r *http.Request){
    player := context.Get(r, "player").(*model.Player)
	idFleet, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16);
    fleet := shipManager.GetFleet(uint16(idFleet));
    
    if player.Id != fleet.Player.Id { // the player does not own the planet
		panic(exception.NewHttpException(http.StatusForbidden, "", nil));
	}
    
    if fleet.IsOnJourney() {
		panic(exception.NewHttpException(400, "Fleet already on journey", nil));
	}
    
    data := utils.DecodeJsonRequest(r)["steps"].([]map[string]interface{});
    var planetIds []uint16;
    var xPos []float64;
    var yPos []float64;
    
    for i,_ := range data {
        if data[i]["planetId"].(float64) == 0. && (data[i]["x"].(float64) == math.NaN() || data[i]["y"].(float64) == math.NaN() ){
            panic(exception.NewHttpException(400, "step not well defined", nil));
        }
        planetIds = append(planetIds,uint16(data[i]["planetId"].(float64)));
        xPos = append(xPos,data[i]["x"].(float64));
        yPos = append(xPos,data[i]["y"].(float64));
    }
    
    steps := shipManager.SendFleetOnJourney(planetIds,xPos,yPos,fleet);
    
    utils.SendJsonResponse(w, 202, steps);
    
}

func AddStepsToJourney (w http.ResponseWriter, r *http.Request){
    player := context.Get(r, "player").(*model.Player)
	idFleet, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16);
    fleet := shipManager.GetFleet(uint16(idFleet));
    
    if player.Id != fleet.Player.Id { // the player does not own the planet
		panic(exception.NewHttpException(http.StatusForbidden, "", nil));
	}
    
    if ! fleet.IsOnJourney() {
		panic(exception.NewHttpException(400, "Fleet is not on journey", nil));
	}
    
    data := utils.DecodeJsonRequest(r)["steps"].([]map[string]interface{});
    var planetIds []uint16;
    var xPos []float64;
    var yPos []float64;
    
    
    for i,_ := range data {
        if data[i]["planetId"].(float64) == 0. && (data[i]["x"].(float64) == math.NaN() || data[i]["y"].(float64) == math.NaN() ){
            panic(exception.NewHttpException(400, "step not well defined", nil));
        }
        planetIds = append(planetIds,uint16(data[i]["planetId"].(float64)));
        xPos = append(xPos,data[i]["x"].(float64));
        yPos = append(xPos,data[i]["y"].(float64));
    }
    
    steps := shipManager.AddStepsToJourney(fleet.Journey,planetIds,xPos,yPos);
    
    utils.SendJsonResponse(w, 202, steps);
    
}

func GetRange(w http.ResponseWriter, r *http.Request){
    
    // ID for later if diffrent fleet have diffrent range
    utils.SendJsonResponse(w, 202, shipManager.GetRange());
}

func GetTimeLaws(w http.ResponseWriter, r *http.Request){
    
    // ID for later if diffrent fleet have diffrent range
    utils.SendJsonResponse(w, 202, shipManager.GetTimeLaws());
}

// TODO remove steps
