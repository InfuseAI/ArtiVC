package log

import (
	"log"
	"os"
)

var logger *log.Logger

func SetDebug(debug bool) {
	if debug {
		logger = log.New(os.Stderr, "[DBG] ", log.Ldate|log.Lmicroseconds)
	} else {
		logger = nil
	}
}

func Debug(v ...interface{}) {
	if logger == nil {
		return
	}

	logger.Print(v...)
}

func Debugf(format string, v ...interface{}) {
	if logger == nil {
		return
	}

	logger.Printf(format, v...)
}

func Debugln(v ...interface{}) {
	if logger == nil {
		return
	}

	logger.Println(v...)
}
