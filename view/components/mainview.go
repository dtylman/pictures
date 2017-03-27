package components

import (
	"github.com/dtylman/pictures/webkit"
	"github.com/dtylman/pictures/webkit/bootstrap"
)

type main struct {
	*webkit.Element
	alert   *webkit.Element
	content *webkit.Element
}

var Root = newMain()

func newMain() *main {
	m := new(main)

	// body
	m.Element = bootstrap.NewContainer(true)

	// header
	// navbar
	navBar := bootstrap.NewNavBar()
	navBar.AddButton(bootstrap.ButtonDefault, "Search").OnEvent(webkit.OnClick, m.btnSearchClick)
	navBar.AddButton(bootstrap.ButtonDefault, "Index").OnEvent(webkit.OnClick, m.btnIndexClick)
	navBar.AddButton(bootstrap.ButtonDefault, "Backup")
	navBar.AddButton(bootstrap.ButtonDefault, "Settings")
	navBar.AddButton(bootstrap.ButtonDefault, "About")
	m.AddElement(navBar.Element)
	// alert
	m.alert = webkit.NewElement("div")
	m.AddElement(m.alert)

	//content
	m.content = bootstrap.NewContainer(true)
	m.AddElement(m.content)

	// footer

	return m
}

func (m *main) btnSearchClick(*webkit.Element, *webkit.EventElement) {
	m.setActiveView(newSearch().Element)
}

func (m *main) btnIndexClick(*webkit.Element, *webkit.EventElement) {
	m.setActiveView(newIndex().Element)
}

func (m *main) setActiveView(view *webkit.Element) {
	//view.onAlert(me)
	m.content.RemoveElements()
	m.content.AddElement(view)
}

func (m *main) addAlert(title string, caption string, alertType string) {
	m.alert.AddElement(bootstrap.NewAlert(title, caption, alertType, true))
}

func (m *main) addAlertError(err error) {
	m.addAlert("Error", err.Error(), bootstrap.AlertDanger)
}
