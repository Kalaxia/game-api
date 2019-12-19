package api

import(
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
		Quantity uint8 `json:"quantity" pg:",use_zero"`
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
    if !fleet.isOnPlanet() {
        panic(NewHttpException(http.StatusBadRequest, "fleet.errors.ship_transfer_on_journey", nil))
    }
    if fleet.Place.Planet.Player.Id != player.Id {
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
	squadronId, _ := strconv.ParseUint(mux.Vars(r)["squadronId"], 10, 16)
	data := DecodeJsonRequest(r)
	fleet := getFleet(uint16(fleetId))
	squadron := fleet.getSquadron(uint32(squadronId))
	player := context.Get(r, "player").(*Player)

	if fleet.Player.Id != player.Id {
		panic(NewHttpException(http.StatusForbidden, "fleets.access_denied", nil))
	}
    if !fleet.isOnPlanet() {
        panic(NewHttpException(http.StatusBadRequest, "fleet.errors.ship_transfer_on_journey", nil))
    }
    if fleet.Place.Planet.Player.Id != player.Id {
        panic(NewHttpException(http.StatusBadRequest, "fleet.errors.ship_transfer_on_foreign_planet", nil))
	}
	squadron.assignShips(uint8(data["quantity"].(float64)))
	w.WriteHeader(204)
	w.Write([]byte(""))
}

func (f *Fleet) createSquadron(data map[string]interface{}) *FleetSquadron {
	f.Squadrons = f.getSquadrons()
	sm := f.Player.getShipModel(uint32(data["ship_model_id"].(float64)))
	position := data["position"].(map[string]interface{})
	quantity := processSquadronQuantity(uint8(data["quantity"].(float64)))

	squadron := &FleetSquadron{
		FleetId: f.Id,
		Fleet: f,
		ShipModelId: sm.Id,
		ShipModel: sm,
		Quantity: 0,
		Position: f.processSquadronPosition(int8(position["x"].(float64)), int8(position["y"].(float64))),
	}
	if err := Database.Insert(squadron); err != nil {
		panic(NewException("Could not create fleet squadron", err))
	}
	squadron.assignShips(quantity)
	// avoid infinite loop at JSON serialization
	squadron.Fleet = nil
	f.Squadrons = append(f.Squadrons, squadron)
	return squadron
}

func (f *Fleet) hasSquadrons() bool {
    return len(f.Squadrons) > 0
}

func (f *Fleet) getSquadron(id uint32) *FleetSquadron {
	squadron := &FleetSquadron{
		Fleet: f,
		FleetId: f.Id,
	}
	if err := Database.Model(squadron).Relation("ShipModel").Where("fleet_squadron.fleet_id = ?", f.Id).Where("fleet_squadron.id = ?", id).Select(); err != nil {
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
    s.delete()
    for i, squadron := range f.Squadrons {
        if s.Id == squadron.Id {
			f.Squadrons = append(f.Squadrons[:i], f.Squadrons[i+1:]...)
			break
        }
    }
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
	for _, s := range f.Squadrons {
		if s.Position.X == position.X && s.Position.Y == position.Y {
			return false
		}
	}
	return true
}

func (fs *FleetSquadron) assignShips(quantity uint8) {
	hangarGroup := fs.Fleet.Place.Planet.findOrCreateHangarGroup(fs.ShipModel)

	requestedQuantity := int8(quantity) - int8(fs.Quantity)

	if requestedQuantity > 0 && uint16(requestedQuantity) > hangarGroup.Quantity {
		panic(NewHttpException(400, "fleet.errors.not_enough_ships", nil))
	}
	if quantity > FleetSquadronMaxQuantity {
		panic(NewHttpException(400, "fleet.errors.invalid_quantity", nil))
	}
	hangarGroup.addShips(-requestedQuantity)
	if quantity > 0 {
		fs.Quantity = quantity
		fs.update()
		return
	}
	fs.delete()
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