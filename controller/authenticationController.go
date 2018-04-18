package controller

import (
    "net/http"
    "os"
    "time"
    "kalaxia-game-api/exception"
    "kalaxia-game-api/manager"
    "kalaxia-game-api/model"
    "kalaxia-game-api/security"
    "kalaxia-game-api/utils"
    "github.com/dgrijalva/jwt-go"
)

/**
 * This method receives crypted data from the portal, containing user credentials
 * If the user is not registered, we create a player account.
 * In any case, a JWT is returned, used by the player to authenticate
 */
func Authenticate(w http.ResponseWriter, r *http.Request) {
    data := utils.DecodeJsonRequest(r)
    server := manager.GetServerBySignature(data["signature"].(string))
    if server == nil {
        panic(exception.NewHttpException(404, "Server not found", nil))
    }
    player := manager.GetPlayerByUsername(data["username"].(string), server)
    if player == nil {
        player = manager.CreatePlayer(data["username"].(string), server)
    }
    if server.Id != player.Server.Id {
        panic(exception.NewHttpException(http.StatusBadRequest, "Invalid server data", nil))
    }
    token := getNewJWT(player)
    cipherData, key, iv := security.Encrypt([]byte(token))
    w.Header().Set("Application-Key", key)
    w.Header().Set("Application-Iv", iv)
    w.Write(cipherData)
}

func getNewJWT(player *model.Player) string {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "id": player.Id,
        "pseudo": player.Pseudo,
        "server_id": player.Server.Id,
        "created_at": time.Now().Format(time.RFC3339),
    })
    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        panic(exception.NewHttpException(http.StatusInternalServerError, "JWT creation failed", err))
    }
    return tokenString
}
