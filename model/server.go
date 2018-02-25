package model

import "fmt"

type(
  Server struct {
    Id uint16 `json:"id"`
    Name string `json:"name" sql:"type:varchar(100);not null;unique"`
    Type string `json:"type" sql:"type:varchar(20);not null"`
    Signature string `json:"_" sql:"type:varchar(125);not null;unique"`
  }
)

func (s Server) String() string {
  return fmt.Sprintf(
    "Server<Id=%d Name=%q Type=%q Signature=%q>",
    s.Id,
    s.Name,
    s.Type,
    s.Signature,
  )
}
