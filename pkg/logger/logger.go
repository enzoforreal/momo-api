package logger

import (
	"io"
	"log"
	"os"
)

var (
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

func Init() {
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening log file: ", err)
	}

	multi := io.MultiWriter(os.Stdout, logFile)

	InfoLogger = log.New(multi, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(multi, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(multi, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
