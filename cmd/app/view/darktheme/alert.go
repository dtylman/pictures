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
	em:=gowd.NewElementMap()
	for key, alert := range a.alerts {
		a.AddHtml(
			`<div class="col-lg-6">
				<div class="alert alert-danger" id="pnlAlert">
					<strong>Error: </strong> ` + alert +
				`</div>`,em)

		btn := bootstrap.NewButton("button", "x")
		btn.SetClass("close")
		btn.OnEvent(gowd.OnClick, func(sender *gowd.Element, event *gowd.EventElement) {
			delete(a.alerts, key)
			a.updateState()
		})
		em["pnlAlert"].AddElement(btn)
	}
}
