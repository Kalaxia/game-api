package api

import(
    "encoding/json"
    "io"
    "io/ioutil"
    "net/http"
)

func DecodeJsonRequest(r *http.Request) map[string]interface{} {
    var body []byte
    var err error
    if body, err = ioutil.ReadAll(io.LimitReader(r.Body, 1048576)); err != nil {
        panic(NewHttpException(500, "Could not read rrequest body", err))
    }
    if err = r.Body.Close(); err != nil {
        panic(NewHttpException(500, "Could not close request body", err))
    }
    if r.Header.Get("Application-Iv") != "" {
        body = Decrypt(r.Header.Get("Application-Key"), r.Header.Get("Application-Iv"), body)
    }
    var data map[string]interface{}
    if err = json.Unmarshal(body, &data); err != nil {
        panic(NewHttpException(500, "Could not decode request body", err))
    }
    return data
}

func SendJsonResponse(w http.ResponseWriter, code int, data interface{}) {
    w.WriteHeader(code)
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(&data); err != nil {
        panic(NewHttpException(500, "Could not encode response data", err))
    }
}
