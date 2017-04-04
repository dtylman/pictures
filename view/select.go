package view

import "github.com/dtylman/gowd"

type Select struct {
	*gowd.Element
}

func NewSelect() *Select {
	s := new(Select)
	s.Element = gowd.NewElement("select")
	return s
}

func (s *Select) AddOption(caption, value string) *gowd.Element {
	option := gowd.NewElement("option")
	option.AddElement(gowd.NewText(caption))
	option.SetAttribute("value", value)
	s.AddElement(option)
	return option
}
