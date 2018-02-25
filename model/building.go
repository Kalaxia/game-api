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
    Type *BuildingType `json:"type" sql:"-"`
    TypeName string `json:"-" sql:"type"`
    Planet *Planet `json:"planet"`
    PlanetId uint16 `json:"-"`
    Status string `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
  }
  BuildingPlan struct {
    Name string `json:"name"`
    ParentName string `json:"parent"`
    Type string `json:"type"`
    Duration uint `json:"duration"`
    Price []Price `json:"price"`
  }
  BuildingPlansData map[string]BuildingPlan

  BuildingType struct {
    Name string `json:"name"`
    Color string `json:"color"`
  }
  Price struct {
    Type string `json:"type"`
    ResourceType string `json:"resource_type,omitempty"`
    Amount uint `json:"amount"`
  }
  BuildingTypeData struct {
      Color string `json:"color"`
  }
  BuildingTypesData map[string]BuildingTypeData
)
