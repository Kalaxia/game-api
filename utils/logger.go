package utils

import(
    "os"
    "time"
)

func Log(message string) {
    date := time.Now()

    filename := "/var/log/api/" + date.Format("20060102150405") + ".log"

    f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        panic(err)
    }

    defer f.Close()

    if _, err = f.WriteString("[" + date.Format(time.UnixDate) + "] " + message + "\n"); err != nil {
        panic(err)
    }
}
