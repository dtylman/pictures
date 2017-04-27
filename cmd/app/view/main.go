package view

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	"github.com/dtylman/pictures/indexer"
	"github.com/dtylman/pictures/cmd/app/view/darktheme"
)

type main struct {
	*gowd.Element

	menu     *gowd.Element

	alerts   *darktheme.Alerts

	content  *gowd.Element

	//views
	search   view
	indexer  view
	indexing view
	backup   view
	settings view
}

func newMain() *main {
	m := new(main)

	// body
	m.Element = gowd.NewElement("div")
	m.Element.SetID("wrapper")

	// menu
	menu := darktheme.NewMenu()
	menu.AddSideButton("Search", "fa fa-search", m.btnSearchClick)
	//search, google style
	menu.AddSideButton("Browse", "fa fa-list", m.btnSearchClick)
	//albums, locations, timeline (?)
	menu.AddSideButton("Thumbs", "fa fa-image", m.btnSearchClick)
	//show search results in thumbs
	menu.AddSideButton("Actions", "fa fa-cog", m.btnSearchClick)
	//show table with search results, something you can work on
	menu.AddSideButton("Faces", "fa fa-user", m.btnSearchClick)
	//show and manage faces
	menu.AddSideButton("Index", "fa fa-database", m.btnIndexClick)
	//do the indexing
	menu.AddSideButton("Manage", "fa fa-adjust", m.btnBackupClick)
	//backup, restore, remove_from_source
	menu.AddSideButton("Settings", "fa fa-gears", m.btnSettingsClick)
	//about
	menu.AddTopButton("About", "fa fa-question", m.btnSettingsClick)

	m.menu = menu.Element

	btn := menu.AddTopButton("Close", "fa fa-close", nil)
	btn.SetAttribute("onclick", "window.close()")

	m.AddElement(menu.Element)

	m.alerts = darktheme.NewAlerts()
	m.AddElement(bootstrap.NewContainer(true, m.alerts.Element))

	//content
	m.content = bootstrap.NewContainer(true)
	m.AddElement(m.content)

	//views
	m.search = newSearchView()
	m.indexer = newIndexerView()
	m.indexing = newIndexingView()
	m.backup = newBackupView()
	m.settings = newSettingsView()
	// footer
	return m
}

func (m *main) btnSearchClick(*gowd.Element, *gowd.EventElement) {
	m.setActiveView(m.search)
}

func (m *main) btnIndexClick(sender *gowd.Element, e *gowd.EventElement) {
	if indexer.IsRunning() {
		m.setActiveView(m.indexing)
	} else {
		m.setActiveView(m.indexer)
	}
}

func (m *main) btnSettingsClick(*gowd.Element, *gowd.EventElement) {
	m.setActiveView(m.settings)
}

func (m *main) btnBackupClick(*gowd.Element, *gowd.EventElement) {
	m.setActiveView(m.backup)
}

func (m *main) setActiveView(view view) {
	view.updateState()
	m.content.RemoveElements()
	m.content.AddElement(view.getContent())
}

func (m *main) addAlert(title string, caption string, alertType string) {
	m.alerts.Add(caption)
}

func (m *main) addAlertError(err error) {
	m.addAlert("Error", err.Error(), bootstrap.AlertDanger)
}
