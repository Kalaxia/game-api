package handler

import(
  "net/http"
  "strings"
  "kalaxia-game-api/exception"
  "kalaxia-game-api/security"
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
        token := security.AuthenticateJwt(bearerToken[1])
        context.Set(req, "jwt", token.Claims)
        next(w, req)
    })
}
