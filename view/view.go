package view

import "github.com/dtylman/gowd"

type view interface {
	populateToolbar(toolbar *gowd.Element)
	updateState()
	getContent() *gowd.Element
}

var Root *main

func InitializeComponents() {
	Root = newMain()
}
