package view

import (
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/indexer"
	"github.com/dtylman/pictures/webkit"
	"github.com/dtylman/pictures/webkit/bootstrap"
)

type index struct {
	*webkit.Element
	progressBar    *bootstrap.ProgressBar
	btnStop        *webkit.Element
	btnStart       *webkit.Element
	btnAddFolder   *bootstrap.FileButton
	chkLocation    *bootstrap.Checkbox
	chkReIndex     *bootstrap.Checkbox
	inputMapQuest  *bootstrap.Input
	SourceFolders  *webkit.Element
	parentControls parentControls
}

func newIndex(p parentControls) *index {
	i := &index{parentControls: p}
	i.Element = bootstrap.NewElement("div", "form-horizontal")
	i.chkLocation = bootstrap.NewCheckBox("Include Locations", false)
	i.chkReIndex = bootstrap.NewCheckBox("Reindex Existing Items", false)

	i.SourceFolders = webkit.NewElement("div")

	i.inputMapQuest = bootstrap.NewInput(bootstrap.InputTypeText, "MapQuest API Key")
	i.inputMapQuest.SetHelpText("Required for Geolocation")
	i.inputMapQuest.SetPlaceHolder("API KEY...")
	i.inputMapQuest.OnEvent(webkit.OnChange, i.inputMapChanged)

	i.progressBar = bootstrap.NewProgressBar()

	i.btnStart = bootstrap.NewButton(bootstrap.ButtonPrimary, "Start")
	i.btnStart.OnEvent(webkit.OnClick, i.btnStartClicked)

	i.btnStop = bootstrap.NewButton(bootstrap.ButtonPrimary, "Stop")
	i.btnStop.OnEvent(webkit.OnClick, i.btnStopClicked)

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
	i.parentControls.Render()
}

func (i *index) btnSourceFolderDelete(sender *webkit.Element, event *webkit.EventElement) {
	path := sender.Object.(string)
	conf.RemoveSourceFolder(path)
	i.saveConfig(false)
}

func (i *index) btnStartClicked(sender *webkit.Element, event *webkit.EventElement) {
	err := indexer.Start(indexer.Options{
		IndexLocation:   i.chkLocation.Checked(),
		ReIndex:         i.chkReIndex.Checked(),
		ProgressHandler: i.updateIndexerProgress,
	})
	if err != nil {
		i.parentControls.addAlertError(err)
	}
	i.updateState()
}

func (i *index) btnStopClicked(sender *webkit.Element, event *webkit.EventElement) {
	indexer.Stop()
	i.updateState()
}

func (i *index) btnAddFolderChanged(sender *webkit.Element, event *webkit.EventElement) {
	path := i.btnAddFolder.GetValue()
	conf.AddSourceFolder(path)
	i.saveConfig(false)
}

func (i *index) inputMapChanged(sender *webkit.Element, event *webkit.EventElement) {
	conf.Options.MapQuestAPIKey = i.inputMapQuest.GetValue()
	i.saveConfig(true)
}

func (i *index) saveConfig(showAlert bool) {
	if showAlert {
		i.parentControls.addAlert("", "Configuration Saved.", bootstrap.AlertInfo)
	}
	err := conf.Save()
	if err != nil {
		i.parentControls.addAlertError(err)
		return
	}
	i.updateState()
}
