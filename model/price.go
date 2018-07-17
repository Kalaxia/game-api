package model

const(
    PRICE_TYPE_RESOURCE = "resource"
    PRICE_TYPE_MONEY = "credits"
    PRICE_TYPE_POINTS = "points"
)

type (
    Price struct {
      Type string `json:"type"`
      Resource string `json:"resource"`
      Amount uint32 `json:"amount"`
    }
)
