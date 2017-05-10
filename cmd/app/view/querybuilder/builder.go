package querybuilder

import "github.com/dtylman/gowd"

type Builder struct {
	*gowd.Element
}

func NewBuilder() *Builder {
	b := new(Builder)
	b.Element = gowd.NewElement("div")
	b.AddElement(newRow().Element)
	return b
}
