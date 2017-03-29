package view

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
)

type search struct {
	*gowd.Element
	parent parentControls
}

func newSearch(p parentControls) *search {
	sv := new(search)
	sv.Element = bootstrap.NewContainer(true)
	sv.parent = p
	return sv
}
