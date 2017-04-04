package view

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
)

type main struct {
	*gowd.Element
	menu    *mainMenu
	toolbar *gowd.Element
	alerts  *gowd.Element
	content *gowd.Element

	//views
	search view
	index  view
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

	m.toolbar = bootstrap.NewElement("div", "navbar btn-toolbar")
	m.toolbar.SetAttribute("style", "margin-top: 5px;")
	m.AddElement(m.toolbar)

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
	m.setActiveView(m.search)
}

func (m *main) btnIndexClick(sender *gowd.Element, e *gowd.EventElement) {
	m.setActiveView(m.index)
}

func (m *main) setActiveView(view view) {
	view.updateState()

	m.toolbar.RemoveElements()
	view.populateToolbar(m.toolbar)

	m.content.RemoveElements()
	m.content.AddElement(view.getContent())
}

func (m *main) addAlert(title string, caption string, alertType string) {
	m.alerts.AddElement(bootstrap.NewAlert(title, caption, alertType, true))
}

func (m *main) addAlertError(err error) {
	m.addAlert("Error", err.Error(), bootstrap.AlertDanger)
}
