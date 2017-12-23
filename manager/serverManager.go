package manager

import(
  "kalaxia-game-api/database"
  "kalaxia-game-api/model/server"
)

func GetServerBySignature(signature string) *model.Server {
  server := model.Server{Signature: signature}
  if err := database.Connection.Model(&server).Where("signature = ?", signature).Select(); err != nil {
    return nil
  }
  return &server
}

func CreateServer(name string, serverType string, signature string, mapSize uint16) model.Server {
  server := model.Server{
    Name: name,
    Type: serverType,
    Signature: signature,
  }
  if err := database.Connection.Insert(&server); err != nil {
    panic(err)
  }
  GenerateMap(&server, mapSize)
  return server
}
