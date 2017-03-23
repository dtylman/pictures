package main

import (
	"fmt"
	"github.com/dtylman/pictures/presentation"
	"github.com/dtylman/pictures/webkit"
)

func main() {
	err := webkit.Run(presentation.MainView)
	if err != nil {
		fmt.Println(err)
	}
}
