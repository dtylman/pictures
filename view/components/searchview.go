package components

import (
	"github.com/dtylman/pictures/webkit"
	"github.com/dtylman/pictures/webkit/bootstrap"
)

type search struct {
	*webkit.Element
}

func newSearch() *search {
	sv := new(search)
	sv.Element = bootstrap.NewContainer(true)
	return sv
}
