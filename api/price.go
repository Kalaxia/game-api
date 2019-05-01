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
