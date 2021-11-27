package logs

import (
	logs "log"
	"os"
)

// log = logs.New(logFile, "", logs.Ldate|logs.Ltime)
var log = logs.New(os.Stdout, "", logs.Ldate|logs.Ltime)

func Println(v ...interface{}) {
	logs.Println(v...)
}

func Printf(format string, v ...interface{}) {
	logs.Printf(format, v...)
}