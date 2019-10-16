package api

type(
	PlanetHangarGroup struct{
		tableName struct{} `json:"-" pg:"map__planet_hangar_groups"`

		Id uint32 `json:"id"`
		LocationId uint16 `json:"-"`
		Location *Planet `json:"location"`
		ModelId uint `json:"-"`
		Model *ShipModel `json:"model"`
		Quantity uint16 `json:"quantity"`
	}
)

func (p *Planet) addShips(sm *ShipModel, quantity uint8) {
	if hg := p.getHangarGroup(sm); hg != nil {
		hg.Quantity += uint16(quantity)
		hg.update()
		return
	}
	p.createHangarGroup(sm, uint16(quantity))
}

func (p *Planet) createHangarGroup(sm *ShipModel, quantity uint16) *PlanetHangarGroup {
	hg := &PlanetHangarGroup{
		Location: p,
		LocationId: p.Id,
		Model: sm,
		ModelId: sm.Id,
		Quantity: quantity,
	}
	if err := Database.Insert(hg); err != nil {
		panic(NewException("Could not create planet hangar group", err))
	}
	return hg
}

func (p *Planet) getHangarGroup(sm *ShipModel) *PlanetHangarGroup {
	hg := &PlanetHangarGroup{
		Location: p,
		Model: sm,
	}
	if err := Database.Model(hg).Where("location_id = ?", p.Id).Where("model_id = ?", sm.Id).Select(); err != nil {
		return nil
	}
	return hg
}

func (p *Planet) getHangarGroups() []PlanetHangarGroup {
    groups := make([]PlanetHangarGroup, 0)

    if err := Database.Model(&groups).Relation("Model").Where("location_id = ?", p.Id).Select(); err != nil {
		panic(NewHttpException(404, "planet not found", err))
    }
    return groups
}

func (hg *PlanetHangarGroup) update() {
	if err := Database.Update(hg); err != nil {
		panic(NewException("Could not update planet hangar group", err))
	}
}

func (hg *PlanetHangarGroup) delete() {
	if err := Database.Delete(hg); err != nil {
		panic(NewException("Could not delete planet hangar group", err))
	}
}