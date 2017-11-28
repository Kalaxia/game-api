package controller

import (
  "fmt"
  "net/http"
  //"encoding/json"
	"io"
	"io/ioutil"
  //"kalaxia-game-api/security"
)

func AuthenticatePlayer(w http.ResponseWriter, r *http.Request) {
  var body []byte
  var err error
	if body, err = ioutil.ReadAll(io.LimitReader(r.Body, 1048576)); err != nil {
    panic(err)
  }
	if err = r.Body.Close(); err != nil {
    panic(err)
  }

  fmt.Println(body);

  //decryptedData := "test"

}
