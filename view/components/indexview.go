package components

import (
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/webkit"
	"github.com/dtylman/pictures/webkit/bootstrap"
)

type index struct {
	*webkit.Element
	progressBar   *bootstrap.ProgressBar
	btnStop       *webkit.Element
	btnStart      *webkit.Element
	btnAddFolder  *bootstrap.FileButton
	chkLocation   *bootstrap.Checkbox
	chkReIndex    *bootstrap.Checkbox
	inputMapQuest *bootstrap.Input
	SourceFolders *webkit.Element
}

func newIndex() *index {
	iv := new(index)
	iv.Element = bootstrap.NewElement("div", "form-horizontal")

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
	iv.AddElement(iv.chkLocation.Element)
	iv.AddElement(iv.chkReIndex.Element)
	iv.AddElement(iv.inputMapQuest.Element)

	iv.AddElement(iv.SourceFolders)
	iv.AddElement(iv.btnAddFolder.Element)

	iv.AddElement(iv.btnStart)
	iv.AddElement(iv.btnStop)

	iv.AddElement(iv.progressBar.Element)

	//set initial state
	//if indexer.isrunning,
	//indexer.setProgressHandler(me)
	//indexer.start()
	//indexer.stop()....
	iv.onIndexerStopped()
	iv.onConfigChanged()
	return iv
}

//OnIndexerStopped happens then indexer had started
func (iv *index) onIndexerStopped() {
	iv.btnStop.Hide()
	iv.btnStart.Show()
	iv.progressBar.Hide()
}

//OnIndexerStarted happens when indexer had started
func (iv *index) onIndexerStarted() {
	iv.btnStop.Show()
	iv.btnStart.Hide()
	iv.progressBar.Show()
}

func (iv *index) onConfigChanged() {
	iv.inputMapQuest.SetValue(conf.Options.MapQuestAPIKey)
	iv.SourceFolders.RemoveElements()
	for _, path := range conf.Options.SourceFolders {
		iv.SourceFolders.AddElement(NewSourceFolder(path, iv.btnSourceFolderDelete))
	}
}

func (iv *index) btnSourceFolderDelete(sender *webkit.Element, event *webkit.EventElement) {
	path := sender.Object.(string)
	conf.RemoveSourceFolder(path)
	iv.saveConfig(false)
}

func (iv *index) btnAddFolderChanged(sender *webkit.Element, event *webkit.EventElement) {
	path := iv.btnAddFolder.GetValue()
	conf.AddSourceFolder(path)
	iv.saveConfig(false)
}

func (iv *index) inputMapChanged(sender *webkit.Element, event *webkit.EventElement) {
	conf.Options.MapQuestAPIKey = iv.inputMapQuest.GetValue()
	iv.saveConfig(true)
}

func (iv *index) saveConfig(showAlert bool) {
	if showAlert {
		//		Root.addAlert("", "Configuration Saved.", bootstrap.AlertInfo)
	}
	err := conf.Save()
	if err != nil {
		//		Root.addAlertError(err)
		return
	}
	iv.onConfigChanged()
}
