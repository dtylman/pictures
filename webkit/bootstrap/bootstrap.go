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

func NewLabel(forwhat string, text string) *webkit.Element {
	label := webkit.NewElement("label")
	label.SetAttribute("for", forwhat)
	label.AddElement(webkit.NewText(text))
	return label
}

func NewInput(placeholder string) *webkit.Element {
	input := NewElement("input", "form-control")
	input.SetAttribute("type", "text")
	if placeholder != "" {
		input.SetAttribute("placeholder", placeholder)
	}
	return input
}

func NewLinkButton(caption string) *webkit.Element {
	linkBtn := webkit.NewElement("a")
	linkBtn.SetAttribute("href", "#")
	linkBtn.AddElement(webkit.NewText(caption))
	return linkBtn
}
