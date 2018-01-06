package controller

import (
  "fmt"
  "net/http"
  "encoding/json"
	"io"
	"io/ioutil"
  "os"
  "time"
  "kalaxia-game-api/manager"
  "kalaxia-game-api/model/player"
  "kalaxia-game-api/security"
  "github.com/dgrijalva/jwt-go"
)

/**
 * This method receives crypted data from the portal, containing user credentials
 * If the user is not registered, we create a player account.
 * In any case, a JWT is returned, used by the player to authenticate
 */
func Authenticate(w http.ResponseWriter, r *http.Request) {
  var body []byte
  var err error
	if body, err = ioutil.ReadAll(io.LimitReader(r.Body, 1048576)); err != nil {
    panic(err)
  }
	if err = r.Body.Close(); err != nil {
    panic(err)
  }
  jsonData := security.Decrypt(r.Header.Get("Application-Key"), r.Header.Get("Application-Iv"), body)
  var data map[string]string
  if err = json.Unmarshal(jsonData, &data); err != nil {
    panic(err)
  }
  server := manager.GetServerBySignature(data["signature"])
  if server == nil {
    w.WriteHeader(http.StatusNotFound)
    return
  }
  player := manager.GetPlayerByUsername(data["username"], server)
  if player == nil {
    player = manager.CreatePlayer(data["username"], server)
  }
  if server.Id != player.Server.Id {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("Invalid server data"))
    return
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
    tokenString, error := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if error != nil {
        fmt.Println(error)
    }
    return tokenString
}
