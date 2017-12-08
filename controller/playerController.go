package controller

import (
  "net/http"

  "github.com/dgrijalva/jwt-go"
  "github.com/gorilla/context"
)

func GetCurrentPlayer(w http.ResponseWriter, r *http.Request) {
  claims := context.Get(r, "decoded")
  data := claims.(jwt.MapClaims)
  w.Write([]byte("Hello " + data["pseudo"].(string) + " !"))
}
