package view

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
)

type main struct {
	*gowd.Element
	menu    *mainMenu
	alerts  *gowd.Element
	content *gowd.Element

	//views
	search *search
	index  *index
}

func newMain() *main {
	m := new(main)

	// body
	m.Element = bootstrap.NewContainer(true)
	//menu
	m.menu = newMainMenu()
	m.menu.btnSearch.OnEvent(gowd.OnClick, m.btnSearchClick)
	m.menu.btnIndex.OnEvent(gowd.OnClick, m.btnIndexClick)
	m.AddElement(m.menu.Element)

	// alerts
	m.alerts = gowd.NewElement("div")
	m.AddElement(m.alerts)

	//content
	m.content = bootstrap.NewContainer(true)
	m.AddElement(m.content)

	//views
	m.index = newIndex()
	m.search = newSearch()
	// footer
	return m
}

func (m *main) btnSearchClick(*gowd.Element, *gowd.EventElement) {
	m.setActiveView(m.search.Element)
}

func (m *main) btnIndexClick(sender *gowd.Element, e *gowd.EventElement) {
	m.index.updateState()
	m.setActiveView(m.index.Element)
}

func (m *main) setActiveView(view *gowd.Element) {
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
