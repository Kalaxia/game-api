package api

import(
	"encoding/json"
	"io/ioutil"
)

type(
	Culture struct {
		Identifier string `json:"identifier"`
		Name string `json:"name"`
	}
)

var culturesData map[string]Culture

func InitCultures() {
    defer CatchException(nil)
    culturesDataJSON, err := ioutil.ReadFile("/go/src/kalaxia-game-api/resources/demography_cultures.json")
    if err != nil {
        panic(NewException("Can't open cultures configuration file", err))
    }
    if err := json.Unmarshal(culturesDataJSON, &culturesData); err != nil {
        panic(NewException("Can't read cultures configuration file", err))
    }
}