package view

import "github.com/dtylman/gowd"

type view interface {
	updateState()
	getContent() *gowd.Element
}

var Root *main

func InitializeComponents() {
	Root = newMain()
}
