package view

import (
	"fmt"
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
)

//parentControls is an interface for the main page
type parentControls interface {
	addAlert(title string, caption string, alertType string)
	addAlertError(err error)
	Render() error
}

type main struct {
	*gowd.Element
	menu    *mainMenu
	alerts  *gowd.Element
	content *gowd.Element

	//views
	search *search
	index  *index
}

var root = newMain()

func newMain() *main {
	m := new(main)

	// body
	m.Element = bootstrap.NewContainer(true)
	btn := bootstrap.NewButton(bootstrap.ButtonPrimary, "lala")
	btn.Object = "lala"
	btn.OnEvent(gowd.OnClick, m.koko)
	m.AddElement(btn)
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
	m.index = newIndex(m)
	m.search = newSearch(m)
	// footer
	return m
}

func (m *main) koko(sender *gowd.Element, elemnt *gowd.EventElement) {

	for i := 0; i < 100; i++ {
		m.menu.btnIndex.SetText(fmt.Sprintf("counting %v", i))
		m.Render()
	}
}

func (m *main) btnSearchClick(*gowd.Element, *gowd.EventElement) {
	m.menu.SetActiveElement(m.menu.btnSearch)
	m.setActiveView(m.search.Element)
}

func (m *main) btnIndexClick(sender *gowd.Element, e *gowd.EventElement) {
	m.index.updateState()
	m.menu.SetActiveElement(m.menu.btnIndex)
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

//RootElement returns the root "body" container
func RootElement() *gowd.Element {
	return root.Element
}
