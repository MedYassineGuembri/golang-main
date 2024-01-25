package logger

import (
	"log"
	"os"
)

var Logger *log.Logger

func init() {
    file, err := os.OpenFile("requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }

    Logger = log.New(file, "", log.LstdFlags)
}