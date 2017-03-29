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
	m.Navbar = bootstrap.NewNavBar()

	m.btnSearch = m.AddButton(bootstrap.ButtonDefault, "Search")
	m.btnIndex = m.AddButton(bootstrap.ButtonDefault, "Index")
	m.btnBackup = m.AddButton(bootstrap.ButtonDefault, "Backup")
	m.btnSettings = m.AddButton(bootstrap.ButtonDefault, "Settings")
	m.btnAbout = m.AddButton(bootstrap.ButtonDefault, "About")

	return m
}
