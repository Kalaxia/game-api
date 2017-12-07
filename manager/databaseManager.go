package manager

import(
  "fmt"
  "os"
  "github.com/go-pg/pg"
)

var db *pg.DB

func init() {
  options := &pg.Options{
    Network: "tcp",
    Addr: "kalaxia_postgresql:5432",
    User: os.Getenv("POSTGRES_USER"),
    Password: os.Getenv("POSTGRES_PASSWORD"),
    Database: os.Getenv("POSTGRES_DB"),
  }
  db = pg.Connect(options)
  fmt.Println("Database connection initialized")
}
