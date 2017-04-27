package main

import (
	"fmt"

	"github.com/dtylman/pictures/cmd/app/darktheme"
)

func worky() error {
	m := darktheme.NewMenu()
	err := m.Render()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := worky()
	if err != nil {
		fmt.Println(err)
	}

}
