package darktheme

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
)

type Alerts struct {
	*gowd.Element
	alerts map[int]string
	key    int
}

func NewAlerts() *Alerts {
	a := new(Alerts)
	a.Element = gowd.NewElement("row")
	a.alerts = make(map[int]string)
	return a
}

func (a *Alerts) Add(text string) {
	a.key++
	a.alerts[a.key] = text
	a.updateState()
}

func (a *Alerts) updateState() {
	a.RemoveElements()
	for key, alert := range a.alerts {
		elems, _ := a.AddHtml(`<div class="col-lg-6">
						<div class="alert alert-danger">
						<strong>Error: </strong> ` + alert + `</div>`)

		btn := bootstrap.NewButton("button", "x")
		btn.SetClass("close")
		btn.OnEvent(gowd.OnClick, func(sender *gowd.Element, event *gowd.EventElement) {
			delete(a.alerts, key)
			a.updateState()
		})
		elems[0].Kids[1].AddElement(btn)
	}
}
