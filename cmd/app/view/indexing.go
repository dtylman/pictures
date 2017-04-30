package view

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	"github.com/dtylman/pictures/indexer"
	"github.com/dtylman/pictures/tasklog"
	"fmt"
	"github.com/dtylman/pictures/cmd/app/view/darktheme"
)

type indexingView struct {
	*gowd.Element
	pnlStatus   *bootstrap.Panel
	progressBar *bootstrap.ProgressBar
}

func newIndexingView() *indexingView {
	iv := new(indexingView)
	iv.Element = gowd.NewElement("div")

	iv.pnlStatus = bootstrap.NewPanel(bootstrap.PanelDefault)
	iv.pnlStatus.AddTitle("Indexing...")

	iv.progressBar = bootstrap.NewProgressBar()
	iv.progressBar.SetAttribute("style", "height: 40px")

	iv.AddElement(iv.pnlStatus.Element)
	iv.AddElement(gowd.NewStyledText("Progress:", gowd.Paragraph))
	iv.AddElement(iv.progressBar.Element)

	tasklog.RegisterHandler(tasklog.IndexerTask, iv.updateIndexerProgress)

	return iv
}

func (iv *indexingView) btnStopClicked(sender *gowd.Element, event *gowd.EventElement) {
	indexer.Stop()
}

func (iv *indexingView) updateState() {
	if !indexer.IsRunning() {
		iv.progressBar.SetText("")
		iv.progressBar.SetValue(0, 0)
		Root.setActiveView(Root.indexer)
	}
}

func (iv *indexingView) populateToolbar(menu*darktheme.Menu) {
	menu.AddTopButton("Stop", "fa fa-stop", iv.btnStopClicked)
}

func (iv *indexingView) updateIndexerProgress(status tasklog.Task) {
	if !status.Running {
		iv.progressBar.SetText("")
		iv.progressBar.SetValue(0, 0)
		Root.setActiveView(Root.indexer)
	} else {
		iv.pnlStatus.Body.RemoveElements()
		for _, msg := range status.Messages {
			iv.pnlStatus.AddToBody(gowd.NewStyledText(msg, gowd.SmallText))
		}
		if status.Pos != 0 {
			iv.progressBar.SetText(fmt.Sprintf("%d / %d", status.Pos, status.Total))
			iv.progressBar.SetValue(status.Pos, status.Total)
		}
	}
	Root.Render()
}

func (iv *indexingView) getContent() *gowd.Element {
	return iv.Element
}