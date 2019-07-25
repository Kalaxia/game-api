package api

type(
	Technology struct {
		TableName struct{} `json:"-" sql:"technology__technologies"`

		Id uint16 `json:"id"`
		Name string `json:"name"`
		Children []*Technology `json:"children"`
	}

	ResearchState struct {
		TableName struct{} `json:"-" sql:"technology__research_states"`

		Points uint8 `json:"points"`
		CurrentPoints uint8 `json:"current_points"`
	}
)
