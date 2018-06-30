package model

type (
    Price struct {
      Type string `json:"type"`
      Resource string `json:"resource"`
      Amount uint `json:"amount"`
    }
)
