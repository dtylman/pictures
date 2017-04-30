package view

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/pictures/cmd/app/view/darktheme"
)

type view interface {
	updateState()
	getContent() *gowd.Element
	populateToolbar(menu*darktheme.Menu)
}

