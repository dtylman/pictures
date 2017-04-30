package view

import (
	"fmt"
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	"github.com/dtylman/pictures/indexer/picture"
	"github.com/dtylman/pictures/cmd/app/view/darktheme"
)

type image struct {
	*gowd.Element
}

func newImage() *image {
	i := new(image)
	i.Element = bootstrap.NewElement("div", "well")

	i.updateState()

	return i
}

func (i *image) populateToolbar(menu*darktheme.Menu) {
	menu.AddButton(menu.TopLeft,"Prev", "fa fa-arrow-left", i.btnPrevClicked)
	menu.AddButton(menu.TopLeft,"Next", "fa fa-arrow-right", i.btnNextClicked)
}

func (i *image) btnPrevClicked(sender *gowd.Element, event *gowd.EventElement) {
	activeSearch.PrevImage()
	i.updateState()
}

func (i *image) btnNextClicked(sender *gowd.Element, event *gowd.EventElement) {
	activeSearch.NextImage()
	i.updateState()
}

func (i *image) updateState() {
	i.Element.RemoveElements()
	col := bootstrap.NewColumn(bootstrap.ColumnLarge, 9)
	row := bootstrap.NewRow()
	row.AddElement(col)
	mimeType := activeSearch.ActiveImage.MimeType
	if picture.MimeIs(mimeType, picture.Image) {
		img := bootstrap.NewElement("img", "img-responsive")
		img.SetAttribute("src", fmt.Sprintf("file:///%s", activeSearch.ActiveImage.Path))

		col.AddElement(img)
	} else if picture.MimeIs(mimeType, picture.Video) {
		vid := bootstrap.NewElement("video", "")
		vid.SetAttribute("src", fmt.Sprintf("file:///%s", activeSearch.ActiveImage.Path))
		vid.SetAttribute("type", mimeType)
		col.AddElement(vid)

	}

	table := bootstrap.NewTable(bootstrap.TableStripped)
	col = bootstrap.NewColumn(bootstrap.ColumnLarge, 3)
	col.AddElement(table.Element)
	row.AddElement(col)

	pnl := bootstrap.NewPanel(bootstrap.PanelDefault)
	pnl.AddTitle(activeSearch.ActiveImage.Path)
	pnl.AddToBody(row)

	i.AddElement(pnl.Element)
}
func (i *image) getContent() *gowd.Element {
	return i.Element
}
