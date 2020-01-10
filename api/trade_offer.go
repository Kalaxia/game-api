package api

import(
    "net/http"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "math"
    "strconv"
    "time"
)

const(
    TradeOfferOperationSell = "sell"
    TradeOfferOperationBuy = "buy"
)

type(
    OfferInterface interface {
        cancel()
        update()
        delete()
    }
    Offer struct {
        tableName struct{} `pg:"trade__offers"`

        Id uint32 `json:"id"`
        Type string `json:"type" pg:"-"`
        Operation string `json:"operation"`
        LocationId uint16 `json:"-"`
        Location *Planet `json:"location"`
        DestinationId uint16 `json:"-"`
        Destination *Planet `json:"destination"`
        Price uint16 `json:"price"`
        CreatedAt time.Time `json:"created_at"`
        AcceptedAt time.Time `json:"accepted_at"`
    }
    ResourceOffer struct {
        Offer `pg:",inherit"`

        Resource string `json:"resource"`
        Quantity uint16 `json:"quantity"`
        LotQuantity uint16 `json:"lot_quantity"`
    }
    ModelOffer struct {
        Offer `pg:",inherit"`

        ModelId uint `json:"-"`
        Model *ShipModel `json:"model"`
    }
    ShipOffer struct {
        Offer `pg:",inherit"`

        Quantity uint16 `json:"quantity"`
        LotQuantity uint16 `json:"lot_quantity"`
        ModelId uint `json:"-"`
        Model *ShipModel `json:"model"`
    }
)

func CreateOffer(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)
    planetId, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := player.getPlanet(uint16(planetId))

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
    planetId := data["planet_id"].(float64)
    player := context.Get(r, "player").(*Player)
    planet := player.getPlanet(uint16(planetId))

    offer := getOffer(uint32(offerId))
    if offer == nil {
        panic(NewHttpException(404, "Offer not found", nil))
    }
    offer.accept(planet, data)

    w.WriteHeader(204)
    w.Write([]byte(""))
}


func getOffer(id uint32) *ResourceOffer {
    offer := &ResourceOffer{}
    if err := Database.Model(offer).Relation("Location.Player.Faction").
        Relation("Location.System").
        Relation("Location.Storage").Where("offer.id = ?", id).Select(); err != nil {
        panic(NewHttpException(404, "Offer not found", err))
    }
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
    o.Offer.cancel()
}

func (o *Offer) cancel() {
    o.delete()

    WsHub.sendBroadcast(&WsMessage{ Action: "cancelTradeOffer", Data: o })
}

func searchOffers(data map[string]interface{}) []*ResourceOffer {
    offers := make([]*ResourceOffer, 0)

    operation := data["operation"].(string)

    if err := Database.Model(&offers).Relation("Location.Player.Faction").Relation("Location.System").Where("operation = ?", operation).Select(); err != nil {
        panic(NewHttpException(500, "Resource offers could not be retrieved", err))
    }
    return offers
}

func (p *Planet) createOffer(data map[string]interface{}) OfferInterface {
    if operation := data["operation"].(string); operation != TradeOfferOperationBuy && operation != TradeOfferOperationSell {
        panic(NewHttpException(400, "trade.offers.invalid_operation", nil))
    }
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

func (p *Planet) newOffer(data map[string]interface{}) Offer {
    return Offer{
        Type: data["good_type"].(string),
        Operation: data["operation"].(string),
        Price: uint16(data["price"].(float64)),
        LocationId: p.Id,
        Location: p,
        CreatedAt: time.Now(),
    }
}

func (p *Planet) createResourceOffer(data map[string]interface{}) *ResourceOffer {
    quantity := data["quantity"].(float64)
    offer := &ResourceOffer{
        Offer: p.newOffer(data),
        Resource: data["resource"].(string),
        Quantity: uint16(quantity),
        LotQuantity: uint16(data["lot_quantity"].(float64)),
    }
    if quantity < 100 {
        panic(NewHttpException(400, "trade.offers.invalid_quantity", nil))
    }
    if offer.LotQuantity < 1 {
        panic(NewHttpException(400, "trade.offers.invalid_lot_quantity", nil))
    }
    if offer.Price < 1 || offer.Price > 1000 {
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
    shipModel := getShipModel(data["model"].(uint32))

    offer := &ShipOffer{
        Offer: p.newOffer(data),
        Quantity: data["quantity"].(uint16),
        LotQuantity: data["lot_quantity"].(uint16),
        Model: shipModel,
        ModelId: shipModel.Id,
    }
    if err := Database.Insert(offer); err != nil {
        panic(NewHttpException(500, "Ship offer could not be created", err))
    }
    return offer
}

func (p *Planet) createModelOffer(data map[string]interface{}) *ModelOffer {
    shipModel := getShipModel(data["model"].(uint32))

    offer := &ModelOffer{
        Offer: p.newOffer(data),
        ModelId: shipModel.Id,
        Model: shipModel,
    }
    if err := Database.Insert(offer); err != nil {
        panic(NewHttpException(500, "Ship model offer could not be created", err))
    }
    return offer
}

func (o *Offer) checkProposal(p *Planet, data map[string]interface{}) {
    if o.Location.Player.Id == p.Player.Id {
        panic(NewHttpException(400, "You can't accept your own offers", nil))
    }
}

func (o *Offer) performTransaction(p *Player, price int32) {
    finalPrice, gain, purchaseTaxes, salesTaxes := o.processTaxes(p, price)
    if !p.updateWallet(-finalPrice) {
        panic(NewHttpException(400, "Not enough money", nil))
    }
    o.Location.Player.updateWallet(gain)
    p.update()
    o.Location.Player.update()
    o.performFactionTransaction(p, purchaseTaxes, salesTaxes)
}

func (o *Offer) performFactionTransaction(p *Player, purchaseTaxes, salesTaxes int32) {
    if o.Operation == TradeOfferOperationBuy {
        o.Location.Player.Faction.Wallet += purchaseTaxes
        p.Faction.Wallet += salesTaxes
    } else {
        p.Faction.Wallet += purchaseTaxes
        o.Location.Player.Faction.Wallet += salesTaxes
    }
    p.Faction.update()
    o.Location.Player.Faction.update()
}

func (o *Offer) processTaxes(p *Player, price int32) (finalPrice, gain, purchaseTaxes, salesTaxes int32) {
    if p.Faction.Id != o.Location.Player.Faction.Id {
        finalPrice, gain, purchaseTaxes, salesTaxes = o.processInterfactionTaxes(p, price)
    } else {
        finalPrice, gain, purchaseTaxes, salesTaxes = o.processIntrafactionTaxes(p, price)
    }
    return
}

func (o *Offer) processInterfactionTaxes(p *Player, price int32) (int32, int32, int32, int32) {
    var purchaseTaxesPercent, salesTaxesPercent uint8

    offerFactionRelationship := o.Location.Player.Faction.getFactionRelation(p.Faction)
    proposerFactionRelationship := p.Faction.getFactionRelation(o.Location.Player.Faction)

    if o.Operation == TradeOfferOperationBuy {
        purchaseTaxesPercent = offerFactionRelationship.PurchaseTradeTax
        salesTaxesPercent = proposerFactionRelationship.SalesTradeTax
    } else {
        purchaseTaxesPercent = proposerFactionRelationship.PurchaseTradeTax
        salesTaxesPercent = offerFactionRelationship.SalesTradeTax
    }
    return applyTaxes(price, purchaseTaxesPercent, salesTaxesPercent)
}

func (o *Offer) processIntrafactionTaxes(p *Player, price int32) (int32, int32, int32, int32) {
    purchaseTaxesPercent := uint8(p.Faction.getSettings(FactionSettingsPurchaseTaxes).Value)
    salesTaxesPercent := uint8(p.Faction.getSettings(FactionSettingsSalesTaxes).Value)

    return applyTaxes(price, purchaseTaxesPercent, salesTaxesPercent)
}

func applyTaxes(price int32, purchaseTaxesPercent, salesTaxesPercent uint8) (finalPrice, gain, purchaseTaxes, salesTaxes int32) {
    purchaseTaxes = int32(math.Ceil(float64(price) * (float64(purchaseTaxesPercent) / 100)))
    salesTaxes = int32(math.Ceil(float64(price) * (float64(salesTaxesPercent) / 100)))

    finalPrice = price + purchaseTaxes
    gain = price - salesTaxes

    return
}

func (o *Offer) accept(p *Planet, data map[string]interface{}) {
    panic(NewException("This code should not be executed", nil))
}

func (o *ResourceOffer) accept(p *Planet, data map[string]interface{}) {
    o.Offer.checkProposal(p, data)

    nbLots := uint16(data["nb_lots"].(float64))
    quantity := nbLots * o.LotQuantity
    price := int32(math.Ceil(float64(o.Price) * float64(quantity)))
    if quantity % o.LotQuantity > 0 {
        panic(NewHttpException(400, "There can be no extra resource out of lots", nil))
    }
    if quantity > o.Quantity {
        panic(NewHttpException(400, "You can't demand more lots than available", nil))
    }
    o.Quantity -= quantity
    o.performTransaction(p.Player, price)
    p.Storage.storeResource(o.Resource, int16(quantity))
    p.Storage.update()
    o.addInfluence(p, price)
    o.notifyAcceptation(p.Player, quantity, price)
    if o.Quantity == 0 {
        o.delete()
        return
	} 
    o.update()
}

func (o *Offer) addInfluence(p *Planet, price int32) {
    baseInfluence := math.Ceil(float64(price) / 200)
    modifier := func(pt *PlanetTerritory, influence uint16) {
        pt.EconomicInfluence += influence
    }

    o.Location.addInfluence(o.Location.Player, uint16(baseInfluence), modifier)
    o.Location.addInfluence(p.Player, uint16(math.Ceil(baseInfluence * 0.25)), modifier)
    p.addInfluence(o.Location.Player, uint16(math.Ceil(baseInfluence * 0.75)), modifier)
    p.addInfluence(p.Player, uint16(math.Ceil(baseInfluence / 2)), modifier)
}

func (o *ResourceOffer) notifyAcceptation(p *Player, quantity uint16, price int32) {
    o.Offer.notifyAcceptation()
    o.Location.Player.notify(
        NotificationTypeTrade,
        "notifications.trade.accepted_offer",
        map[string]interface{}{
            "offer_resource": o.Resource,
            "player_id": p.Id,
            "player_pseudo": p.Pseudo,
            "player_faction_id": p.Faction.Id,
            "player_faction_name": p.Faction.Name,
            "quantity": quantity,
            "price": price,
        },
    )
}

func (o *Offer) notifyAcceptation() {
    WsHub.sendTo(o.Location.Player, &WsMessage{ Action: "updateWallet", Data: map[string]uint32{
        "wallet": o.Location.Player.Wallet,
    }})
    WsHub.sendBroadcast(&WsMessage{ Action: "updateTradeOffer", Data: o })
}

func (o *Offer) update() {
    if err := Database.Update(o); err != nil {
        panic(NewHttpException(500, "Offer could not be accepted", err))
    }
}

func (o *Offer) delete() {
	if err := Database.Delete(o); err != nil {
		panic(NewHttpException(500, "Offer could not be deleted", err))
	}
}