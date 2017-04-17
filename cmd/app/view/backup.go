package view

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/backuper"
	"github.com/dtylman/pictures/tasklog"
	"fmt"
)

type backupView struct {
	*gowd.Element
	btnStart          *gowd.Element
	progressBar       *bootstrap.ProgressBar
	btnStop           *gowd.Element
	inputBackupFolder *bootstrap.FormInput
}

func newBackupView() *backupView {
	bv := new(backupView)
	bv.Element = gowd.NewElement("div")

	bv.btnStart = bootstrap.NewButton(bootstrap.ButtonPrimary, "Start")
	bv.btnStart.OnEvent(gowd.OnClick, bv.btnStartClick)

	bv.inputBackupFolder = bootstrap.NewFormInput(bootstrap.InputTypeText, "Backup folder")

	pnl := bootstrap.NewPanel(bootstrap.PanelDefault)
	pnl.AddTitle("Backup")
	pnl.AddToBody(bv.inputBackupFolder.Element)

	bv.AddElement(pnl.Element)

	bv.progressBar = bootstrap.NewProgressBar()
	bv.progressBar.SetAttribute("style", "height: 40px")
	bv.btnStop = bootstrap.NewButton(bootstrap.ButtonPrimary, "Stop")
	bv.btnStop.OnEvent(gowd.OnClick, bv.btnStopClicked)

	bv.AddElement(gowd.NewStyledText("Progress:", gowd.Paragraph))
	bv.AddElement(bv.progressBar.Element)

	tasklog.RegisterHandler(tasklog.BackuperTask, bv.updateProgress)
	return bv
}

func (bv *backupView) updateState() {
	bv.inputBackupFolder.SetValue(conf.Options.BackupFolder)
	if backuper.IsRunning() {
		bv.btnStart.Disable()
		bv.btnStop.Enable()
	} else {
		bv.btnStart.Enable()
		bv.btnStop.Disable()
	}
}

func (bv *backupView) populateToolbar(toolbar *gowd.Element) {
	toolbar.AddElement(bv.btnStart)
	toolbar.AddElement(bv.btnStop)
}

func (bv *backupView) getContent() *gowd.Element {
	return bv.Element
}

func (bv *backupView) btnStartClick(sender *gowd.Element, event *gowd.EventElement) {
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

func (bv*backupView) btnStopClicked(sender *gowd.Element, event *gowd.EventElement) {
	backuper.Stop()
}

func (iv *backupView) updateProgress(status tasklog.Task) {
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