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
    SystemId uint16 `json:"-"`
    System *System `json:"system"`
    OrbitId uint16 `json:"-"`
    Orbit *SystemOrbit `json:"orbit"`
    PlayerId uint16 `json:"-"`
    Player *Player `json:"player"`
    Resources []PlanetResource `json:"resources"`
    StorageId uint16 `json:"-"`
    Storage *Storage `json:"storage"`
    Relations []DiplomaticRelation `json:"relations" sql:"-"`
    Buildings []Building `json:"buildings" sql:"-"`
    NbBuildings uint8 `json:"nb_buildings" sql:"-"`
    AvailableBuildings []BuildingPlan `json:"available_buildings" sql:"-"`
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

      Id uint16 `json:"id"`
      Capacity uint16 `json:"capacity"`
      Resources map[string]uint16 `json:"resources"`
  }
  PlanetData struct {
    Resources map[string]uint8
  }
  PlanetsData map[string]PlanetData
)
