package tasklog

import (
	"fmt"
	"os"
)

func Println(v ...interface{}) {
	fmt.Fprintln(os.Stderr,v)
}

func Printf(format string, v ...interface{}) {
	Println(fmt.Sprintf(format, v...))
}

func Error(err error) {
	fmt.Fprintf(os.Stderr,err.Error())
}
