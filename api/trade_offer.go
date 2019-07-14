package api

import(
    "net/http"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "math"
    "strconv"
    "time"
)

type(
    OfferInterface interface {
        getTotalPrice() float32
        cancel()
        update()
        delete()
    }
    Offer struct {
        Id uint32 `json:"id"`
        Type string `json:"type" sql:"-"`
        Operation string `json:"operation"`
        LocationId uint16 `json:"-"`
        Location *Planet `json:"location"`
        DestinationId uint16 `json:"-"`
        Destination *Planet `json:"destination"`
        CreatedAt time.Time `json:"created_at"`
        AcceptedAt time.Time `json:"accepted_at"`
    }
    ResourceOffer struct {
        TableName struct{} `json:"-" sql:"trade__resource_offers"`

        Offer

        Resource string `json:"resource"`
        Quantity uint16 `json:"quantity"`
        LotQuantity uint16 `json:"lot_quantity"`
        Price float32 `json:"price"`
    }
    ModelOffer struct {
        TableName struct{} `json:"-" sql:"trade__model_offers"`

        Offer

        ModelId uint `json:"-"`
        Model *ShipModel `json:"model"`
        Price uint16 `json:"price"`
    }
    ShipOffer struct {
        TableName struct{} `json:"-" sql:"trade__ship_offers"`

        ModelOffer

        Quantity uint16 `json:"quantity"`
        LotQuantity uint16 `json:"lot_quantity"`
    }
)

func CreateOffer(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)
    planetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := getPlayerPlanet(uint16(planetId), player.Id)

    if planet.Player.Id != player.Id {
        panic(NewHttpException(403, "You do not control this planet", nil))
    }
    SendJsonResponse(w, 201, planet.createOffer(DecodeJsonRequest(r)))
}

func SearchOffers(w http.ResponseWriter, r *http.Request) {
    SendJsonResponse(w, 200, searchOffers(DecodeJsonRequest(r)))
}

func CancelOffer(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)
    offerId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    offer := getOffer(uint32(offerId))

	if (offer.Location.Player.Id != player.Id) {
		panic(NewHttpException(403, "planets.forbidden", nil))
	}
    offer.cancel()

    w.WriteHeader(204)
    w.Write([]byte(""))
}

func GetOffer(w http.ResponseWriter, r *http.Request) {
    offerId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    SendJsonResponse(w, 200, getOffer(uint32(offerId)))
}

func AcceptOffer(w http.ResponseWriter, r *http.Request) {
    data := DecodeJsonRequest(r)
    offerId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    nbLots := uint16(data["nb_lots"].(float64))
    planetId := data["planet_id"].(float64)
    player := context.Get(r, "player").(*Player)
    planet := getPlayerPlanet(uint16(planetId), player.Id)

    planet.acceptOffer(uint32(offerId), nbLots)

    w.WriteHeader(204)
    w.Write([]byte(""))
}


func getOffer(id uint32) *ResourceOffer {
    offer := &ResourceOffer{}
    if err := Database.Model(offer).Column("Location.Player.Faction", "Location.System", "Location.Storage").Where("resource_offer.id = ?", id).Select(); err != nil {
        panic(NewHttpException(404, "Offer not found", err))
    }
    offer.Type = "resources"
    return offer
}
//
// func GetPlanetOffers(location *model.Planet) []model.OfferInterface {
//     offers := make([]model.Offer, 0)
//     if err := Database.Model(&offers).Where("location_id = ?", location.Id).Select(); err != nil {
//         panic(NewHttpException(500, "Planet offers could not be retrieved", err))
//     }
//     return offers
// }
//
// func GetDestinationOffers(location *model.Planet) []model.OfferInterface {
//     offers := make([]model.Offer, 0)
//     if err := Database.Model(&offers).Where("destination_id = ?", location.Id).Select(); err != nil {
//         panic(NewHttpException(500, "Planet offers could not be retrieved", err))
//     }
//     return offers
// }

func (o *ResourceOffer) cancel() {
    o.Location.Storage.storeResource(o.Resource, int16(o.Quantity))
    o.Location.Storage.update()
    o.delete()

    WsHub.sendBroadcast(&WsMessage{ Action: "cancelTradeOffer", Data: o })
}

func (o *ShipOffer) cancel() {
    o.delete()

    WsHub.sendBroadcast(&WsMessage{ Action: "cancelTradeOffer", Data: o })
}

func (o ModelOffer) cancel() {
    o.delete()
    
    WsHub.sendBroadcast(&WsMessage{ Action: "cancelTradeOffer", Data: o })
}

func searchOffers(data map[string]interface{}) []OfferInterface {
    offers := make([]OfferInterface, 0)

    operation := data["operation"].(string)

    resourceOffers := make([]*ResourceOffer, 0)
    if err := Database.Model(&resourceOffers).Column("Location.Player.Faction", "Location.System").Where("operation = ?", operation).Select(); err != nil {
        panic(NewHttpException(500, "Resource offers could not be retrieved", err))
    }
    // shipOffers := make([]*model.ShipOffer, 0)
    // if err := Database.Model(&shipOffers).Where("operation = ?", operation).Select(); err != nil {
    //     panic(NewHttpException(500, "Ship offers could not be retrieved", err))
    // }
    // modelOffers := make([]*model.ModelOffer, 0)
    // if err := Database.Model(&modelOffers).Where("operation = ?", operation).Select(); err != nil {
    //     panic(NewHttpException(500, "Model offers could not be retrieved", err))
    // }

    for _, offer := range resourceOffers {
        offers = append(offers, offer)
    }
    // for _, offer := range shipOffers {
    //     offers = append(offers, offer)
    // }
    // for _, offer := range modelOffers {
    //     offers = append(offers, offer)
    // }

    return offers
}

func (p *Planet) createOffer(data map[string]interface{}) OfferInterface {
    var offer OfferInterface
    switch goodType := data["good_type"].(string); goodType {
        case "resources":
            offer = p.createResourceOffer(data);
        case "ships":
            offer = p.createShipOffer(data);
        case "models":
            offer = p.createModelOffer(data);
    }
    WsHub.sendBroadcast(&WsMessage{ Action: "addTradeOffer", Data: offer })
    return offer
}

func (p *Planet) createResourceOffer(data map[string]interface{}) *ResourceOffer {
    offer := &ResourceOffer{
        Resource: data["resource"].(string),
        Quantity: uint16(data["quantity"].(float64)),
        LotQuantity: uint16(data["lot_quantity"].(float64)),
        Price: float32(data["price"].(float64)),
    }
    offer.Type = "resources"
    offer.Operation = data["operation"].(string)
    offer.LocationId = p.Id
    offer.Location = p
    offer.CreatedAt = time.Now()
    if offer.Price < 1 {
        panic(NewHttpException(400, "trade.offers.invalid_price", nil))
    }
    if offer.Quantity < offer.LotQuantity {
        panic(NewHttpException(400, "Lot quantity cannot be lesser than total quantity", nil))
    }
    if !p.Storage.storeResource(offer.Resource, -int16(offer.Quantity)) {
        panic(NewHttpException(400, "Not enough resources in storage", nil))
    }
    p.Storage.update()
    if err := Database.Insert(offer); err != nil {
        panic(NewHttpException(500, "Resource offer could not be created", err))
    }
    return offer
}

func (p *Planet) createShipOffer(data map[string]interface{}) *ShipOffer {
    shipModel := p.Player.getShipModel(data["model"].(uint32))

    offer := &ShipOffer{
        Quantity: data["quantity"].(uint16),
        LotQuantity: data["lot_quantity"].(uint16),
    }
    offer.Operation = data["operation"].(string)
    offer.LocationId = p.Id
    offer.Location = p
    offer.CreatedAt = time.Now()
    offer.Price = data["price"].(uint16)
    offer.Model = shipModel
    offer.ModelId = shipModel.Id
    if err := Database.Insert(offer); err != nil {
        panic(NewHttpException(500, "Ship offer could not be created", err))
    }
    return offer
}

func (p *Planet) createModelOffer(data map[string]interface{}) *ModelOffer {
    shipModel := p.Player.getShipModel(data["model"].(uint32))

    offer := &ModelOffer{
        Price: data["price"].(uint16),
        ModelId: shipModel.Id,
        Model: shipModel,
    }
    offer.Operation = data["operation"].(string)
    offer.LocationId = p.Id
    offer.Location = p
    offer.CreatedAt = time.Now()
    if err := Database.Insert(offer); err != nil {
        panic(NewHttpException(500, "Ship model offer could not be created", err))
    }
    return offer
}

func (p *Planet) acceptOffer(offerId uint32, nbLots uint16) {
    offer := getOffer(offerId)
    if offer == nil {
        panic(NewHttpException(404, "Offer not found", nil))
    }
    if offer.Location.Player.Id == p.Player.Id {
        panic(NewHttpException(400, "You can't accept your own offers", nil))
    }

    quantity := nbLots * offer.LotQuantity
    price := int32(math.Ceil(float64(offer.Price) * float64(quantity)))

    if quantity % offer.LotQuantity > 0 {
        panic(NewHttpException(400, "There can be no extra resource out of lots", nil))
    }
    if quantity > offer.Quantity {
        panic(NewHttpException(400, "You can't demand more lots than available", nil))
    }
    if !p.Player.updateWallet(-price) {
        panic(NewHttpException(400, "Not enough money", nil))
    }
    offer.Location.Player.updateWallet(price)
    p.Player.update()
    offer.Location.Player.update()

    p.Storage.storeResource(offer.Resource, int16(quantity))
    p.Storage.update()

    offer.Quantity -= quantity
    WsHub.sendTo(offer.Location.Player, &WsMessage{ Action: "updateWallet", Data: map[string]uint32{
        "wallet": offer.Location.Player.Wallet,
    }})
    offer.Location.Player.notify(
        NotificationTypeTrade,
        "notifications.trade.accepted_offer",
        map[string]interface{}{
            "offer": offer,
            "player": p.Player,
            "quantity": quantity,
            "price": price,
        },
    )
    WsHub.sendBroadcast(&WsMessage{ Action: "updateTradeOffer", Data: offer })
    if offer.Quantity == 0 {
        offer.delete()
        return
	} 
    offer.update()
}

func (o *ResourceOffer) update() {
    if err := Database.Update(o); err != nil {
        panic(NewHttpException(500, "Offer could not be accepted", err))
    }
}

func (o *ShipOffer) update() {
    if err := Database.Update(o); err != nil {
        panic(NewHttpException(500, "Offer could not be accepted", err))
    }
}

func (o *ModelOffer) update() {
    if err := Database.Update(o); err != nil {
        panic(NewHttpException(500, "Offer could not be accepted", err))
    }
}

func (o *ResourceOffer) delete() {
	if err := Database.Delete(o); err != nil {
		panic(NewHttpException(500, "Offer could not be deleted", err))
	}
}

func (o *ShipOffer) delete() {
	if err := Database.Delete(o); err != nil {
		panic(NewHttpException(500, "Offer could not be deleted", err))
	}
}

func (o *ModelOffer) delete() {
	if err := Database.Delete(o); err != nil {
		panic(NewHttpException(500, "Offer could not be deleted", err))
	}
}

func (o *ResourceOffer) getTotalPrice() float32 {
    return o.Price
}

func (o *ShipOffer) getTotalPrice() float32 {
    return float32(o.Price)
}

func (o *ModelOffer) getTotalPrice() float32 {
    return float32(o.Price)
}