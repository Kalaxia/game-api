package database

import(
  "fmt"
  "os"
  "github.com/go-pg/pg"
)

var Connection *pg.DB

func init() {
  options := &pg.Options{
    Network: "tcp",
    Addr: os.Getenv("POSTGRES_HOST") + ":5432",
    User: os.Getenv("POSTGRES_USER"),
    Password: os.Getenv("POSTGRES_PASSWORD"),
    Database: os.Getenv("POSTGRES_DB"),
  }
  Connection = pg.Connect(options)
  fmt.Println("Database connection initialized")
}
