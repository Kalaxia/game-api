package handler

import(
    "net/http"
    "kalaxia-game-api/manager"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/security"
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
        player := manager.GetPlayer(security.GetJwtPlayerId(claims.(jwt.MapClaims)), true)
        if player == nil {
            panic(exception.NewHttpException(http.StatusInternalServerError, "Unavailable player account", nil))
        }
        context.Set(req, "player", player)
        next(w, req)
    })
}
