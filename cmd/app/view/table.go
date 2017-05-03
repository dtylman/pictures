package view

import (
	"fmt"
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	"github.com/dtylman/pictures/cmd/app/view/darktheme"
	"github.com/dtylman/pictures/indexer/thumbs"
	"log"
)

type table struct {
	*gowd.Element
}

func newTableView() *table {
	t := new(table)
	t.Element = gowd.NewElement("div")
	return t
}

func (t *table) updateState() {
	t.RemoveElements()
	if activeSearch == nil {
		t.AddElement(gowd.NewStyledText("Need to search first", gowd.EmphasizedText))
		return
	}
	table := bootstrap.NewTable(bootstrap.TableStripped)
	t.AddElement(table.Element)

	table.AddHeader("")
	table.AddHeader("Album")
	table.AddHeader("Path")
	table.AddHeader("File Time")
	table.AddHeader("Taken")
	table.AddHeader("Location")
	table.AddHeader("Objects")
	table.AddHeader("Faces")
	table.AddHeader("")

	for _, image := range activeSearch.Results {
		row := table.AddRow()
		td := gowd.NewElement("td")
		td.AddElement(bootstrap.NewInput(bootstrap.InputTypeCheckbox))
		row.AddElement(td)
		row.AddElement(table.NewCell(image.Album))
		row.AddElement(table.NewCell(image.Path))
		row.AddElement(table.NewCell(image.FileTime.String()))
		row.AddElement(table.NewCell(image.Taken.String()))
		row.AddElement(table.NewCell(image.Location))
		row.AddElement(table.NewCell(image.Objects))
		row.AddElement(table.NewCell(image.Faces))
		img := gowd.NewElement("img")
		thumb, err := thumbs.MakeThumb(image.Path, image.MD5, false)
		if err != nil {
			log.Println(err)
		}
		img.SetAttribute("src", fmt.Sprintf("file://"+thumb))
		img.SetAttribute("style", "height: 60px;")
		td = gowd.NewElement("td")
		td.AddElement(img)
		row.AddElement(td)
	}
}

func (t *table) getContent() *gowd.Element {
	return t.Element
}

func (t *table) populateToolbar(menu *darktheme.Menu) {

}
