package indexview

import (
	"github.com/dtylman/pictures/webkit"
	"github.com/dtylman/pictures/webkit/bootstrap"
)

var (
	View        *webkit.Element
	ProgressBar *bootstrap.ProgressBar
	BtnStop     *webkit.Element
	BtnStart    *webkit.Element
	ChkLocation *bootstrap.Checkbox
	ChkReindex  *bootstrap.Checkbox
)

func init() {
	View = webkit.NewElement("div")
	ProgressBar = bootstrap.NewProgressBar()
	BtnStop = bootstrap.NewButton(bootstrap.ButtonPimary, "Stop")
	ChkLocation = bootstrap.NewCheckBox("Include Locations", false)
	ChkReindex = bootstrap.NewCheckBox("Reindex Existing Items", false)
	BtnStart = bootstrap.NewButton(bootstrap.ButtonPimary, "Start")

	View.AddElement(ProgressBar.Element)
	View.AddElement(ChkLocation.Element)
	View.AddElement(ChkReindex.Element)
	View.AddElement(BtnStart)
	View.AddElement(BtnStop)
}
