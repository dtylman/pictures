package view

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/pictures/cmd/app/view/darktheme"
)

type aboutView struct {
	*gowd.Element
}

func newAboutView() *aboutView {
	a := new(aboutView)
	a.Element = gowd.NewStyledText("What about that?", gowd.Heading1)
	return a
}

func (a *aboutView) updateState() {

}

func (a *aboutView) getContent() *gowd.Element {
	return a.Element
}

func (a *aboutView) populateToolbar(menu *darktheme.Menu) {

}
