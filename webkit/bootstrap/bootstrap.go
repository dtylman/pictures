package bootstrap

import (
	"github.com/dtylman/pictures/webkit"
)

func NewContainer(fluid bool) *webkit.Element {
	if fluid {
		return webkit.NewElement("div", "container-fluid")
	}
	return webkit.NewElement("div", "container")
}

func NewForm() *webkit.Element {
	return webkit.NewElement("form", "")
}

func NewFormGroup() *webkit.Element {
	return webkit.NewElement("div", "form-group")
}

func NewLabel(forwhat string, text string) *webkit.Element {
	label := webkit.NewElement("label", "")
	label.SetAttribute("for", forwhat)
	label.Data = text
	return label
}

func NewInput(placeholder string) *webkit.Element {
	input := webkit.NewElement("input", "form-control")
	input.SetAttribute("type", "text")
	if placeholder != "" {
		input.SetAttribute("placeholder", placeholder)
	}
	return input
}

func NewButton(buttontype string, caption string) *webkit.Element {
	btn := webkit.NewElement("button", "btn "+buttontype)
	btn.SetText("caption")
	return btn
}
