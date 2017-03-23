package indexview

import (
	"github.com/dtylman/pictures/webkit"
	"github.com/dtylman/pictures/webkit/bootstrap"
)

var (
	View          *webkit.Element
	ProgressBar   *bootstrap.ProgressBar
	BtnStop       *webkit.Element
	BtnStart      *webkit.Element
	BtnSave       *webkit.Element
	ChkLocation   *bootstrap.Checkbox
	ChkReindex    *bootstrap.Checkbox
	InputMapQuest *bootstrap.Input
	SourceFolders *webkit.Element
)

func init() {
	View = bootstrap.NewElement("form", "form-horizontal")
	ChkLocation = bootstrap.NewCheckBox("Include Locations", false)
	ChkReindex = bootstrap.NewCheckBox("Reindex Existing Items", false)
	InputMapQuest = bootstrap.NewInput(bootstrap.InputTypeText, "MapQuest API Key")
	InputMapQuest.SetHelpText("Required for Geolocation")
	InputMapQuest.SetPlaceHolder("API KEY...")
	SourceFolders = webkit.NewElement("div")

	ProgressBar = bootstrap.NewProgressBar()
	BtnStart = bootstrap.NewButton(bootstrap.ButtonPimary, "Start")
	BtnStop = bootstrap.NewButton(bootstrap.ButtonPimary, "Stop")
	BtnSave = bootstrap.NewButton(bootstrap.ButtonDefault, "Save")
	//BtnAddFolder = <input  type="file" id="fileDialog" nwdirectory="true" onchange="alert(this.value);" />

	View.AddElement(ChkLocation.Element)
	View.AddElement(ChkReindex.Element)
	View.AddElement(InputMapQuest.Element)
	View.AddElement(SourceFolders)
	View.AddElement(ProgressBar.Element)
	View.AddElement(BtnStart)
	View.AddElement(BtnStop)
	View.AddElement(BtnSave)
}

func OnIndexerStopped() {
	BtnStop.Hide()
	BtnStart.Show()
	ProgressBar.Hide()
}

func OnIndexerStarted() {
	BtnStop.Show()
	BtnStart.Hide()
	ProgressBar.Show()
}
