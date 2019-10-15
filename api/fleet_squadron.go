package api

import(
	"math"
	"net/http"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"strconv"
)

const FleetSquadronMaxQuantity = 50

type(
	FleetSquadron struct {
		tableName struct{} `json:"-" pg:"fleet__squadrons"`

		Id uint32 `json:"id"`
		FleetId uint16 `json:"-"`
		Fleet *Fleet `json:"fleet"`
		ShipModelId uint `json:"-"`
		ShipModel *ShipModel `json:"ship_model"`
		Quantity uint8 `json:"quantity"`
		CombatInitiative uint16 `json:"combat_initiative" pg:"-,use_zero"`
		CombatPosition *FleetGridPosition `json:"combat_position" pq:"type:jsonb,use_zero"`
		Position *FleetGridPosition `json:"position" pg:"type:jsonb"`
	}

	FleetGridPosition struct {
		X int8 `json:"x"`
		Y int8 `json:"y"`
	}
)

func CreateFleetSquadron(w http.ResponseWriter, r *http.Request) {
	fleetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    fleet := getFleet(uint16(fleetId))
	player := context.Get(r, "player").(*Player)

	if fleet.Player.Id != player.Id {
		panic(NewHttpException(http.StatusForbidden, "fleets.access_denied", nil))
	}
    if fleet.isOnJourney() {
        panic(NewHttpException(http.StatusBadRequest, "fleet.errors.ship_transfer_on_journey", nil))
    }
    if fleet.Location.Player.Id != player.Id {
        panic(NewHttpException(http.StatusBadRequest, "fleet.errors.ship_transfer_on_foreign_planet", nil))
	}
	SendJsonResponse(w, 200, fleet.createSquadron(DecodeJsonRequest(r)))
}

func GetFleetSquadrons(w http.ResponseWriter, r *http.Request) {
	fleetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    fleet := getFleet(uint16(fleetId))
	player := context.Get(r, "player").(*Player)

	if fleet.Player.Id != player.Id {
		panic(NewHttpException(http.StatusForbidden, "fleets.access_denied", nil))
	}
	SendJsonResponse(w, 200, fleet.getSquadrons())
}

func AssignFleetSquadronShips(w http.ResponseWriter, r *http.Request) {
	fleetId, _ := strconv.ParseUint(mux.Vars(r)["fleetId"], 10, 16)
	squadronId, _ := strconv.ParseUint(mux.Vars(r)["fleetId"], 10, 16)
	data := DecodeJsonRequest(r)
	fleet := getFleet(uint16(fleetId))
	squadron := fleet.getSquadron(uint32(squadronId))
	player := context.Get(r, "player").(*Player)
	shipModel := player.getShipModel(uint32(data["ship_model_id"].(float64)))

	if fleet.Player.Id != player.Id {
		panic(NewHttpException(http.StatusForbidden, "fleets.access_denied", nil))
	}
    if fleet.isOnJourney() {
        panic(NewHttpException(http.StatusBadRequest, "fleet.errors.ship_transfer_on_journey", nil))
    }
    if fleet.Location.Player.Id != player.Id {
        panic(NewHttpException(http.StatusBadRequest, "fleet.errors.ship_transfer_on_foreign_planet", nil))
	}
	squadron.assignShips(shipModel, uint8(data["quantity"].(float64)))
	w.WriteHeader(204)
	w.Write([]byte(""))
}

func (f *Fleet) createSquadron(data map[string]interface{}) *FleetSquadron {
	sm := f.Player.getShipModel(uint32(data["ship_model_id"].(float64)))

	squadron := &FleetSquadron{
		FleetId: f.Id,
		Fleet: f,
		ShipModelId: sm.Id,
		ShipModel: sm,
		Quantity: processSquadronQuantity(uint8(data["quantity"].(float64))),
		Position: f.processSquadronPosition(int8(data["x"].(float64)), int8(data["y"].(float64))),
	}
	if err := Database.Insert(squadron); err != nil {
		panic(NewException("Could not create fleet squadron", err))
	}
	f.Squadrons = append(f.Squadrons, squadron)
	return squadron
}

func (f *Fleet) hasSquadrons() bool {
    return len(f.Squadrons) > 0
}

func (f *Fleet) getSquadron(id uint32) *FleetSquadron {
	squadron := &FleetSquadron{}
	if err := Database.Model(squadron).Relation("ShipModel").Where("squadron.fleet_id = ?", f.Id).Where("squadron.id = ?", id).Select(); err != nil {
		panic(NewException("Could not retrieve fleet squadron", err))
	}
	return squadron
}

func (f *Fleet) getSquadrons() []*FleetSquadron {
    squadrons := make([]*FleetSquadron, 0)
    if err := Database.Model(&squadrons).Relation("ShipModel").Where("fleet_id = ?", f.Id).Select(); err != nil {
        panic(NewException("Could not retrieve fleet squadrons", err))
    }
    return squadrons
}

func (f *Fleet) deleteSquadron(s *FleetSquadron) {
    for i, squadron := range f.Squadrons {
        if s.Id == squadron.Id {
            f.Squadrons = append(f.Squadrons[:i], f.Squadrons[i+1:]...)
        }
    }
    s.delete()
}

func processSquadronQuantity(quantity uint8) uint8 {
	if !isValidSquadronQuantity(quantity) {
		panic(NewHttpException(400, "ships.squadrons.invalid_quantity", nil))
	}
	return quantity
}

func isValidSquadronQuantity(quantity uint8) bool {
	return quantity <= FleetSquadronMaxQuantity && quantity > 0
}

func (f *Fleet) processSquadronPosition(x, y int8) *FleetGridPosition {
	position := &FleetGridPosition{
		X: x,
		Y: y,
	}
	if !f.isValidSquadronPosition(position) {
		panic(NewHttpException(400, "ships.squadrons.invalid_position", nil))
	}
	return position
}

func (f *Fleet) isValidSquadronPosition(position *FleetGridPosition) bool {
	isXEven := position.X % 2 == 0
	isYEven := position.Y % 2 == 0
	if (isXEven && !isYEven) || (!isXEven && isYEven) {
		return false
	}
	for _, s := range f.Squadrons {
		if s.Position.X == position.X && s.Position.Y == position.Y {
			return false
		}
	}
	return true
}

func (fs *FleetSquadron) assignShips(sm *ShipModel, quantity uint8) {
	hangarGroup := fs.Fleet.Location.getHangarGroup(sm)
	if hangarGroup == nil {
		panic(NewHttpException(400, "fleet.errors.invalid_ship_type", nil))
	}
	
	requestedQuantity := quantity - fs.Quantity

	if uint16(requestedQuantity) > hangarGroup.Quantity {
		panic(NewHttpException(400, "fleet.errors.invalid_ship_number", nil))
	}
	if quantity > FleetSquadronMaxQuantity {
		panic(NewHttpException(400, "fleet.errors.invalid_ship_number", nil))
	}
	if requestedQuantity > 0 {
		hangarGroup.Quantity -= uint16(requestedQuantity)
		fs.Quantity += requestedQuantity
	} else {
		q := math.Abs(float64(requestedQuantity))
		hangarGroup.Quantity += uint16(q)
		fs.Quantity -= uint8(q)
	}
	fs.update()
	if hangarGroup.Quantity == 0 {
		hangarGroup.delete()
	} else {
		hangarGroup.update()
	}
}

func (fs *FleetSquadron) update() {
	if err := Database.Update(fs); err != nil {
		panic(NewException("Could not update fleet squadron", err))
	}
}

func (fs *FleetSquadron) delete() {
	if err := Database.Delete(fs); err != nil {
		panic(NewException("Could not remove fleet squadron", err))
	}
}