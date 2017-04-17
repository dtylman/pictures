package tasklog

import (
	"sync"
)

const (
	IndexerTask = "Indexer"
	BackuperTask = "Backuper"
)

type Task struct {
	Name     string
	Pos      int
	Total    int
	Messages []string
	Running  bool
}

type TaskHandler func(status Task)

var (
	handlers map[string][]TaskHandler
	mutex sync.Mutex
)

func init() {
	handlers = make(map[string][]TaskHandler)
}

func RegisterHandler(taskName string, handler TaskHandler) {
	mutex.Lock()
	defer mutex.Unlock()
	list, exists := handlers[taskName]
	if !exists {
		list = make([]TaskHandler, 0)
	}
	list = append(list, handler)
	handlers[taskName] = list
}

func getHandlers(taskName string) []TaskHandler {
	mutex.Lock()
	defer mutex.Unlock()
	return handlers[taskName]
}

func Status(taskName string, running bool, pos int, total int, messages ...string) {
	task := Task{Name:taskName, Running:running, Pos:pos, Total:total, Messages:messages}
	for _, handler := range getHandlers(taskName) {
		handler(task)
	}
}

func StatusMessage(taskName string, messages...string) {
	Status(taskName, true, 0, 0, messages...)
}
