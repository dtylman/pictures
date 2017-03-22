package webkit

import (
	"golang.org/x/net/html"
)

type Element struct {
	*html.Node
	onEvent EventHandler
}

func NewElement(tag string, class string) *Element {
	elem := &Element{Node: new(html.Node)}
	elem.Type = html.ElementNode
	elem.Data = tag
	elem.Attr = make([]html.Attribute, 0)
	if class != "" {
		elem.SetAttribute("class", class)
	}
	return elem
}

func (e *Element) SetAttributes(event *EventElement) {
	e.Attr = make([]html.Attribute, 0)
	for key, value := range event.Attributes {
		e.SetAttribute(key, value)
	}
}

func (e *Element) SetText(text string) {
	if e.FirstChild == nil {
		e.AppendChild(NewText(text))
		return
	}
	if e.FirstChild.Type == html.TextNode {
		e.FirstChild.Data = text
	}
}

func (e *Element) SetAttribute(key, val string) {
	if e.Attr == nil {
		e.Attr = make([]html.Attribute, 0)
	}
	for i := range e.Attr {
		if e.Attr[i].Key == key {
			e.Attr[i].Val = val
			return
		}
	}
	e.Attr = append(e.Attr, html.Attribute{Key: key, Val: val})
}

func (e *Element) GetAttribute(key string) (string, bool) {
	if e.Attr == nil {
		return "", false
	}
	for _, a := range e.Attr {
		if a.Key == key {
			return a.Val, true
		}
	}
	return "", false
}

func (e *Element) RegisterHandler(event string, handler EventHandler) {
	e.SetAttribute(event, `fire_event(this);`)
	e.onEvent = handler
}

func NewText(text string) *html.Node {
	return &html.Node{Type: html.TextNode, Data: text}
}
