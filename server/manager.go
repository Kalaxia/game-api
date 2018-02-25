package server

import(
    "kalaxia-game-api/database"
)

func GetServerBySignature(signature string) *Server {
    server := Server{Signature: signature}
    if err := database.Connection.Model(&server).Where("signature = ?", signature).Select(); err != nil {
        return nil
    }
    return &server
}

func CreateServer(name string, serverType string, signature string) *Server {
    server := &Server{
        Name: name,
        Type: serverType,
        Signature: signature,
    }
    if err := database.Connection.Insert(server); err != nil {
        panic(err)
    }
    return server
}
