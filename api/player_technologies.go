package api

type(
	Technology struct {
		TableName struct{} `json:"-" sql:"technology__technologies"`

		Id uint16 `json:"id"`
		Name string `json:"name"`
		Children []*Technology `json:"children"`
	}
)
