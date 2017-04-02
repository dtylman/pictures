package view

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/indexer"
)

type index struct {
	*gowd.Element
	progressBar   *bootstrap.ProgressBar
	btnStop       *gowd.Element
	btnStart      *gowd.Element
	btnAddFolder  *bootstrap.FileButton
	chkLocation   *bootstrap.Checkbox
	chkReIndex    *bootstrap.Checkbox
	inputMapQuest *bootstrap.FormInput
	SourceFolders *gowd.Element
}

func newIndex() *index {
	i := new(index)
	i.Element = gowd.NewElement("div")

	i.chkLocation = bootstrap.NewCheckBox("Include Locations", false)
	i.chkReIndex = bootstrap.NewCheckBox("Reindex Existing Items", false)

	i.SourceFolders = gowd.NewElement("div")

	i.inputMapQuest = bootstrap.NewFormInput(bootstrap.InputTypeText, "MapQuest API Key")
	i.inputMapQuest.SetHelpText("Required for Geolocation")
	i.inputMapQuest.SetPlaceHolder("API KEY...")
	i.inputMapQuest.OnEvent(gowd.OnChange, i.inputMapChanged)

	i.progressBar = bootstrap.NewProgressBar()

	i.btnStart = bootstrap.NewButton(bootstrap.ButtonPrimary, "Start")
	i.btnStart.OnEvent(gowd.OnClick, i.btnStartClicked)

	i.btnStop = bootstrap.NewButton(bootstrap.ButtonPrimary, "Stop")
	i.btnStop.OnEvent(gowd.OnClick, i.btnStopClicked)

	i.btnAddFolder = bootstrap.NewFileButton(bootstrap.ButtonDefault, "Add folder", true)
	i.btnAddFolder.OnChange(i.btnAddFolderChanged)

	i.AddElement(i.chkLocation.Element)
	i.AddElement(i.chkReIndex.Element)
	i.AddElement(i.inputMapQuest.Element)

	i.AddElement(i.SourceFolders)
	i.AddElement(i.btnAddFolder.Element)

	i.AddElement(i.btnStart)
	i.AddElement(i.btnStop)

	i.AddElement(i.progressBar.Element)

	i.updateState()

	return i
}

func (i *index) updateState() {
	i.SourceFolders.RemoveElements()
	for _, path := range conf.Options.SourceFolders {
		i.SourceFolders.AddElement(NewSourceFolder(path, i.btnSourceFolderDelete))
	}
	i.inputMapQuest.SetValue(conf.Options.MapQuestAPIKey)
	if indexer.IsRunning() {
		i.btnStart.Disable()
		i.btnStop.Enable()
	} else {
		i.btnStart.Enable()
		i.btnStop.Disable()
	}

}

func (i *index) updateIndexerProgress(progress indexer.IndexerProgress) {
	i.progressBar.SetText(progress.Text())
	i.progressBar.SetValue(progress.Percentage())
	if !progress.Running {
		i.updateState()
	}
	Root.Render()
}

func (i *index) btnSourceFolderDelete(sender *gowd.Element, event *gowd.EventElement) {
	path := sender.Object.(string)
	conf.RemoveSourceFolder(path)
	i.saveConfig(false)
}

func (i *index) btnStartClicked(sender *gowd.Element, event *gowd.EventElement) {
	err := indexer.Start(indexer.Options{
		IndexLocation:   i.chkLocation.Checked(),
		ReIndex:         i.chkReIndex.Checked(),
		ProgressHandler: i.updateIndexerProgress,
	})
	if err != nil {
		Root.addAlertError(err)
	}
	i.updateState()
}

func (i *index) btnStopClicked(sender *gowd.Element, event *gowd.EventElement) {
	indexer.Stop()
	i.updateState()
}

func (i *index) btnAddFolderChanged(sender *gowd.Element, event *gowd.EventElement) {
	path := i.btnAddFolder.GetValue()
	conf.AddSourceFolder(path)
	i.saveConfig(false)
}

func (i *index) inputMapChanged(sender *gowd.Element, event *gowd.EventElement) {
	conf.Options.MapQuestAPIKey = i.inputMapQuest.GetValue()
	i.saveConfig(true)
}

func (i *index) saveConfig(showAlert bool) {
	if showAlert {
		Root.addAlert("", "Configuration Saved.", bootstrap.AlertInfo)
	}
	err := conf.Save()
	if err != nil {
		Root.addAlertError(err)
		return
	}
	i.updateState()
}
