package api

import (
    "github.com/gorilla/context"
    "net/http"
    "os"
    "time"
    "github.com/dgrijalva/jwt-go"
)

/**
 * This method receives crypted data from the portal, containing user credentials
 * If the user is not registered, we create a player account.
 * In any case, a JWT is returned, used by the player to authenticate
 */
func Authenticate(w http.ResponseWriter, r *http.Request) {
    data := DecodeJsonRequest(r)
    server := getServerBySignature(data["signature"].(string))
    if server == nil {
        panic(NewHttpException(404, "Server not found", nil))
    }
    player := server.getPlayerByUsername(data["username"].(string))
    if player == nil {
        player = server.createPlayer(data["username"].(string))
    }
    token := player.getNewJWT()
    cipherData, key, iv := Encrypt([]byte(token))
    w.Header().Set("Application-Key", key)
    w.Header().Set("Application-Iv", iv)
    w.Write(cipherData)
}

func (p *Player) getNewJWT() string {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "id": p.Id,
        "pseudo": p.Pseudo,
        "server_id": p.Server.Id,
        "created_at": time.Now().Format(time.RFC3339),
    })
    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        panic(NewHttpException(http.StatusInternalServerError, "JWT creation failed", err))
    }
    return tokenString
}

func AuthorizationHandler(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        claims := context.Get(req, "jwt")
        // This case means that the JWT was not mandatory on this call
        if claims == nil {
            next(w, req)
            return
        }
        if player := getJwtPlayer(claims.(jwt.MapClaims)); player != nil {
			context.Set(req, "player", player)
            next(w, req)
            return
		}
		panic(NewHttpException(http.StatusInternalServerError, "Unavailable player account", nil))
    })
}
