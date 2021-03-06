package api

import(
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "net/http"
    "strconv"
)

const PlanetTypeArctic = "arctic"
const PlanetTypeDesert = "desert"
const PlanetTypeOceanic = "oceanic"
const PlanetTypeTemperate = "temperate"
const PlanetTypeTropical = "tropical"
const PlanetTypeRocky = "rocky"
const PlanetTypeVolcanic = "volcanic"

type(
	Planet struct {
		tableName struct{} `pg:"map__planets"`
	
		Id uint16 `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
        Population uint `json:"population" pg:",use_zero,notnull"`
        PopulationGrowth float64 `json:"population_growth" pg:"-"`
        PublicOrder uint8 `json:"public_order" pg:",use_zero,notnull"`
        PublicOrderGrowth int8 `json:"public_order_growth" pg:"-"`
        TaxRate uint8 `json:"tax_rate" pg:",use_zero,notnull"`
        FactionId uint16 `json:"-"`
        Faction *Faction `json:"faction"`
		SystemId uint16 `json:"-"`
		System *System `json:"system"`
		OrbitId uint16 `json:"-"`
		Orbit *SystemOrbit `json:"orbit"`
		PlayerId uint16 `json:"-"`
		Player *Player `json:"player"`
        Resources []PlanetResource `json:"resources"`
        ResourcesProduction map[string]*ResourceProduction `json:"resources_production" pg:"-"`
		StorageId uint16 `json:"-"`
		Storage *Storage `json:"storage"`
		SettingsId uint16 `json:"-"`
		Settings *PlanetSettings `json:"settings"`
		Relations []*DiplomaticRelation `json:"relations"`
		Buildings []*Building `json:"buildings"`
		NbBuildings uint8 `json:"nb_buildings" pg:"-"`
        AvailableBuildings []BuildingPlan `json:"available_buildings" pg:"-"`
        Territories []*PlanetTerritory `json:"territories"`
	}
	
	PlanetResource struct {
		tableName struct{} `pg:"map__planet_resources"`
	
		Name string `json:"name"`
		Density uint8 `json:"density"`
		PlanetId uint16 `json:"-"`
		Planet *Planet `json:"planet"`
    }

	PlanetData struct {
	  Picto string
	  Image string
	  Resources map[string]uint8
	}
	PlanetsData map[string]PlanetData
)

func GetPlanet(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)
    planet := getPlanet(uint16(id))
    if planet.PlayerId == player.Id {
        planet = player.getPlanet(planet.Id)
    }
    SendJsonResponse(w, 200, planet)
}

func (p *Planet) changeOwner(player *Player) {
	if p.Player != nil && p.Player.CurrentPlanetId == p.Id {
        p.Player.relocate()
	}
	p.Player = player
    p.PlayerId = player.Id
    p.update()
}

func (p *Planet) update() {
    if err := Database.Update(p); err != nil {
        panic(NewException("Planet could not be updated", err))
    }
}

func getPlanet(id uint16) *Planet {
    planet := &Planet{}
    if err := Database.
        Model(planet).
        Relation("Player.Faction").
        Relation("Settings").
        Relation("Relations.Player.Faction").
        Relation("Relations.Faction").
        Relation("Resources").
        Relation("System.Map").
        Relation("Storage").
        Where("planet.id = ?", id).
        Select(); err != nil {
            return nil
    }
    return planet
}

func (s *System) getPlanets() []Planet {
    planets := make([]Planet, 0)
    if err := Database.
        Model(&planets).
        Relation("Orbit").
        Relation("Player.Faction").
        Where("planet.system_id = ?", s.Id).
        Select(); err != nil {
            panic(NewHttpException(404, "System not found", err))
    }
    return planets
}

func (p *Player) getPlanets() []Planet {
    planets := make([]Planet, 0)
    if err := Database.
        Model(&planets).
        Relation("Player.Faction").
        Relation("System").
        Relation("Resources").
        Relation("Settings").
        Where("planet.player_id = ?", p.Id).
        Select(); err != nil {
            panic(NewHttpException(404, "Player not found", err))
    }
    return planets
}

func (p *Player) getPlanet(id uint16) *Planet {
    planet := &Planet{}
    if err := Database.
        Model(planet).
        Relation("Player.Faction").
        Relation("Settings").
        Relation("Relations.Player.Faction").
        Relation("Relations.Faction").
        Relation("Resources").
        Relation("System").
        Relation("Storage").
        Where("planet.id = ?", id).
        Where("Player.id = ?", p.Id).
        Select(); err != nil {
            panic(NewHttpException(404, "Planet not found", err))
    }
    planet.getOwnerData()
    return planet
}

func (p *Planet) getOwnerData() {
    p.Buildings, p.AvailableBuildings = p.getBuildings()
    p.NbBuildings = 7
    p.ResourcesProduction = p.getProducedResources()
    p.PopulationGrowth = p.calculatePopulationGrowth()
    p.PublicOrderGrowth = p.calculatePublicOrderGrowth()
}

func (p *Planet) getResource(name string) *PlanetResource {
    for _, r := range p.Resources {
        if r.Name == name {
            return &r
        }
    }
    return nil
}

func (p *Planet) getFaction() *Faction {
    if p.Player != nil {
        return p.Player.Faction
    }
    if p.Faction != nil {
        return p.Faction
    }
    return nil
}