package view

import (
	"fmt"
	"github.com/dtylman/pictures/webkit"
	"github.com/dtylman/pictures/webkit/bootstrap"
)

type mainView struct {
	Root    *webkit.Element
	content *webkit.Element
}

var MainView mainView

func init() {
	MainView.init()
}

func (mv *mainView) init() {
	// body
	mv.Root = bootstrap.NewContainer(true)

	// header
	// navbar
	navBar := bootstrap.NewNavBar()
	navBar.AddButton(bootstrap.ButtonDefault, "Search").OnEvent(webkit.OnClick, mv.btnSearchClick)
	navBar.AddButton(bootstrap.ButtonDefault, "Index").OnEvent(webkit.OnClick, mv.btnIndexClick)
	navBar.AddButton(bootstrap.ButtonDefault, "Backup")
	navBar.AddButton(bootstrap.ButtonDefault, "Settings")
	navBar.AddButton(bootstrap.ButtonDefault, "About")
	mv.Root.AddElement(navBar.Element)

	//content
	mv.content = bootstrap.NewContainer(true)
	mv.Root.AddElement(mv.content)

	// footer
}

func (mv *mainView) btnSearchClick(*webkit.Element, *webkit.EventElement) {
	mv.content.RemoveElements()
	mv.content.AddElement(webkit.NewText("Search"))
	txt := webkit.NewText("haha")
	mv.content.AddElement(txt)
	go func() {
		for i := 0; i < 1000; i++ {
			txt.SetText(fmt.Sprintf("haha %d", i))
			mv.Root.Render()
		}
	}()
}

func (mv *mainView) btnIndexClick(*webkit.Element, *webkit.EventElement) {
	mv.setActiveView(IndexView.Root)
}

func (mv *mainView) setActiveView(view *webkit.Element) {
	mv.content.RemoveElements()
	mv.content.AddElement(view)
}
