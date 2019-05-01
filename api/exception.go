package api

type (
    Exception struct {
        Message string `json:"message"`
        Error error `json:"-"`
    }
    HttpException struct {
        Exception
        Code int `json:"code"`
    }
)

func NewHttpException(code int, message string, err error) *HttpException {
    return &HttpException {
        Exception: Exception{
            Message: message,
            Error: err,
        },
        Code: code,
    }
}

func NewException(message string, err error) *Exception {
    return &Exception {
        Message: message,
        Error: err,
    }
}
