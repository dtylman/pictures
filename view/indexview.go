package view

import (
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/webkit"
	"github.com/dtylman/pictures/webkit/bootstrap"
)

type indexView struct {
	Root          *webkit.Element
	progressBar   *bootstrap.ProgressBar
	btnStop       *webkit.Element
	btnStart      *webkit.Element
	btnSave       *webkit.Element
	chkLocation   *bootstrap.Checkbox
	chkReIndex    *bootstrap.Checkbox
	inputMapQuest *bootstrap.Input
	SourceFolders *webkit.Element
}

var IndexView indexView

func init() {
	IndexView.init()
}

func (iv *indexView) init() {
	iv.Root = bootstrap.NewElement("form", "form-horizontal")
	iv.chkLocation = bootstrap.NewCheckBox("Include Locations", false)
	iv.chkReIndex = bootstrap.NewCheckBox("Reindex Existing Items", false)
	iv.inputMapQuest = bootstrap.NewInput(bootstrap.InputTypeText, "MapQuest API Key")
	iv.inputMapQuest.SetHelpText("Required for Geolocation")
	iv.inputMapQuest.SetPlaceHolder("API KEY...")
	iv.SourceFolders = webkit.NewElement("div")

	iv.progressBar = bootstrap.NewProgressBar()
	iv.btnStart = bootstrap.NewButton(bootstrap.ButtonPimary, "Start")
	iv.btnStart.OnEvent(webkit.OnClick, iv.btnSourceFolderDelete)
	//iv.btnStop = bootstrap.NewButton(bootstrap.ButtonPimary, "Stop")
	//iv.btnSave = bootstrap.NewButton(bootstrap.ButtonDefault, "Save")
	//BtnAddFolder = <input  type="file" id="fileDialog" nwdirectory="true" onchange="alert(this.value);" />

	iv.Root.AddElement(iv.chkLocation.Element)
	iv.Root.AddElement(iv.chkReIndex.Element)
	iv.Root.AddElement(iv.inputMapQuest.Element)
	iv.Root.AddElement(iv.SourceFolders)
	iv.Root.AddElement(iv.progressBar.Element)
	iv.Root.AddElement(iv.btnStart)
	//iv.Root.AddElement(iv.btnStop)
	//iv.Root.AddElement(iv.btnSave)
}

//OnIndexerStopped happens then indexer had started
func (iv *indexView) onIndexerStopped() {
	iv.btnStop.Hide()
	iv.btnStart.Show()
	iv.progressBar.Hide()
}

//OnIndexerStarted happens when indexer had started
func (iv *indexView) onIndexerStarted() {
	iv.btnStop.Show()
	iv.btnStart.Hide()
	iv.progressBar.Show()
}

func (iv *indexView) onConfigChanged() {
	iv.inputMapQuest.SetValue(conf.Options.MapQuestAPIKey)
	iv.SourceFolders.RemoveElements()
	for _, path := range conf.Options.SourceFolders {
		iv.SourceFolders.AddElement(NewSourceFolder(path, iv.btnSourceFolderDelete))
	}
}

func (iv *indexView) btnSourceFolderDelete(sender *webkit.Element, event *webkit.EventElement) {
}
