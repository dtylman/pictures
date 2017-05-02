package view

import (
	"fmt"
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	"github.com/dtylman/pictures/backuper"
	"github.com/dtylman/pictures/cmd/app/view/darktheme"
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/tasklog"
)

type manageView struct {
	*gowd.Element
	progressBar       *bootstrap.ProgressBar
	inputBackupFolder *bootstrap.FormInput
	dialogRFS         *darktheme.Dialog
	inputSourceFolder *gowd.Element
}

func newManageView() *manageView {
	mv := new(manageView)
	mv.Element = gowd.NewElement("div")

	mv.inputBackupFolder = bootstrap.NewFormInput(bootstrap.InputTypeText, "Backup folder")

	pnl := bootstrap.NewPanel(bootstrap.PanelDefault)
	pnl.AddTitle("Backup")
	pnl.AddToBody(mv.inputBackupFolder.Element)

	mv.AddElement(pnl.Element)

	mv.progressBar = bootstrap.NewProgressBar()
	mv.progressBar.SetAttribute("style", "height: 40px")

	mv.AddElement(gowd.NewStyledText("Progress:", gowd.Paragraph))
	mv.AddElement(mv.progressBar.Element)

	tasklog.RegisterHandler(tasklog.BackuperTask, mv.updateProgress)

	mv.dialogRFS = darktheme.NewDialog("Remove From Source")
	mv.dialogRFS.Body.AddHtml(`<p>Select a source folder to delete indexed images from<p>`)
	mv.inputSourceFolder=gowd.NewElement("input")
	mv.inputSourceFolder.SetAttribute("type","file")
	mv.inputSourceFolder.SetAttribute("nwdirectory","")
	mv.dialogRFS.Body.AddElement(mv.inputSourceFolder)
	btn := bootstrap.NewButton(bootstrap.ButtonPrimary, "OK")
	btn.OnEvent(gowd.OnClick, mv.btnRemoveFromSrcClicked)
	mv.dialogRFS.Footer.AddElement(btn)

	mv.AddElement(mv.dialogRFS.Element)
	return mv
}

func (bv *manageView) updateState() {
	bv.inputBackupFolder.SetValue(conf.Options.BackupFolder)
}

func (bv *manageView) populateToolbar(menu *darktheme.Menu) {
	menu.AddButton(menu.TopLeft, "Backup", "fa fa-play", bv.btnBackupClicked)
	btn:=menu.AddButton(menu.TopLeft, "Remove from source", "fa fa-remove", nil)
	btn.SetAttribute("data-toggle","modal")
	btn.SetAttribute("data-target","#"+bv.dialogRFS.GetID())

	menu.AddButton(menu.TopLeft, "Stop", "fa fa-stop", bv.btnStopClicked)
}

func (bv *manageView) getContent() *gowd.Element {
	return bv.Element
}

func (bv *manageView) btnRemoveFromSrcClicked(sender *gowd.Element, event *gowd.EventElement) {
	sourceFolder:=bv.inputSourceFolder.GetValue()

	bv.AddElement(gowd.NewStyledText("This will erase all pictures from "+sourceFolder,gowd.BoldText))
}

func (bv *manageView) btnBackupClicked(sender *gowd.Element, event *gowd.EventElement) {
	conf.Options.BackupFolder = bv.inputBackupFolder.GetValue()
	err := conf.Save()
	if err != nil {
		Root.addAlertError(err)
		return
	}
	err = backuper.Start()
	if err != nil {
		Root.addAlertError(err)
	}
}

func (bv *manageView) btnStopClicked(sender *gowd.Element, event *gowd.EventElement) {
	backuper.Stop()
}

func (iv *manageView) updateProgress(status tasklog.Task) {
	if !status.Running {
		iv.progressBar.SetText("")
		iv.progressBar.SetValue(0, 0)
	} else {
		if status.Pos != 0 {
			iv.progressBar.SetText(fmt.Sprintf("%d / %d", status.Pos, status.Total))
			iv.progressBar.SetValue(status.Pos, status.Total)
		}
	}
	Root.Render()
}
