package manager

import(
  "kalaxia-game-api/model/server"
)

func GetServerBySignature(signature string) *model.Server {
  server := model.Server{Signature: signature}
  if err := db.Model(&server).Where("signature = ?", signature).Select(); err != nil {
    return nil
  }
  return &server
}

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
