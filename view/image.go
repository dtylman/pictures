package view

import (
	"fmt"
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	"github.com/dtylman/pictures/server/model"
)

type image struct {
	*gowd.Element
	activeSearch *model.Search
}

func newImage(activeSearch *model.Search) *image {
	i := new(image)
	i.Element = bootstrap.NewElement("div", "well")
	i.activeSearch = activeSearch

	img := gowd.NewElement("img")
	img.SetAttribute("src", fmt.Sprintf("file:///%s", activeSearch.ActiveImage.Path))

	imageLink := bootstrap.NewLinkButton("")
	imageLink.SetClass("thumbnail")
	imageLink.SetAttribute("width", "100%")
	imageLink.AddElement(img)

	col := bootstrap.NewColumn(bootstrap.ColumnLarge, 9)
	col.AddElement(imageLink)

	row := bootstrap.NewRow()
	row.AddElement(col)

	table := bootstrap.QuickTable("", i.activeSearch.ActiveImage.Details)
	col = bootstrap.NewColumn(bootstrap.ColumnLarge, 3)
	col.AddElement(table.Element)
	row.AddElement(col)

	pnl := bootstrap.NewPanel(bootstrap.PanelDefault)
	pnl.AddTitle(activeSearch.ActiveImage.Name)
	pnl.AddBody(row)

	i.AddElement(pnl.Element)
	return i
}
