package api

import(
    "log"
    "net/http"
)

func ErrorHandler(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        defer CatchHttpException(w)
        next(w, req)
    })
}

func CatchHttpException(w http.ResponseWriter) {
    r := recover()
    if r == nil {
        return
    }
    if exception, ok := r.(*HttpException); ok {
        message := ""
        if exception.Exception.Error != nil {
            message = "; [Error]: " + exception.Exception.Error.Error()
        }
        if exception.Exception.Message != "" || message != "" {
            log.Println("[Exception]: " + exception.Exception.Message + message)
        }
        SendJsonResponse(w, exception.Code, exception)
        return
    }
    if exception, ok := r.(*Exception); ok {
        message := ""
        if exception.Error != nil {
            message = "; [Error]: " + exception.Error.Error()
        }
        if exception.Message != "" || message != "" {
            log.Println("[Exception]: " + exception.Message + message)
        }
        SendJsonResponse(w, 500, exception)
        return
    }
    if err, ok := r.(error); ok {
        log.Println("[Error]: " + err.Error())
        SendJsonResponse(w, 500, "Internal server error")
        return
    }
    panic(r)
}

func CatchException() {
    r := recover()
    if r == nil {
        return
    }
    if exception, ok := r.(*Exception); ok {
        message := ""
        if exception.Error != nil {
            message = "; [Error]: " + exception.Error.Error()
        }
        if exception.Message != "" || message != "" {
            log.Println("[Exception]: " + exception.Message + message)
        }
        return
    }
    if err, ok := r.(error); ok {
        log.Println("[Error]: " + err.Error())
        return
    }
    panic(r)
}
