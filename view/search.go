package view

import (
	"github.com/dtylman/pictures/webkit"
	"github.com/dtylman/pictures/webkit/bootstrap"
)

type search struct {
	*webkit.Element
	parent parentControls
}

func newSearch(p parentControls) *search {
	sv := new(search)
	sv.Element = bootstrap.NewContainer(true)
	sv.parent = p
	return sv
}
