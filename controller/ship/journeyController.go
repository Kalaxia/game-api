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
)

func GetJourney (w http.ResponseWriter, r *http.Request){
	
	player := context.Get(r, "player").(*model.Player)
	idFleet, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	fleet := shipManager.GetFleetOnJourney (uint16(idFleet))
	
	if player.Id != fleet.Player.Id { // the player does not own the planet
		panic(exception.NewHttpException(http.StatusForbidden, "", nil))
	}
	if !fleet.IsOnJourney() {
		panic(exception.NewHttpException(400, "This journey has ended", nil))
	}
	
	utils.SendJsonResponse(w, 200,fleet.Journey)
}

func GetFleetSteps (w http.ResponseWriter, r *http.Request){
    player := context.Get(r, "player").(*model.Player)
    idFleet, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
	fleet := shipManager.GetFleetOnJourney (uint16(idFleet))
    
    if player.Id != fleet.Player.Id { // the player does not own the planet
		panic(exception.NewHttpException(http.StatusForbidden, "", nil))
	}
	if !fleet.IsOnJourney() {
		panic(exception.NewHttpException(400, "This journey has ended", nil))
	}
    
    utils.SendJsonResponse(w, 200,shipManager.GetStepsByJourneyId(fleet.Journey.Id))
}


func SendFleetOnJourney (w http.ResponseWriter, r *http.Request){
    player := context.Get(r, "player").(*model.Player)
	idFleet, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    fleet := shipManager.GetFleet(uint16(idFleet))
    
    if player.Id != fleet.Player.Id { // the player does not own the planet
		panic(exception.NewHttpException(http.StatusForbidden, "", nil))
	}
    
    if fleet.IsOnJourney() {
		panic(exception.NewHttpException(400, "Fleet already on journey", nil))
	}
    
    data := utils.DecodeJsonRequest(r)["steps"].([]interface{})
    
    utils.SendJsonResponse(w, 201, shipManager.SendFleetOnJourney(fleet, data))
}

func AddStepsToJourney (w http.ResponseWriter, r *http.Request){
    player := context.Get(r, "player").(*model.Player)
	idFleet, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    fleet := shipManager.GetFleet(uint16(idFleet))
    
    if player.Id != fleet.Player.Id { // the player does not own the planet
		panic(exception.NewHttpException(http.StatusForbidden, "", nil))
	}
    
    if ! fleet.IsOnJourney() {
		panic(exception.NewHttpException(400, "Fleet is not on journey", nil))
	}
    
    data := utils.DecodeJsonRequest(r)["steps"].([]interface{})
    
    steps := shipManager.AddStepsToJourney(fleet, data)
    
    utils.SendJsonResponse(w, 202, steps)
}

func GetRange(w http.ResponseWriter, r *http.Request){
    // ID for later if diffrent fleet have diffrent range
    utils.SendJsonResponse(w, 200, shipManager.GetRange())
}

func GetTimeLaws(w http.ResponseWriter, r *http.Request){
    // ID for later if diffrent fleet have diffrent range
    utils.SendJsonResponse(w, 200, shipManager.GetTimeLaws())
}

func RemoveStepAndFollowingFormJourneyAssociatedWithFleet (w http.ResponseWriter, r *http.Request){
    // Cancel a journey form a setp, it remove this step and evryone after this one
    player := context.Get(r, "player").(*model.Player)
    idFleet, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    fleet := shipManager.GetFleet(uint16(idFleet))
    
    if player.Id != fleet.Player.Id{
        panic(exception.NewHttpException(http.StatusForbidden, "", nil))
    }
    if !fleet.IsOnJourney(){
        panic(exception.NewHttpException(400, "Fleet is not on journey", nil))
    }
    
    stepId,_ := strconv.ParseUint(mux.Vars(r)["idStep"], 10, 16)
    step := shipManager.GetStep(uint16(stepId))
    
    if fleet.JourneyId != step.JourneyId{
        panic(exception.NewHttpException(400, "This step is not linked to the fleet journey", nil))
    }
    if fleet.Journey.CurrentStepId == step.Id{
        panic(exception.NewHttpException(400, "current step cannot be canceled", nil))
    }
    if fleet.Journey.CurrentStep.StepNumber >= step.StepNumber{
        panic(exception.NewHttpException(400, "cannot remove step with smaler step number that the current one ", nil))
    }
    
    shipManager.RemoveStepsAndFollowingFromJourney(fleet.Journey,step)
    utils.SendJsonResponse(w,204,"Deleted");
    
}
