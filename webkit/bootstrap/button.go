package bootstrap

import "github.com/dtylman/pictures/webkit"

const (
	ButtonDefault = "btn btn-default"
	ButtonPimary  = "btn btn-primary"
)

func NewButton(buttontype string, caption string) *webkit.Element {
	btn := NewElement("button", "btn "+buttontype)
	btn.SetText(caption)
	return btn
}
