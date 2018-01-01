package handler

import(
  "net/http"
  "time"
  "kalaxia-game-api/manager"
  "github.com/dgrijalva/jwt-go"
  "github.com/gorilla/context"
)

func AuthorizationHandler(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        claims := context.Get(req, "jwt")
        // This case means that the JWT was not mandatory on this call
        if claims == nil {
          next(w, req)
          return
        }
        data := claims.(jwt.MapClaims)
        created_at, _ := time.Parse(time.RFC3339, data["created_at"].(string))
        if created_at.Add(time.Hour * time.Duration(2)).Before(time.Now()) {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("Expired JWT"))
        }
        player := manager.GetPlayer(uint16(data["id"].(float64)))
        if player != nil {
            context.Set(req, "player", player)
            next(w, req)
        } else {
            w.WriteHeader(http.StatusInternalServerError)
            w.Write([]byte("Unavailable player account"))
        }
    })
}
