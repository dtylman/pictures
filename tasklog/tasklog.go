package tasklog

import (
	"fmt"
)

func Println(v ...interface{}) {
}

func Printf(format string, v ...interface{}) {
	Println(fmt.Sprintf(format, v...))
}

func Error(err error) {

}
