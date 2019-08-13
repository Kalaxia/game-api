package api

import(
    "log"
    "net/http"
    "runtime"
    "strconv"
)

func ErrorHandler(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        defer CatchHttpException(w)
        next(w, req)
    })
}

func CatchHttpException(w http.ResponseWriter) {
    if isError, code, message := CatchException(recover()); isError {
        SendJsonResponse(w, code, &HttpException{
            Exception: Exception{ Message: message },
            Code: code,
        })
    }
}

func CatchException(r interface{}) (bool, int, string) {
    if r == nil {
        r = recover()
        if r == nil {
            return false, 200, ""
        }
    }
    finalMessage := "Internal server error"
    code := 500

    if exception, ok := r.(*HttpException); ok {
        message := ""
        if exception.Exception.Error != nil {
            message = "; [Error]: " + exception.Exception.Error.Error()
        }
        if exception.Exception.Message != "" || message != "" {
            log.Println("[Exception]: ", exception.Exception.Message + message)
        }
        code = exception.Code
        finalMessage = exception.Exception.Message
    } else if exception, ok := r.(*Exception); ok {
        message := ""
        if exception.Error != nil {
            message = "; [Error]: " + exception.Error.Error()
        }
        if exception.Message != "" || message != "" {
            log.Println("[Exception]: ", exception.Message + message)
        }
    } else if err, ok := r.(error); ok {
        log.Println("[Error]: " + err.Error())
    } else {
        log.Println("[Unknown Error]: " , r)
    }
    pc := make([]uintptr, 10)
    if n := runtime.Callers(0, pc); code >= 500 && n > 0 {
        logFrames(pc, n)
    }
    return true, code, finalMessage
}

func logFrames(pc []uintptr, n int) {
    pc = pc[:n]
    frames := runtime.CallersFrames(pc)

    for {
        frame, more := frames.Next()
        log.Println(frame.Function + " - " + frame.File + ":l" + strconv.Itoa(frame.Line))
        if !more {
            break
        }
    }
}