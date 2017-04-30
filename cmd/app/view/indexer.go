package view

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/indexer"
	"github.com/dtylman/pictures/cmd/app/view/darktheme"
)

//<a href="https://developer.mapquest.com/plan_purchase/steps/business_edition/business_edition_free/register">MapQuest </a>Key:

type indexerView struct {
	*gowd.Element

	btnAddFolder  *bootstrap.FileButton
	chkLocation   *bootstrap.Checkbox
	chkDeleteDB   *bootstrap.Checkbox
	chkWithFaces  *bootstrap.Checkbox
	chkWithObjs   *bootstrap.Checkbox
	chkQuicksacn  *bootstrap.Checkbox
	inputMapQuest *bootstrap.FormInput
	SourceFolders *gowd.Element
}

func newIndexerView() *indexerView {
	i := new(indexerView)
	i.Element = gowd.NewElement("div")

	i.chkQuicksacn = bootstrap.NewCheckBox("Quick Scan", false)
	i.chkLocation = bootstrap.NewCheckBox("With Locations", false)
	i.chkWithObjs = bootstrap.NewCheckBox("With Objects", true)
	i.chkWithFaces = bootstrap.NewCheckBox("With Faces", true)
	i.chkDeleteDB = bootstrap.NewCheckBox("Delete Existing Data", false)

	i.SourceFolders = gowd.NewElement("div")

	i.inputMapQuest = bootstrap.NewFormInput(bootstrap.InputTypeText, "MapQuest API Key")
	i.inputMapQuest.SetHelpText("Required for Geolocation")
	i.inputMapQuest.SetPlaceHolder("API KEY...")
	i.inputMapQuest.OnEvent(gowd.OnChange, i.inputMapChanged)

	i.btnAddFolder = bootstrap.NewFileButton(bootstrap.ButtonDefault, "Add folder", true)
	i.btnAddFolder.OnChange(i.btnAddFolderChanged)

	pnl := bootstrap.NewPanel(bootstrap.PanelDefault)
	pnl.AddTitle("Indexing Options:")
	pnl.AddToBody(i.chkQuicksacn.Element)
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

	i.updateState()

	return i
}

func (i*indexerView) populateToolbar(menu*darktheme.Menu) {
	menu.AddTopButton("Start", "fa fa-start", i.btnStartClicked)
}

func (i *indexerView) getContent() *gowd.Element {
	return i.Element
}

func (i *indexerView) updateState() {
	if indexer.IsRunning() {
		Root.setActiveView(Root.indexing)
		return
	}
	i.SourceFolders.RemoveElements()
	for _, path := range conf.Options.SourceFolders {
		i.SourceFolders.AddElement(NewSourceFolder(path, i.btnSourceFolderDelete))
	}
	i.inputMapQuest.SetValue(conf.Options.MapQuestAPIKey)
}

func (i *indexerView) btnSourceFolderDelete(sender *gowd.Element, event *gowd.EventElement) {
	path := sender.Object.(string)
	conf.RemoveSourceFolder(path)
	i.saveConfig(false)
}

func (i *indexerView) btnStartClicked(sender *gowd.Element, event *gowd.EventElement) {
	err := indexer.Start(indexer.Options{
		WithLocation:   i.chkLocation.Checked(),
		DeleteDatabase:         i.chkDeleteDB.Checked(),
		WithFaces: i.chkWithFaces.Checked(),
		WithObjects: i.chkWithObjs.Checked(),
		QuickScan: i.chkQuicksacn.Checked(),
	})
	if err != nil {
		Root.addAlertError(err)
	}
	i.updateState()
}

func (i *indexerView) btnAddFolderChanged(sender *gowd.Element, event *gowd.EventElement) {
	path := i.btnAddFolder.GetValue()
	conf.AddSourceFolder(path)
	i.saveConfig(false)
}

func (i *indexerView) inputMapChanged(sender *gowd.Element, event *gowd.EventElement) {
	conf.Options.MapQuestAPIKey = i.inputMapQuest.GetValue()
	i.saveConfig(true)
}

func (i *indexerView) saveConfig(showAlert bool) {
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
