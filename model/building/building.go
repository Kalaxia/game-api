package model

import(
  mapModel "kalaxia-game-api/model/map"
)

type(
  Building struct {
    TableName struct{} `json:"-" sql:"map__planet_buildings"`

    Name string `json:"name"`
    Type *BuildingType `json:"type"`
    Planet *mapModel.Planet `json:"planet"`
    PlanetId uint16 `json:"-"`
  }
  BuildingPlan struct {
    Name string `json:"name"`
    ParentName string `json:"parent"`
    Type string `json:"type"`
    Duration uint `json:"duration"`
    Price []*Price `json:"price"`
  }
  BuildingPlansData map[string]BuildingPlan

  BuildingType struct {
    Name string `json:"name"`
    Color string `json:"color"`
  }
  Price struct {
    Type string `json:"type"`
    Amount int `json:"type"`
  }
  ResourcePrice struct {
    Price `json:"price"`
    ResourceType string `json:"type"`
  }
  BuildingTypeData struct {
      Color string
  }
  BuildingTypesData map[string]BuildingTypeData
)
