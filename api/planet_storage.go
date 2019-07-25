package api

type Storage struct {
	TableName struct{} `json:"-" sql:"map__planet_storage"`

	Id uint16 `json:"-"`
	Capacity uint16 `json:"capacity"`
	Resources map[string]uint16 `json:"resources"`
}

func (s *Storage) hasResource(resource string, quantity uint16) bool {
    q, ok := s.Resources[resource]
    return ok && q >= quantity
}

func (s *Storage) storeResourceProduction(planet *Planet) {
    for _, resource := range planet.Resources {
        s.storeResource(resource.Name, int16(resource.Density) * 10)
    }
}

func (s *Storage) storeResource(resource string, quantity int16) bool {
    var currentStock uint16
    var newStock int16
    var isset bool
    if currentStock, isset = s.Resources[resource]; !isset {
        currentStock = 0
    }
    if newStock = int16(currentStock) + quantity; newStock > int16(s.Capacity) {
        newStock = int16(s.Capacity)
    }
    if newStock < 0 {
        return false
    }
    s.Resources[resource] = uint16(newStock)
    return true
}

func (s *Storage) update() {
    if err := Database.Update(s); err != nil {
        panic(NewException("Planet storage could not be updated", err))
    }
}