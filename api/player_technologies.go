package api

type(
	Technology struct {
		tableName struct{} `pg:"technology__technologies"`

		Id uint16 `json:"id"`
		Name string `json:"name"`
		Children []*Technology `json:"children"`
	}
)
