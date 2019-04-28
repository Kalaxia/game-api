package handler

import(
    "net/http"
    "time"
    "kalaxia-game-api/manager"
    "kalaxia-game-api/exception"
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
        createdAt, _ := time.Parse(time.RFC3339, data["created_at"].(string))
        if createdAt.Add(time.Hour * time.Duration(2)).Before(time.Now()) {
            panic(exception.NewHttpException(http.StatusUnauthorized, "Expired JWT", nil))
        }
        player := manager.GetPlayer(uint16(data["id"].(float64)), true)
        if player == nil {
            panic(exception.NewHttpException(http.StatusInternalServerError, "Unavailable player account", nil))
        }
        context.Set(req, "player", player)
        next(w, req)
    })
}
