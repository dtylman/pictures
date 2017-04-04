package view

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
)

type mainMenu struct {
	*bootstrap.Navbar
	btnSearch   *gowd.Element
	btnIndex    *gowd.Element
	btnBackup   *gowd.Element
	btnSettings *gowd.Element
	btnAbout    *gowd.Element
}

func newMainMenu() *mainMenu {
	m := new(mainMenu)
	m.Navbar = bootstrap.NewNavBar(bootstrap.NavbarDefault)

	list := m.Navbar.AddList()
	list.SetClass("nav-pills")

	m.btnSearch = bootstrap.NewButton(bootstrap.ButtonDefault, "Search")
	m.btnSearch.SetClass("navbar-btn")
	list.AddItem(m.btnSearch)

	m.btnIndex = bootstrap.NewButton(bootstrap.ButtonDefault, "Index")
	m.btnIndex.SetClass("navbar-btn")
	list.AddItem(m.btnIndex)

	m.btnBackup = bootstrap.NewButton(bootstrap.ButtonDefault, "Backup")
	m.btnBackup.SetClass("navbar-btn")
	list.AddItem(m.btnBackup)

	m.btnSettings = bootstrap.NewButton(bootstrap.ButtonDefault, "Settings")
	m.btnSettings.SetClass("navbar-btn")
	list.AddItem(m.btnSettings)

	list = m.Navbar.AddList()
	list.SetClass("navbar-right")
	list.SetClass("nav-pills")

	m.btnAbout = bootstrap.NewButton(bootstrap.ButtonDefault, "About")
	m.btnAbout.SetClass("navbar-btn")
	list.AddItem(m.btnAbout)

	btnClose := bootstrap.NewButton(bootstrap.ButtonDefault, "Close")
	btnClose.SetAttribute("onclick", "window.close()")
	btnClose.SetClass("navbar-btn")
	list.AddItem(btnClose)

	return m
}
