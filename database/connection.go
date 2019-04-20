package database

import(
  "os"
  "github.com/go-pg/pg"
  "github.com/go-pg/pg/orm"
)

var Connection orm.DB

func init() {
  // Cast go-pg instance to PgConnection to allow mocks
  Connection = pg.Connect(&pg.Options{
    Network: "tcp",
    Addr: os.Getenv("POSTGRES_HOST") + ":5432",
    User: os.Getenv("POSTGRES_USER"),
    Password: os.Getenv("POSTGRES_PASSWORD"),
    Database: os.Getenv("POSTGRES_DB"),
  })
}
