package handler

import(
  "fmt"
  "os"
  "net/http"
  "strings"

  "github.com/dgrijalva/jwt-go"
  "github.com/gorilla/context"
)

func JwtHandler(next http.HandlerFunc, isProtected bool) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        if (isProtected == false) {
          next(w, req)
          return
        }
        authorizationHeader := req.Header.Get("authorization")
        if authorizationHeader == "" {
          w.WriteHeader(http.StatusUnauthorized)
          w.Write([]byte("An authorization header is required"))
          return
        }
        bearerToken := strings.Split(authorizationHeader, " ")
        if len(bearerToken) != 2 {
          w.WriteHeader(http.StatusBadRequest)
          w.Write([]byte("The Authorization header format is invalid"))
          return
        }
        token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("There was an error")
            }
            return []byte(os.Getenv("JWT_SECRET")), nil
        })
        if err != nil {
            panic(err)
            return
        }
        if token.Valid {
            context.Set(req, "jwt", token.Claims)
            next(w, req)
        } else {
            w.WriteHeader(http.StatusInternalServerError)
            w.Write([]byte("Invalid authorization token"))
        }
    })
}
