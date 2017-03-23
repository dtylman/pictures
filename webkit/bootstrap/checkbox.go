package bootstrap

import (
	"github.com/dtylman/pictures/webkit"
	"strconv"
)

//<div class="checkbox">
//<label>
//<input type="checkbox"> Check me out
//</label>
//</div>

type Checkbox struct {
	*webkit.Element
	chkbox *webkit.Element
	txt    *webkit.Element
}

func NewCheckBox(caption string, checked bool) *Checkbox {
	cb := new(Checkbox)
	cb.Element = NewElement("div", "checkbox")
	lbl := webkit.NewElement("label")
	cb.chkbox = webkit.NewElement("input")
	cb.chkbox.SetAttribute("type", "checkbox")
	cb.chkbox.SetAttribute("checked", strconv.FormatBool(checked))
	lbl.AddElement(cb.chkbox)
	cb.txt = webkit.NewText(caption)
	lbl.AddElement(cb.txt)
	cb.Element.AddElement(lbl)
	return cb
}
