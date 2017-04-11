package view

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/indexer"
	"github.com/dtylman/pictures/tasklog"
	"fmt"
)

//<a href="https://developer.mapquest.com/plan_purchase/steps/business_edition/business_edition_free/register">MapQuest </a>Key:

type index struct {
	*gowd.Element
	progressBar   *bootstrap.ProgressBar
	btnStop       *gowd.Element
	btnStart      *gowd.Element
	btnAddFolder  *bootstrap.FileButton
	chkLocation   *bootstrap.Checkbox
	chkDeleteDB   *bootstrap.Checkbox
	chkWithFaces  *bootstrap.Checkbox
	chkWithObjs   *bootstrap.Checkbox
	inputMapQuest *bootstrap.FormInput
	SourceFolders *gowd.Element
}

func newIndex() *index {
	i := new(index)
	i.Element = gowd.NewElement("div")

	i.chkLocation = bootstrap.NewCheckBox("With Locations", false)
	i.chkWithObjs = bootstrap.NewCheckBox("With Objects", true)
	i.chkWithFaces = bootstrap.NewCheckBox("With Faces", true)
	i.chkDeleteDB = bootstrap.NewCheckBox("Delete Existing Data", false)

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

	pnl := bootstrap.NewPanel(bootstrap.PanelDefault)
	pnl.AddTitle("Indexing Options:")
	pnl.AddToBody(i.chkLocation.Element)
	pnl.AddToBody(i.chkWithObjs.Element)
	pnl.AddToBody(i.chkWithFaces.Element)
	pnl.AddToBody(i.chkDeleteDB.Element)
	pnl.AddToBody(i.inputMapQuest.Element)
	i.AddElement(pnl.Element)

	pnl = bootstrap.NewPanel(bootstrap.PanelDefault)
	title := pnl.AddTitle("Source Folders: ")
	title.AddElement(gowd.NewStyledText("(These will be scanned)", gowd.SmallText))
	pnl.AddToBody(i.SourceFolders)
	pnl.AddToBody(i.btnAddFolder.Element)

	i.AddElement(pnl.Element)

	i.AddElement(gowd.NewStyledText("Progress:", gowd.Paragraph))
	i.AddElement(i.progressBar.Element)

	tasklog.RegisterHandler(tasklog.IndexerTask, i.updateIndexerProgress)

	i.updateState()

	return i
}

func (i *index) getContent() *gowd.Element {
	return i.Element
}

func (i *index) populateToolbar(toolbar *gowd.Element) {
	toolbar.AddElement(i.btnStart)
	toolbar.AddElement(i.btnStop)
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

func (i *index) updateIndexerProgress(status tasklog.Task) {
	i.progressBar.SetText(fmt.Sprintf("%v", status.Messages))
	i.progressBar.SetPercent(100)
	if !status.Running {
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
		WithLocation:   i.chkLocation.Checked(),
		DeleteDatabase:         i.chkDeleteDB.Checked(),
		WithFaces: i.chkWithFaces.Checked(),
		WithObjects: i.chkWithObjs.Checked(),
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
