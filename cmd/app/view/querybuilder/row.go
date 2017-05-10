package querybuilder

import (
	"github.com/dtylman/gowd"
)


type Row struct {
	*gowd.Element

}

func newRow() *Row {
	r := new(Row)
	r.Element = gowd.NewElement("div")
	r.AddElement(gowd.NewText("Picture "))
	r.AddHtml(`<div class="btn-group">
                  <button type="button" class="btn btn-default">Default</button>
                  <button type="button" class="btn btn-default dropdown-toggle" data-toggle="dropdown"><span class="caret"></span></button>
                  <ul class="dropdown-menu">
                    <li><a href="#">Action</a></li>
                    <li><a href="#">Another action</a></li>
                    <li><a href="#">Something else here</a></li>
                    <li class="divider"></li>
                    <li><a href="#">Separated link</a></li>
                  </ul>
                </div>`,nil)
	return r
}