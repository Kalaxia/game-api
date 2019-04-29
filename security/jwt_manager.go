package security

import(
	"fmt"
	"kalaxia-game-api/exception"
	"net/http"
	"os"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func AuthenticateJwt(rawToken string) *jwt.Token {
	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		panic(exception.NewException("JWT could not be created", nil))
	}
	if !token.Valid {
		panic(exception.NewException("Invalid authorization token", nil))
	}
	return token
}

func GetJwtPlayerId(data jwt.MapClaims) uint16 {
	createdAt, _ := time.Parse(time.RFC3339, data["created_at"].(string))
	if createdAt.Add(time.Hour * time.Duration(2)).Before(time.Now()) {
		panic(exception.NewHttpException(http.StatusUnauthorized, "Expired JWT", nil))
	}
	return uint16(data["id"].(float64))
}