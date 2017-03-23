package bootstrap

import (
	"github.com/dtylman/pictures/webkit"
)

func NewElement(tag, class string) *webkit.Element {
	elem := webkit.NewElement(tag)
	elem.SetAttribute("class", class)
	return elem
}

func NewContainer(fluid bool) *webkit.Element {
	if fluid {
		return NewElement("div", "container-fluid")
	}
	return NewElement("div", "container")
}
