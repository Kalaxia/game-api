package player

import(
  "time"
  "kalaxia-game-api/model/server"
  orm "github.com/go-pg/pg"
)

type(
  Player struct {
    Id int16 `json:"id"`
    Username string `json:"_" sql:"type:varchar(180);not null;unique"`
    Pseudo string `json:"pseudo" sql:"type:varchar(180);not null;unique"`
    Server *server.Server `json"_"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
  }
  Players []Player
)

func (p *Player) beforeInsert(db orm.DB) error {
  p.CreatedAt = time.Now()
  p.UpdatedAt = time.Now()
  return nil
}

func (p *Player) beforeUpdate(db orm.DB) error {
  p.UpdatedAt = time.Now()
}
