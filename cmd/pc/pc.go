package main

import (
	"encoding/json"
	"github.com/dtylman/pictures/webkit"
	"github.com/dtylman/pictures/webkit/bootstrap"
	"os"
)

type Page struct {
	*webkit.Body
	btn *webkit.Element
	txt *webkit.Element
}

func (p *Page) init() {
	p.Body = webkit.NewBody(bootstrap.NewContainer(true))
	p.btn = bootstrap.NewButton("btn-default", "click")
	p.btn.RegisterHandler("onclick", p.btnOnClick)
	p.txt = bootstrap.NewInput("hoho")
	p.AddElement(p.txt)
	p.AddElement(p.btn)
}

func (p *Page) btnOnClick(sender *webkit.EventElement) {
	p.AddElement(bootstrap.NewInput("hoho"))
	val, _ := p.txt.GetAttribute("value")
	p.btn.SetText(val)
}

func main() {
	var p Page
	p.init()
	for true {
		p.Render()
		decoder := json.NewDecoder(os.Stdin)
		var event webkit.Event
		decoder.Decode(&event)
		p.ProcessEvent(&event)
	}
}
