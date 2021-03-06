package view

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
)

//<div class="input-group">
//<input type="text" class="form-control" placeholder="Search for...">
//<span class="input-group-btn">
//<button class="btn btn-default" type="button">Go!</button>
//</span>
//</div><!-- /input-group -->

//NewSourceFolder returns a new source folder element
func NewSourceFolder(path string, onRemove gowd.EventHandler) *gowd.Element {
	sf := bootstrap.NewElement("div", "input-group")
	input := bootstrap.NewElement("input", "form-control")
	input.SetAttribute("readonly", "true")
	input.SetAttribute("value", path)
	btnRemove := bootstrap.NewButton(bootstrap.ButtonDefault, "Remove")
	btnRemove.OnEvent(gowd.OnClick, onRemove)
	btnRemove.Object = path
	span := bootstrap.NewElement("span", "input-group-btn")
	span.AddElement(btnRemove)
	sf.AddElement(input)
	sf.AddElement(span)
	return sf
}
