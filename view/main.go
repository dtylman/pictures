package view

import (
	"github.com/dtylman/pictures/webkit"
	"github.com/dtylman/pictures/webkit/bootstrap"
)

//parentControls is an interface for the main page
type parentControls interface {
	addAlert(title string, caption string, alertType string)
	addAlertError(err error)
	Render() error
}

type main struct {
	*webkit.Element
	alerts  *webkit.Element
	content *webkit.Element

	//views
	search *search
	index  *index
}

var root = newMain()

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

	// alerts
	m.alerts = webkit.NewElement("div")
	m.AddElement(m.alerts)

	//content
	m.content = bootstrap.NewContainer(true)
	m.AddElement(m.content)

	//views
	m.index = newIndex(m)
	m.search = newSearch(m)
	// footer
	return m
}

func (m *main) btnSearchClick(*webkit.Element, *webkit.EventElement) {
	m.setActiveView(m.search.Element)
}

func (m *main) btnIndexClick(*webkit.Element, *webkit.EventElement) {
	m.index.updateState()
	m.setActiveView(m.index.Element)
}

func (m *main) setActiveView(view *webkit.Element) {
	//view.onAlert(me)
	m.content.RemoveElements()
	m.content.AddElement(view)
}

func (m *main) addAlert(title string, caption string, alertType string) {
	m.alerts.AddElement(bootstrap.NewAlert(title, caption, alertType, true))
}

func (m *main) addAlertError(err error) {
	m.addAlert("Error", err.Error(), bootstrap.AlertDanger)
}

//RootElement returns the root "body" container
func RootElement() *webkit.Element {
	return root.Element
}
