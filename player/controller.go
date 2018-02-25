package player

import (
    "io"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "os"
    "kalaxia-game-api/server"
    "strconv"
    "github.com/dgrijalva/jwt-go"
    "time"
)

/**
 * This method receives crypted data from the portal, containing user credentials
 * If the user is not registered, we create a player account.
 * In any case, a JWT is returned, used by the player to authenticate
 */
func AuthenticateAction(w http.ResponseWriter, r *http.Request) {
    var body []byte
    var err error
    if body, err = ioutil.ReadAll(io.LimitReader(r.Body, 1048576)); err != nil {
        panic(err)
    }
    if err = r.Body.Close(); err != nil {
        panic(err)
    }
    jsonData := server.Decrypt(r.Header.Get("Application-Key"), r.Header.Get("Application-Iv"), body)
    var data map[string]string
    if err = json.Unmarshal(jsonData, &data); err != nil {
        panic(err)
    }
    server := server.GetServerBySignature(data["signature"])
    if server == nil {
        w.WriteHeader(http.StatusNotFound)
        return
    }
    player := GetPlayerByUsername(data["username"], server)
    if player == nil {
        player = CreatePlayer(data["username"], server)
    }
    if server.Id != player.Server.Id {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Invalid server data"))
        return
    }
    token := getNewJWT(player)
    cipherData, key, iv := server.Encrypt([]byte(token))
    w.Header().Set("Application-Key", key)
    w.Header().Set("Application-Iv", iv)
    w.Write(cipherData)
}

func getNewJWT(player *Player) string {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "id": player.Id,
        "pseudo": player.Pseudo,
        "server_id": player.Server.Id,
        "created_at": time.Now().Format(time.RFC3339),
    })
    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        panic(err)
    }
    return tokenString
}

func GetCurrentPlayerAction(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player")
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(player); err != nil {
        panic(err)
    }
}

func RegisterPlayerAction(w http.ResponseWriter, r *http.Request) {
    var body []byte
    var err error
    if body, err = ioutil.ReadAll(io.LimitReader(r.Body, 1048576)); err != nil {
        panic(err)
    }
    if err = r.Body.Close(); err != nil {
        panic(err)
    }
    var data map[string]string
    if err = json.Unmarshal(body, &data); err != nil {
        panic(err)
    }
    player := context.Get(r, "player").(*Player)
    if player.IsActive == true {
        w.WriteHeader(http.StatusForbidden)
        return
    }
    factionId, _ := strconv.ParseUint(data["faction_id"], 10, 16)
    planetId, _ := strconv.ParseUint(data["planet_id"], 10, 16)
    RegisterPlayer(player, uint16(factionId), uint16(planetId))
    w.Write([]byte(""))
}
