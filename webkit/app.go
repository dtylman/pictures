package webkit

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"os"
	"sync"
	"time"
)

var renderMutex sync.Mutex

func render(e *Element) error {
	node := e.toNode()
	renderMutex.Lock()
	defer renderMutex.Unlock()
	err := html.Render(os.Stdout, node)
	if err != nil {
		return err
	}
	_, err = fmt.Println()
	return err
}

func Run(body *Element) error {
	for true {
		err := body.Render()
		if err != nil {
			return err
		}
		decoder := json.NewDecoder(os.Stdin)
		var event Event
		err = decoder.Decode(&event)
		if err != nil {
			return err
		}
		body.ProcessEvent(&event)
	}
	return nil
}

func Error(err error) {
	fmt.Println("Error: ", err.Error())
	time.Sleep(3 * time.Second)
}
