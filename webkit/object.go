package webkit

import (
	"fmt"
	"io"
)

type Renderable interface {
	Render(writer io.Writer) error
}

type Object struct {
	Parent     *Object
	Tag        string
	Attributes ObjectAttributes
	Components []Renderable
}

func NewObject(tag string, class string, id string) *Object {
	o := &Object{Tag: tag, Attributes: NewAttributes(), Components: make([]Renderable, 0)}
	o.Attributes[AttrClass] = class
	o.Attributes[AttrID] = id
	return o
}

func (o *Object) Render(writer io.Writer) error {
	_, err := fmt.Fprintf(writer, `<%s`, o.Tag)
	if err != nil {
		return err
	}
	err = o.Attributes.Render(writer)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(writer, ">")
	if err != nil {
		return err
	}
	for _, c := range o.Components {
		err = c.Render(writer)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintf(writer, `</%s>`, o.Tag)
	return err
}

func (o *Object) AddChild(r Renderable) {
	if o.Components == nil {
		o.Components = make([]Renderable, 0)
	}
	object, ok := r.(*Object)
	if ok {
		object.Parent = o
	}
	o.Components = append(o.Components, r)
}
