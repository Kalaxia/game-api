package api

import(
	"time"
)

type(
	FactionSettings struct {
		TableName struct{} `json:"-" sql:"faction__settings"`

		Id uint16 `json:"-"`
		Faction *Faction `json:"faction"`
		FactionId uint16 `json:"-"`
		Name string `json:"name"`
		Value int `json:"value"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

const(
	FactionRegimeDemocracy = "democracy"
	FactionSettingsRegime = "regime"
	FactionSettingsPlanetTaxes = "planet_taxes"
)

func (f *Faction) getSettings(name string) *FactionSettings {
	for _, s := range f.Settings {
		if s.Name == name {
			return s
		}
	}
	var settings *FactionSettings
	if err := Database.Model(&settings).Where("faction_id = ?", f.Id).Where("name = ?", name).Select(); err != nil {
		panic(NewException("Could not retrieve faction settings", err))
	}
	return settings
}