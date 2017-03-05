package runningindexer

import (
	"errors"
	"sync"
)

type IndexError struct {
	Path  string
	Error error
}

var status struct {
	running bool
	options Options
	errors  []IndexError
}

var mutex sync.Mutex

func init() {
	resetStatus()
}

func resetStatus() {
	status.running = false
	status.options = Options{}
	status.errors = make([]IndexError, 0)
}

func IsRunning() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return status.running
}

func SetStarted(options Options) error {
	mutex.Lock()
	defer mutex.Unlock()
	if status.running {
		return errors.New("Indexer already running")
	}
	resetStatus()
	status.running = true
	status.options = options
	return nil
}

func SetDone() {
	mutex.Lock()
	defer mutex.Unlock()
	status.running = false
}

func GetOptions() Options {
	mutex.Lock()
	defer mutex.Unlock()
	return status.options
}

func AddError(path string, err error) {
	mutex.Lock()
	defer mutex.Unlock()
	status.errors = append(status.errors, IndexError{Path: path, Error: err})
}
