package model

import(
    "time"
)

type(
    OfferInterface interface {
        getTotalPrice() float32
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

func (o *ResourceOffer) getTotalPrice() float32 {
    return o.Price
}

func (o *ShipOffer) getTotalPrice() float32 {
    return float32(o.Price)
}

func (o *ModelOffer) getTotalPrice() float32 {
    return float32(o.Price)
}
