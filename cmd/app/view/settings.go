package view

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	"github.com/dtylman/pictures/conf"
)

type settingsView struct {
	*gowd.Element
	btnSave         *gowd.Element
	inputDataFolder *bootstrap.FormInput
}

func newSettingsView() *settingsView {
	s := new(settingsView)
	s.Element = gowd.NewElement("div")

	s.btnSave = bootstrap.NewButton(bootstrap.ButtonPrimary, "Save")
	s.btnSave.OnEvent(gowd.OnClick, s.btnSaveClick)

	s.inputDataFolder = bootstrap.NewFormInput(bootstrap.InputTypeText, "Data folder")

	pnl := bootstrap.NewPanel(bootstrap.PanelDefault)
	pnl.AddTitle("Settings")
	pnl.AddToBody(s.inputDataFolder.Element)

	s.AddElement(pnl.Element)
	return s
}

func (sv *settingsView) updateState() {
	sv.inputDataFolder.SetValue(conf.Options.DataFolder)
}

func (sv *settingsView) populateToolbar(toolbar *gowd.Element) {
	toolbar.AddElement(sv.btnSave)
}

func (sv *settingsView) getContent() *gowd.Element {
	return sv.Element
}

func (sv *settingsView) btnSaveClick(sender *gowd.Element, event *gowd.EventElement) {
	conf.Options.DataFolder = sv.inputDataFolder.GetValue()
	conf.Save()
}