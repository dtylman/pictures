package view

import (
	"fmt"
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	"github.com/dtylman/pictures/indexer/picture"
	"github.com/dtylman/pictures/model"
)

type image struct {
	*gowd.Element
	activeSearch *model.Search
}

func newImage(activeSearch *model.Search) *image {
	i := new(image)
	i.Element = bootstrap.NewElement("div", "well")
	i.activeSearch = activeSearch

	i.updateState()

	return i
}

func (i *image) populateToolbar(toolbar *gowd.Element) {
	btnPrev := bootstrap.NewButton(bootstrap.ButtonDefault, "Prev")
	btnPrev.OnEvent(gowd.OnClick, i.btnPrevClicked)
	toolbar.AddElement(bootstrap.NewColumn(bootstrap.ColumnLarge, 1, btnPrev))

	btnNext := bootstrap.NewButton(bootstrap.ButtonDefault, "Next")
	btnNext.OnEvent(gowd.OnClick, i.btnNextClicked)
	toolbar.AddElement(bootstrap.NewColumn(bootstrap.ColumnLarge, 1, btnNext))

}

func (i *image) btnPrevClicked(sender *gowd.Element, event *gowd.EventElement) {
	i.activeSearch.PrevImage()
	i.updateState()
}

func (i *image) btnNextClicked(sender *gowd.Element, event *gowd.EventElement) {
	i.activeSearch.NextImage()
	i.updateState()
}

func (i *image) updateState() {
	i.Element.RemoveElements()
	col := bootstrap.NewColumn(bootstrap.ColumnLarge, 9)
	row := bootstrap.NewRow()
	row.AddElement(col)
	mimeType := i.activeSearch.ActiveImage.MimeType
	if picture.MimeIs(mimeType, picture.Image) {
		img := bootstrap.NewElement("img", "img-responsive")
		img.SetAttribute("src", fmt.Sprintf("file:///%s", i.activeSearch.ActiveImage.Path))

		col.AddElement(img)
	} else if picture.MimeIs(mimeType, picture.Video) {
		vid := bootstrap.NewElement("video", "")
		vid.SetAttribute("src", fmt.Sprintf("file:///%s", i.activeSearch.ActiveImage.Path))
		vid.SetAttribute("type", mimeType)
		col.AddElement(vid)

	}

	table := bootstrap.NewTable(bootstrap.TableStripped)
	col = bootstrap.NewColumn(bootstrap.ColumnLarge, 3)
	col.AddElement(table.Element)
	row.AddElement(col)

	pnl := bootstrap.NewPanel(bootstrap.PanelDefault)
	pnl.AddTitle(i.activeSearch.ActiveImage.Path)
	pnl.AddToBody(row)

	i.AddElement(pnl.Element)
}
func (i *image) getContent() *gowd.Element {
	return i.Element
}
