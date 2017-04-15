package tasklog

import "log"

func Error(err error){
	log.Println(err)
}

func ErrorF(format string, v ...interface{}){
	log.Printf(format, v...)
}
