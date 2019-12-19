package api

import(
	"time"
)

type(
	FactionSettings struct {
		tableName struct{} `json:"-" pg:"faction__settings"`

		Id uint16 `json:"-"`
		Faction *Faction `json:"faction"`
		FactionId uint16 `json:"-"`
		IsPublic bool `json:"is_public" pg:",notnull,use_zero"`
		Name string `json:"name"`
		Value int `json:"value"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

const(
	FactionRegimeDemocracy = "democracy"
	FactionSettingsRegime = "regime"
	FactionSettingsPlanetTaxes = "planet_taxes"
	FactionSettingsMotionDuration = "motion_duration"
)

func (f *Faction) getSettings(name string) *FactionSettings {
	for _, s := range f.Settings {
		if s.Name == name {
			return s
		}
	}
	settings := &FactionSettings{}
	if err := Database.Model(settings).Where("faction_id = ?", f.Id).Where("name = ?", name).Select(); err != nil {
		panic(NewException("Could not retrieve faction settings", err))
	}
	return settings
}

func (f *Faction) getAllSettings(publicOnly bool) []*FactionSettings {
	settings := make([]*FactionSettings, 0)

	query := Database.Model(&settings).Where("faction_id = ?", f.Id)

	if publicOnly {
		query.Where("is_public = ?", true)
	}
	if err := query.Select(); err != nil {
		panic(NewException("Could not retrieve all faction settings", err))
	}
	return settings
}