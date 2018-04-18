package handler

import(
  "fmt"
  "os"
  "net/http"
  "strings"
  "kalaxia-game-api/exception"
  "github.com/dgrijalva/jwt-go"
  "github.com/gorilla/context"
)

func JwtHandler(next http.HandlerFunc, isProtected bool) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        if isProtected == false {
          next(w, req)
          return
        }
        authorizationHeader := req.Header.Get("authorization")
        if authorizationHeader == "" {
          panic(exception.NewHttpException(http.StatusUnauthorized, "An authorization header is required", nil))
        }
        bearerToken := strings.Split(authorizationHeader, " ")
        if len(bearerToken) != 2 {
          panic(exception.NewHttpException(http.StatusBadRequest, "The Authorization header format is invalid", nil))
        }
        token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("there was an error")
            }
            return []byte(os.Getenv("JWT_SECRET")), nil
        })
        if err != nil {
            panic(exception.NewHttpException(500, "JWT could not be created", nil))
        }
        if token.Valid {
            context.Set(req, "jwt", token.Claims)
            next(w, req)
        } else {
            panic(exception.NewHttpException(http.StatusInternalServerError, "Invalid authorization token", nil))
        }
    })
}
