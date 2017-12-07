package manager

import(
  "kalaxia-game-api/model/server"
)

func CreateServer(name string, serverType string, signature string) model.Server {
  server := model.Server{
    Name: name,
    Type: serverType,
    Signature: signature,
  }
  if err := db.Insert(&server); err != nil {
    panic(err)
  }
  return server
}
