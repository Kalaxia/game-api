package tradeManager

import(
    "kalaxia-game-api/database"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/model"
    "kalaxia-game-api/manager"
    "kalaxia-game-api/manager/ship"
    "math"
    "time"
)

func GetOffer(id uint32) *model.ResourceOffer {
    offer := &model.ResourceOffer{}
    if err := database.Connection.Model(offer).Column("Location.Player.Faction", "Location.System", "Location.Storage").Where("resource_offer.id = ?", id).Select(); err != nil {
        panic(exception.NewHttpException(404, "Offer not found", err))
    }
    offer.Type = "resources"
    return offer
}
//
// func GetPlanetOffers(location *model.Planet) []model.OfferInterface {
//     offers := make([]model.Offer, 0)
//     if err := database.Connection.Model(&offers).Where("location_id = ?", location.Id).Select(); err != nil {
//         panic(exception.NewHttpException(500, "Planet offers could not be retrieved", err))
//     }
//     return offers
// }
//
// func GetDestinationOffers(location *model.Planet) []model.OfferInterface {
//     offers := make([]model.Offer, 0)
//     if err := database.Connection.Model(&offers).Where("destination_id = ?", location.Id).Select(); err != nil {
//         panic(exception.NewHttpException(500, "Planet offers could not be retrieved", err))
//     }
//     return offers
// }

func CancelOffer(offer *model.ResourceOffer, player *model.Player) {
    if offer.Location.Player.Id != player.Id {
        panic(exception.NewHttpException(403, "You do not own this offer", nil))
    }
    manager.UpdateStorageResource(offer.Location.Storage, offer.Resource, int16(offer.Quantity))
    manager.UpdatePlanetStorage(offer.Location)

    if err := database.Connection.Delete(offer); err != nil {
        panic(exception.NewHttpException(500, "Could not delete offer", err))
    }
}

func SearchOffers(data map[string]interface{}) []model.OfferInterface {
    offers := make([]model.OfferInterface, 0)

    operation := data["operation"].(string)

    resourceOffers := make([]*model.ResourceOffer, 0)
    if err := database.Connection.Model(&resourceOffers).Column("Location.Player.Faction", "Location.System").Where("operation = ?", operation).Select(); err != nil {
        panic(exception.NewHttpException(500, "Resource offers could not be retrieved", err))
    }
    // shipOffers := make([]*model.ShipOffer, 0)
    // if err := database.Connection.Model(&shipOffers).Where("operation = ?", operation).Select(); err != nil {
    //     panic(exception.NewHttpException(500, "Ship offers could not be retrieved", err))
    // }
    // modelOffers := make([]*model.ModelOffer, 0)
    // if err := database.Connection.Model(&modelOffers).Where("operation = ?", operation).Select(); err != nil {
    //     panic(exception.NewHttpException(500, "Model offers could not be retrieved", err))
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

func CreateOffer(planet *model.Planet, data map[string]interface{}) model.OfferInterface {
    var offer model.OfferInterface
    switch goodType := data["good_type"].(string); goodType {
        case "resources":
            offer = createResourceOffer(planet, data);
        case "ships":
            offer = createShipOffer(planet, data);
        case "models":
            offer = createModelOffer(planet, data);
    }
    return offer
}

func createResourceOffer(location *model.Planet, data map[string]interface{}) *model.ResourceOffer {
    offer := &model.ResourceOffer{
        Resource: data["resource"].(string),
        Quantity: uint16(data["quantity"].(float64)),
        LotQuantity: uint16(data["lot_quantity"].(float64)),
        Price: float32(data["price"].(float64)),
    }
    offer.Type = "resources"
    offer.Operation = data["operation"].(string)
    offer.LocationId = location.Id
    offer.Location = location
    offer.CreatedAt = time.Now()
    if offer.Quantity < offer.LotQuantity {
        panic(exception.NewHttpException(400, "Lot quantity cannot be lesser than total quantity", nil))
    }
    if !manager.UpdateStorageResource(location.Storage, offer.Resource, -int16(offer.Quantity)) {
        panic(exception.NewHttpException(400, "Not enough resources in storage", nil))
    }
    manager.UpdatePlanetStorage(location)
    if err := database.Connection.Insert(offer); err != nil {
        panic(exception.NewHttpException(500, "Resource offer could not be created", err))
    }
    return offer
}

func createShipOffer(location *model.Planet, data map[string]interface{}) *model.ShipOffer {
    shipModel := shipManager.GetShipModel(location.Player.Id, data["model"].(uint32))

    offer := &model.ShipOffer{
        Quantity: data["quantity"].(uint16),
        LotQuantity: data["lot_quantity"].(uint16),
    }
    offer.Operation = data["operation"].(string)
    offer.LocationId = location.Id
    offer.Location = location
    offer.CreatedAt = time.Now()
    offer.Price = data["price"].(uint16)
    offer.Model = shipModel
    offer.ModelId = shipModel.Id
    if err := database.Connection.Insert(offer); err != nil {
        panic(exception.NewHttpException(500, "Ship offer could not be created", err))
    }
    return offer
}

func createModelOffer(location *model.Planet, data map[string]interface{}) *model.ModelOffer {
    shipModel := shipManager.GetShipModel(location.Player.Id, data["model"].(uint32))

    offer := &model.ModelOffer{
        Price: data["price"].(uint16),
        ModelId: shipModel.Id,
        Model: shipModel,
    }
    offer.Operation = data["operation"].(string)
    offer.LocationId = location.Id
    offer.Location = location
    offer.CreatedAt = time.Now()
    if err := database.Connection.Insert(offer); err != nil {
        panic(exception.NewHttpException(500, "Ship model offer could not be created", err))
    }
    return offer
}

func AcceptOffer(offerId uint32, planet *model.Planet, nbLots uint16) {
    offer := GetOffer(offerId)
    if offer == nil {
        panic(exception.NewHttpException(404, "Offer not found", nil))
    }
    if offer.Location.Player.Id == planet.Player.Id {
        panic(exception.NewHttpException(400, "You can't accept your own offers", nil))
    }

    quantity := nbLots * offer.LotQuantity
    price := int32(math.Ceil(float64(offer.Price) * float64(quantity)))

    if quantity % offer.LotQuantity > 0 {
        panic(exception.NewHttpException(400, "There can be no extra resource out of lots", nil))
    }
    if quantity > offer.Quantity {
        panic(exception.NewHttpException(400, "You can't demand more lots than available", nil))
    }
    if !manager.UpdatePlayerWallet(planet.Player, -price) {
        panic(exception.NewHttpException(400, "Not enough money", nil))
    }
    manager.UpdatePlayerWallet(offer.Location.Player, price)
    manager.UpdatePlayer(planet.Player)
    manager.UpdatePlayer(offer.Location.Player)

    manager.UpdateStorageResource(planet.Storage, offer.Resource, int16(quantity))
    manager.UpdatePlanetStorage(planet)

    offer.Quantity -= quantity
    if offer.Quantity == 0 {
        if err := database.Connection.Delete(offer); err != nil {
            panic(exception.NewHttpException(500, "Offer could not be deleted", err))
        }
        return
    } 
    if err := database.Connection.Update(offer); err != nil {
        panic(exception.NewHttpException(500, "Offer could not be accepted", err))
    }
}
