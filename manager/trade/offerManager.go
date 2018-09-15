package tradeManager

import(
    "kalaxia-game-api/database"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/model"
    "kalaxia-game-api/manager/ship"
    "time"
)

// func GetOffer(id uint32) model.OfferInterface {
//     offer := model.Offer{ Id: id }
//     if err := database.Connection.Select(&offer); err != nil {
//         panic(exception.NewHttpException(404, "Offer not found", err))
//     }
//     return offer
// }
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
    offer.Operation = data["operation"].(string)
    offer.LocationId = location.Id
    offer.Location = location
    offer.CreatedAt = time.Now()
    if offer.Quantity < offer.LotQuantity {
        panic(exception.NewHttpException(400, "Lot quantity cannot be lesser than total quantity", nil))
    }
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

// func AcceptOffer(offerId uint32, player *model.Player) {
//     offer := GetOffer(offerId)
//
//     if offer.Price > player.Wallet {
//         panic(exception.NewHttpException(400, "Not enough money", nil))
//     }
//     if err := database.Connection.Update(offer); err != nil {
//         panic(exception.NewHttpException(500, "Offer could not be accepted", err))
//     }
// }
