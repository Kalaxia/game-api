package model

type(
  Faction struct {
    TableName struct{} `json:"-" sql:"faction__factions"`

    Id uint16 `json:"id"`
    Name string `json:"name"`
    Slug string `json:"slug"`
    Description string `json:"description"`
    Colors map[string]string `json:"colors"`
    Banner string `json:"banner"`
    ServerId uint16 `json:"-"`
    Server *Server `json:"-"`
    Relations []*FactionRelation `json:"relations"`
  }
)
