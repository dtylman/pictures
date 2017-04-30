package view

import "github.com/dtylman/gowd"

type toolbarButton struct {
	caption string
	icon    string
	handler gowd.EventHandler
}

type toolbar struct {
	buttons []toolbarButton
}

func newToolbar() *toolbar {
	return &toolbar{buttons:make([]toolbarButton, 0)}
}

func (t*toolbar) add(caption string, icon string, handler gowd.EventHandler) {
	t.buttons = append(t.buttons, toolbarButton{caption:caption, icon:icon, handler:handler})
}