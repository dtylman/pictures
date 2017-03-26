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
	btnAddFolder  *bootstrap.FileButton
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
	iv.Root = bootstrap.NewElement("div", "form-horizontal")

	iv.chkLocation = bootstrap.NewCheckBox("Include Locations", false)
	iv.chkReIndex = bootstrap.NewCheckBox("Reindex Existing Items", false)

	iv.inputMapQuest = bootstrap.NewInput(bootstrap.InputTypeText, "MapQuest API Key")
	iv.inputMapQuest.SetHelpText("Required for Geolocation")
	iv.inputMapQuest.SetPlaceHolder("API KEY...")
	iv.inputMapQuest.OnEvent(webkit.OnChange, iv.inputMapChanged)

	iv.SourceFolders = webkit.NewElement("div")

	iv.progressBar = bootstrap.NewProgressBar()
	iv.btnStart = bootstrap.NewButton(bootstrap.ButtonPrimary, "Start")
	iv.btnStop = bootstrap.NewButton(bootstrap.ButtonPrimary, "Stop")
	iv.btnAddFolder = bootstrap.NewFileButton(bootstrap.ButtonDefault, "Add folder", true)
	iv.btnAddFolder.OnChange(iv.btnAddFolderChanged)
	iv.Root.AddElement(iv.chkLocation.Element)
	iv.Root.AddElement(iv.chkReIndex.Element)
	iv.Root.AddElement(iv.inputMapQuest.Element)

	iv.Root.AddElement(iv.SourceFolders)
	iv.Root.AddElement(iv.btnAddFolder.Element)

	iv.Root.AddElement(iv.btnStart)
	iv.Root.AddElement(iv.btnStop)

	iv.Root.AddElement(iv.progressBar.Element)

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
	path := sender.Object.(string)
	conf.RemoveSourceFolder(path)
	err := conf.Save()
	if err != nil {
		MainView.addAlertError(err)
		return
	}
	OnConfigChanged()
}

func (iv *indexView) btnAddFolderChanged(sender *webkit.Element, event *webkit.EventElement) {
	path := iv.btnAddFolder.GetValue()
	conf.AddSourceFolder(path)
	OnConfigChanged()
}

func (iv *indexView) inputMapChanged(sender *webkit.Element, event *webkit.EventElement) {
	///conf.Options.MapQuestAPIKey = iv.inputMapQuest.GetValue()
	MainView.addAlert("News:", "Want to save this: "+iv.inputMapQuest.GetValue(), bootstrap.AlertInfo)

}
