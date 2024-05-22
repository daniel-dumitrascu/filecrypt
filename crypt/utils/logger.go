package utils

import (
	"log"
	"os"
	"sync"
)

type logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
}

func (lg *logger) Info(v ...any) {
	lg.infoLogger.Println(v...)
}

func (lg *logger) Error(v ...any) {
	lg.errorLogger.Println(v...)
}

func (lg *logger) Fatal(v ...any) {
	lg.fatalLogger.Println(v...)
	os.Exit(1)
}

var lock = &sync.Mutex{}
var loggerInstance Log

func GetLogger() Log {
	lock.Lock()
	defer lock.Unlock()
	if loggerInstance == nil {
		loggerInstance = &logger{
			infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
			errorLogger: log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime),
			fatalLogger: log.New(os.Stdout, "FATAL: ", log.Ldate|log.Ltime),
		}
	}

	return loggerInstance
}
