package indexview

import (
	"github.com/dtylman/pictures/webkit"
	"github.com/dtylman/pictures/webkit/bootstrap"
)

//<div class="input-group">
//<input type="text" class="form-control" placeholder="Search for...">
//<span class="input-group-btn">
//<button class="btn btn-default" type="button">Go!</button>
//</span>
//</div><!-- /input-group -->

type SourceFolder struct {
	*webkit.Element
	ButtonRemove *webkit.Element
}

func NewSourceFolder(path string) *SourceFolder {
	sf := new(SourceFolder)
	sf.Element = bootstrap.NewElement("div", "input-group")
	input := bootstrap.NewElement("input", "form-control")
	input.SetAttribute("readonly", "true")
	input.SetAttribute("value", path)
	sf.ButtonRemove = bootstrap.NewButton(bootstrap.ButtonDefault, "Remove")
	span := bootstrap.NewElement("span", "input-group-btn")
	span.AddElement(sf.ButtonRemove)
	sf.AddElement(input)
	sf.AddElement(span)
	return sf
}
