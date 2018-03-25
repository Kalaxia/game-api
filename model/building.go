package model

import(
    "time"
)

const BuildingStatusConstructing = "constructing"
const BuildingStatusOperational = "operational"
const BuildingStatusDestroying = "destroying"

type(
  Building struct {
    TableName struct{} `json:"-" sql:"map__planet_buildings"`

    Id uint32 `json:"id"`
    Name string `json:"name"`
    Type string `json:"type" sql:"type"`
    Planet *Planet `json:"planet"`
    PlanetId uint16 `json:"-"`
    Status string `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    BuiltAt time.Time `json:"built_at"`
  }
  BuildingPlan struct {
    Name string `json:"name"`
    ParentName string `json:"parent"`
    Type string `json:"type"`
    Duration uint `json:"duration"`
    Price []Price `json:"price"`
  }
  BuildingPlansData map[string]BuildingPlan

  Price struct {
    Type string `json:"type"`
    Amount uint `json:"amount"`
  }
)
