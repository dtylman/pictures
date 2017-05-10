package darktheme

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
)

type Dialog struct {
	*gowd.Element
	Header *gowd.Element
	Body   *gowd.Element
	Footer *gowd.Element
}

func NewDialog(title string) *Dialog {
	d := new(Dialog)
	d.Header = bootstrap.NewElement("div", "modal-header")
	d.Body = bootstrap.NewElement("div", "modal-body")
	d.Footer = bootstrap.NewElement("div", "modal-footer")
	d.Element = bootstrap.NewElement("div", "modal fade",
		bootstrap.NewElement("div", "modal-dialog",
			bootstrap.NewElement("div", "modal-content", d.Header, d.Body, d.Footer)))
	d.Element.SetAttribute("role", "dialog")
	d.Header.AddHtml(`<button type="button" class="close" data-dismiss="modal">&times;</button>`,nil)
	d.Header.AddElement(gowd.NewStyledText(title, gowd.Heading4))
	d.Footer.AddHtml(`<button type="button" class="btn btn-default" data-dismiss="modal">Close</button>`,nil)
	return d
}
