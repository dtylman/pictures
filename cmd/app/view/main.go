package view

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	"github.com/dtylman/pictures/cmd/app/view/darktheme"
	"github.com/dtylman/pictures/indexer"
	"github.com/dtylman/pictures/model"
)

var (
	Root         *main
	activeSearch *model.Search
)

func InitializeComponents() {
	Root = newMain()
}

type main struct {
	*gowd.Element

	menu *darktheme.Menu

	alerts *darktheme.Alerts

	content *gowd.Element

	//views
	search   view
	thumb    view
	table    view
	indexer  view
	indexing view
	backup   view
	settings view
	about    view
}

func newMain() *main {
	m := new(main)

	// body
	m.Element = gowd.NewElement("div")
	m.Element.SetID("wrapper")

	// menu
	m.menu = darktheme.NewMenu()
	m.AddElement(m.menu.Element)
	m.popluateMenu()

	m.alerts = darktheme.NewAlerts()
	m.AddElement(bootstrap.NewContainer(true, m.alerts.Element))

	//content
	m.content = bootstrap.NewContainer(true)
	m.AddElement(m.content)

	//views
	m.search = newSearchView()
	m.thumb = newThumbView()
	m.table = newTableView()
	m.indexer = newIndexerView()
	m.indexing = newIndexingView()
	m.backup = newBackupView()
	m.settings = newSettingsView()
	m.about = newAboutView()

	return m
}

func (m *main) popluateMenu() {
	// build static buttons...
	m.menu.AddButton(m.menu.Side, "Search", "fa fa-search", m.btnSearchClick)
	//search, google style
	m.menu.AddButton(m.menu.Side, "Browse", "fa fa-list", m.btnSearchClick)
	//albums, locations, timeline (?)
	m.menu.AddButton(m.menu.Side, "Thumbs", "fa fa-image", m.btnThumbClick)
	//show search results in thumbs
	m.menu.AddButton(m.menu.Side, "Actions", "fa fa-cog", m.btnTableClick)
	//show table with search results, something you can work on
	m.menu.AddButton(m.menu.Side, "Faces", "fa fa-user", m.btnSearchClick)
	//show and manage faces
	m.menu.AddButton(m.menu.Side, "Index", "fa fa-database", m.btnIndexClick)
	//do the indexing
	m.menu.AddButton(m.menu.Side, "Manage", "fa fa-adjust", m.btnBackupClick)
	//backup, restore, remove_from_source
	m.menu.AddButton(m.menu.Side, "Settings", "fa fa-gears", m.btnSettingsClick)

	//about
	m.menu.AddButton(m.menu.Side, "About", "fa fa-question", m.btnAboutClick)

	btn := m.menu.AddButton(m.menu.TopRight, "Close", "fa fa-close", nil)
	btn.SetAttribute("onclick", "window.close()")
}

func (m *main) btnSearchClick(*gowd.Element, *gowd.EventElement) {
	m.setActiveView(m.search)
}

func (m *main) btnThumbClick(*gowd.Element, *gowd.EventElement) {
	m.setActiveView(m.thumb)
}

func (m *main) btnTableClick(*gowd.Element, *gowd.EventElement) {
	m.setActiveView(m.table)
}

func (m *main) btnIndexClick(sender *gowd.Element, e *gowd.EventElement) {
	if indexer.IsRunning() {
		m.setActiveView(m.indexing)
	} else {
		m.setActiveView(m.indexer)
	}
}

func (m *main) btnAboutClick(*gowd.Element, *gowd.EventElement) {
	m.setActiveView(m.about)
}

func (m *main) btnSettingsClick(*gowd.Element, *gowd.EventElement) {
	m.setActiveView(m.settings)
}

func (m *main) btnBackupClick(*gowd.Element, *gowd.EventElement) {
	m.setActiveView(m.backup)
}

func (m *main) setActiveView(view view) {
	view.updateState()

	m.menu.TopLeft.RemoveElements()

	if view != nil {
		view.populateToolbar(m.menu)
	}

	m.content.RemoveElements()
	m.content.AddElement(view.getContent())
}

func (m *main) addAlert(title string, caption string, alertType string) {
	m.alerts.Add(caption)
}

func (m *main) addAlertError(err error) {
	m.addAlert("Error", err.Error(), bootstrap.AlertDanger)
}
