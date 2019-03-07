package model

const PlanetTypeArtic = "arctic"
const PlanetTypeDesert = "desert"
const PlanetTypeOceanic = "oceanic"
const PlanetTypeTemperate = "temperate"
const PlanetTypeTropical = "tropical"
const PlanetTypeRocky = "rocky"
const PlanetTypeVolcanic = "volcanic"

type(
  Planet struct {
    TableName struct{} `json:"-" sql:"map__planets"`

    Id uint16 `json:"id"`
    Name string `json:"name"`
    Type string `json:"type"`
    Population uint `json:"population"`
    SystemId uint16 `json:"-"`
    System *System `json:"system"`
    OrbitId uint16 `json:"-"`
    Orbit *SystemOrbit `json:"orbit"`
    PlayerId uint16 `json:"-"`
    Player *Player `json:"player"`
    Resources []PlanetResource `json:"resources"`
    StorageId uint16 `json:"-"`
    Storage *Storage `json:"storage"`
    SettingsId uint16 `json:"-"`
    Settings *PlanetSettings `json:"settings"`
    Relations []DiplomaticRelation `json:"relations"`
    Buildings []Building `json:"buildings"`
    NbBuildings uint8 `json:"nb_buildings" sql:"-"`
    AvailableBuildings []BuildingPlan `json:"available_buildings" sql:"-"`
  }
  PlanetSettings struct {
      TableName struct{} `json:"-" sql:"map__planet_settings"`

      Id uint16 `json:"-"`
      ServicesPoints uint8 `json:"services_points" sql:",notnull"`
      BuildingPoints uint8 `json:"building_points" sql:",notnull"`
      MilitaryPoints uint8 `json:"military_points" sql:",notnull"`
      ResearchPoints uint8 `json:"research_points" sql:",notnull"`
  }

  PlanetResource struct {
    TableName struct{} `json:"-" sql:"map__planet_resources"`

    Name string `json:"name"`
    Density uint8 `json:"density"`
    PlanetId uint16 `json:"-"`
    Planet *Planet `json:"planet"`
  }
  Storage struct {
      TableName struct{} `json:"-" sql:"map__planet_storage"`

      Id uint16 `json:"-"`
      Capacity uint16 `json:"capacity"`
      Resources map[string]uint16 `json:"resources"`
  }
  PlanetData struct {
    Picto string
    Image string
    Resources map[string]uint8
  }
  PlanetsData map[string]PlanetData
)

func (s *Storage) HasResource(resource string, quantity uint16) bool {
    if rQuantity, ok := s.Resources[resource]; !ok || rQuantity < quantity {
        return false
    }
    return true
}

func (s *Storage) SubstractResource(resource string, quantity uint16) {
    if _, ok := s.Resources[resource]; ok {
        s.Resources[resource] -= quantity
    }
}
