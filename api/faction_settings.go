package api

import(
	"time"
)

type(
	FactionSettings struct {
		tableName struct{} `pg:"faction__settings"`

		Id uint16 `json:"-"`
		Faction *Faction `json:"faction"`
		FactionId uint16 `json:"-"`
		IsPublic bool `json:"is_public" pg:",notnull,use_zero"`
		Name string `json:"name"`
		Value int `json:"value" pg:",notnull,use_zero"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

const(
	FactionRegimeDemocracy = 1
	FactionSettingsRegime = "regime"
	FactionSettingsMotionDuration = "motion_duration"
	FactionSettingsPlanetTaxes = "planet_taxes"
	FactionSettingsPurchaseTaxes = "purchase_taxes"
	FactionSettingsSalesTaxes = "sales_taxes"
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

func createFactionSettings(factions []*Faction) {
	defaultSettings := map[string]map[string]int{
		FactionSettingsRegime: map[string]int{
			"is_public": 1,
			"value": FactionRegimeDemocracy,
		},
		FactionSettingsMotionDuration: map[string]int{
			"is_public": 0,
			"value": 6,
		},
		FactionSettingsPlanetTaxes: map[string]int{
			"is_public": 0,
			"value": 100,
		},
		FactionSettingsPurchaseTaxes: map[string]int{
			"is_public": 0,
			"value": 4,
		},
		FactionSettingsSalesTaxes: map[string]int{
			"is_public": 0,
			"value": 2,
		},
	} 

	for _, f := range factions {
		for name, data := range defaultSettings {
			s := &FactionSettings{
				Faction: f,
				FactionId: f.Id,
				Name: name,
				Value: data["value"],
				IsPublic: data["is_public"] == 1,
			}
			if err := Database.Insert(s); err != nil {
				panic(NewException("Could not create faction settings", err))
			}
			f.Settings = append(f.Settings, s)
		}
	}
}
