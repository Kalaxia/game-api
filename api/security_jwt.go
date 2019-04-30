package api

import(
	"fmt"
	"net/http"
	"os"
	"github.com/gorilla/context"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
)

func authenticateJwt(rawToken string) *jwt.Token {
	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		panic(NewException("JWT could not be created", nil))
	}
	if !token.Valid {
		panic(NewException("Invalid authorization token", nil))
	}
	return token
}

func getJwtPlayer(data jwt.MapClaims) *Player {
	createdAt, _ := time.Parse(time.RFC3339, data["created_at"].(string))
	if createdAt.Add(time.Hour * time.Duration(2)).Before(time.Now()) {
		panic(NewHttpException(http.StatusUnauthorized, "Expired JWT", nil))
	}
	if player := getPlayer(uint16(data["id"].(float64)), true); player != nil {
		return player
	}
	panic(NewHttpException(http.StatusInternalServerError, "Unavailable player account", nil))
}

func JwtHandler(next http.HandlerFunc, isProtected bool) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        if isProtected == false {
          next(w, req)
          return
        }
        authorizationHeader := req.Header.Get("authorization")
        if authorizationHeader == "" {
          panic(NewHttpException(http.StatusUnauthorized, "An authorization header is required", nil))
        }
        bearerToken := strings.Split(authorizationHeader, " ")
        if len(bearerToken) != 2 {
          panic(NewHttpException(http.StatusBadRequest, "The Authorization header format is invalid", nil))
        }
        token := authenticateJwt(bearerToken[1])
        context.Set(req, "jwt", token.Claims)
        next(w, req)
    })
}