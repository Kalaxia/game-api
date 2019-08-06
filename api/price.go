package api

const(
    PriceTypeResources = "resource"
    PriceTypeMoney = "credits"
    PriceTypePoints = "points"
)

type (
    Price struct {
      Type string `json:"type"`
      Resource string `json:"resource"`
      Amount uint32 `json:"amount"`
    }
)

func (p *Planet) payPrice(prices []Price, quantity uint8) uint8 {
  var points uint8
  for _, price := range prices {
      switch price.Type {
          case PriceTypeMoney:
              if !p.Player.updateWallet(-(int32(price.Amount) * int32(quantity))) {
                  panic(NewHttpException(400, "Not enough money", nil))
              }
              p.Player.update()
              break
          case PriceTypePoints:
              points = uint8(price.Amount) * quantity
              break
          case PriceTypeResources:
              amount := uint16(price.Amount) * uint16(quantity)
              if !p.Storage.hasResource(price.Resource, amount) {
                  panic(NewHttpException(400, "Not enough resources", nil))
              }
              p.Storage.storeResource(price.Resource, -int16(amount))
              break
      }
  }
  return points
}