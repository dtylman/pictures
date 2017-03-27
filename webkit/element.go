package webkit

import (
	"fmt"
	"golang.org/x/net/html"
)

var order int

type Element struct {
	//Parent the parent element
	Parent *Element
	//eventHandlers holds event handlers for events
	eventHandlers map[string]EventHandler
	//Kids child elements
	Kids []*Element
	//Attributes element attributes...
	Attributes []html.Attribute
	//Object arbitrary user object that can be associated with element.
	Object interface{}
	//nodeType the element node type
	nodeType html.NodeType
	//data html tag or text
	data string
	//Hidden if true the element will not be rendered
	Hidden bool
	//renderHash after rendered get a hash to identify if should be rendered again.
	renderHash []byte
}

func NewElement(tag string) *Element {
	elem := &Element{
		data:          tag,
		Attributes:    make([]html.Attribute, 0),
		nodeType:      html.ElementNode,
		Kids:          make([]*Element, 0),
		eventHandlers: make(map[string]EventHandler),
	}
	order++
	elem.setID(fmt.Sprintf("_%s%d", elem.data, order))
	return elem
}

func (e *Element) SetAttributes(event *EventElement) {
	for key, value := range event.Attributes {
		e.SetAttribute(key, value)
	}
}

func (e *Element) AddElement(elem *Element) *Element {
	if elem == nil {
		panic("Cannot add nil element")
	}
	elem.Parent = e
	e.Kids = append(e.Kids, elem)
	return elem
}

func (e *Element) RemoveElements() {
	e.Kids = make([]*Element, 0)
}

func (e *Element) RemoveElement(elem *Element) {
	before := e.Kids
	e.RemoveElements()
	for _, kid := range before {
		if kid.GetID() != elem.GetID() {
			e.AddElement(kid)
		}
	}
}

func (e *Element) SetText(text string) {
	if e.nodeType == html.TextNode {
		e.data = text
	} else {
		e.AddElement(NewText(text))
	}
}

func (e *Element) SetAttribute(key, val string) {
	for i := range e.Attributes {
		if e.Attributes[i].Key == key {
			e.Attributes[i].Val = val
			return
		}
	}
	e.Attributes = append(e.Attributes, html.Attribute{Key: key, Val: val})
}

func (e *Element) GetAttribute(key string) (string, bool) {
	if e.Attributes == nil {
		return "", false
	}
	for _, a := range e.Attributes {
		if a.Key == key {
			return a.Val, true
		}
	}
	return "", false
}

func (e *Element) GetValue() string {
	val, _ := e.GetAttribute("value")
	return val
}

func (e *Element) GetID() string {
	val, _ := e.GetAttribute("id")
	return val
}

func (e *Element) setID(id string) {
	e.SetAttribute("id", id)
}

func (e *Element) Hide() {
	e.Hidden = true
}

func (e *Element) Show() {
	e.Hidden = false
}

func (e *Element) Find(id string) *Element {
	if e.GetID() == id {
		return e
	}
	for i := range e.Kids {
		elem := e.Kids[i].Find(id)
		if elem != nil {
			return elem
		}
	}
	return nil
}

func (e *Element) OnEvent(event string, handler EventHandler) {
	e.SetAttribute(event, fmt.Sprintf(`fire_event('%s',this);`, event))
	e.eventHandlers[event] = handler
}

func (e *Element) toNode() *html.Node {
	node := &html.Node{Attr: e.Attributes, Data: e.data, Type: e.nodeType}
	if e.Hidden {
		node.Type = html.CommentNode // to a void returning null which may crash the caller
		return node
	}
	for i := range e.Kids {
		node.AppendChild(e.Kids[i].toNode())
	}
	return node
}

func (e *Element) Render() error {
	return render(e)
}

func (e *Element) ProcessEvent(event *Event) {
	for _, input := range event.Inputs {
		elementID := input.GetID()
		if elementID != "" {
			e.updateState(elementID, &input)
		}
	}
	e.fireEvent(event.Name, event.Sender.GetID(), &event.Sender)
}

func (e *Element) updateState(elementID string, input *EventElement) {
	for i := range e.Kids {
		e.Kids[i].updateState(elementID, input)
	}
	if e.GetID() == elementID {
		e.SetAttributes(input)
		e.SetAttribute("value", input.Value)
	}
}

func (e *Element) fireEvent(eventName string, senderID string, sender *EventElement) {
	if senderID == "" {
		return
	}
	kids := e.Kids //events may change the kids container
	for i := range kids {
		kids[i].fireEvent(eventName, senderID, sender)
	}
	if e.GetID() == senderID {
		handler, exists := e.eventHandlers[eventName]
		if exists {
			handler(e, sender)
		}
	}
}

func SetAttribute(node *html.Node, key string, val string) {
	if node.Attr == nil {
		node.Attr = make([]html.Attribute, 0)
	}
	for i := range node.Attr {
		if node.Attr[i].Key == key {
			node.Attr[i].Val = val
			return
		}
	}
	node.Attr = append(node.Attr, html.Attribute{Key: key, Val: val})
}
