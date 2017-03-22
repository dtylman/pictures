package webkit

import (
	"fmt"
	"golang.org/x/net/html"
	"os"
)

type Body struct {
	container *Element
	id        int
	Kids      map[string]*Element
}

func NewBody(container *Element) *Body {
	b := &Body{container: container}
	b.Kids = make(map[string]*Element)
	return b
}

func (b *Body) AddElement(elem *Element) {
	elemID := b.idFor(elem)
	elem.SetAttribute("id", elemID)
	b.Kids[elemID] = elem
	b.container.AppendChild(elem.Node)
}

func (b *Body) idFor(elem *Element) string {
	b.id++
	return fmt.Sprintf("_%s%d", elem.Data, b.id)
}

func (b *Body) ProcessEvent(event *Event) {
	for _, input := range event.Inputs {
		kid := b.findElementFromEventElement(&input)
		if kid != nil {
			kid.SetAttributes(&input)
			kid.SetAttribute("value", input.Value)
		}
	}
	b.fireEvent(&event.Sender)
}

func (b *Body) fireEvent(sender *EventElement) {
	kid := b.findElementFromEventElement(sender)
	if kid != nil {
		if kid.onEvent != nil {
			kid.onEvent(sender)
		}
	}
}

func (b *Body) findElementFromEventElement(elem *EventElement) *Element {
	id, exists := elem.Attributes["id"]
	if exists {
		kid, exists := b.Kids[id]
		if exists {
			return kid
		}
	}
	return nil
}

func (b *Body) Render() error {
	return html.Render(os.Stdout, b.container.Node)
}
