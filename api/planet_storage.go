package api

type Storage struct {
	tableName struct{} `pg:"map__planet_storage"`

	Id uint16 `json:"-"`
	Capacity uint16 `json:"capacity"`
	Resources map[string]uint16 `json:"resources" pg:",notnull,use_zero"`
}

func (s *Storage) hasResource(resource string, quantity uint16) bool {
    q, ok := s.Resources[resource]
    return ok && q >= quantity
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

func (p *Planet) createStorage() {
    storage := &Storage{
        Capacity: 5000,
        Resources: make(map[string]uint16, 0),
    }
    if err := Database.Insert(storage); err != nil {
        panic(NewException("Storage could not be created", err))
    }
    p.Storage = storage
    p.StorageId = storage.Id
    p.update()
}

func (s *Storage) update() {
    if err := Database.Update(s); err != nil {
        panic(NewException("Planet storage could not be updated", err))
    }
}