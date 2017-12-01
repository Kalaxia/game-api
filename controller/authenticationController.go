package controller

import (
  "fmt"
  "net/http"
  "encoding/json"
	"io"
	"io/ioutil"
  "kalaxia-game-api/security"
  "github.com/dgrijalva/jwt-go"
)

/**
 * This method receives crypted data from the portal, containing user credentials
 * If the user is not registered, we create a player account.
 * In any case, a JWT is returned, used by the player to authenticate
 */
func AuthenticatePlayer(w http.ResponseWriter, r *http.Request) {
  var body []byte
  var err error
	if body, err = ioutil.ReadAll(io.LimitReader(r.Body, 1048576)); err != nil {
    panic(err)
  }
	if err = r.Body.Close(); err != nil {
    panic(err)
  }
  jsonData := security.Decrypt(body)
  var data map[string]interface{}
  if err = json.Unmarshal(jsonData, &data); err != nil {
    panic(err)
  }
  token := getNewJWT(data)
  //json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
  w.Write(security.Encrypt([]byte(token)))
}

func getNewJWT(data map[string]interface{}) string {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": data["username"],
    })
    tokenString, error := token.SignedString([]byte("secret"))
    if error != nil {
        fmt.Println(error)
    }
    return tokenString
}
