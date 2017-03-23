package presentation

import (
	"fmt"
	"github.com/dtylman/pictures/presentation/indexview"
	"github.com/dtylman/pictures/webkit"
	"github.com/dtylman/pictures/webkit/bootstrap"
)

var (
	MainView *webkit.Element
	Content  *webkit.Element
)

func init() {
	// body
	MainView = bootstrap.NewContainer(true)

	// header
	// navbar
	navBar := bootstrap.NewNavBar()
	navBar.AddButton(bootstrap.ButtonDefault, "Search").OnEvent(webkit.OnClick, btnSearchClick)
	navBar.AddButton(bootstrap.ButtonDefault, "Index").OnEvent(webkit.OnClick, btnIndexClick)
	navBar.AddButton(bootstrap.ButtonDefault, "Backup")
	navBar.AddButton(bootstrap.ButtonDefault, "Settings")
	navBar.AddButton(bootstrap.ButtonDefault, "About")
	MainView.AddElement(navBar.Element)

	//content
	Content = bootstrap.NewContainer(true)
	MainView.AddElement(Content)

	// footer
}

func btnSearchClick(*webkit.Element, *webkit.EventElement) {
	Content.RemoveElements()
	Content.AddElement(webkit.NewText("Search"))
	txt := webkit.NewText("haha")
	Content.AddElement(txt)
	go func() {
		for i := 0; i < 1000; i++ {
			txt.SetText(fmt.Sprintf("haha %d", i))
			MainView.Render()
		}
	}()
}

func btnIndexClick(*webkit.Element, *webkit.EventElement) {
	showView(indexview.View)
}

func showView(view *webkit.Element) {
	Content.RemoveElements()
	Content.AddElement(view)
	MainView.Render()
}
