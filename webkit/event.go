package webkit

type EventElement struct {
	Value      string            `json:"value"`
	Attributes map[string]string `json:"attributes"`
}

type Event struct {
	Sender EventElement   `json:"sender"`
	Inputs []EventElement `json:"inputs"`
}

type EventHandler func(sender *Element, event *EventElement)

const (
	OnClick = "onclick"
)

func (e *EventElement) GetID() string {
	id, _ := e.Attributes["id"]
	return id
}
