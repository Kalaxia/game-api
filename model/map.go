package model

type(
  Map struct {
    TableName struct{} `json:"-" sql:"map__maps"`

    Id uint16 `json:"-"`
    ServerId uint16 `json:"-"`
    Server *Server `json:"-"`
    Systems []System `json:"systems" sql:"-"`
    Size uint16 `json:"size"`
    SectorSize uint16 `json:"sector_size" sql:"-"`
  }
)
