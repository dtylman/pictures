package indexer

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

type IndexError struct {
	Path  string
	Error error
}

type Indexer struct {
	running  bool
	options  Options
	errors   []IndexError
	progress Progress
	mutex    sync.Mutex
}

func (i *Indexer) reset() {
	i.running = false
	i.options = Options{}
	i.errors = make([]IndexError, 0)
	i.progress.reset()
}

func (i *Indexer) IsRunning() bool {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	return i.running
}

func (i *Indexer) SetStarted(options Options) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if i.running {
		return errors.New("Indexer already running")
	}
	i.reset()
	i.running = true
	i.options = options
	return nil
}

func (i *Indexer) SetDone() {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.running = false
}

func (i *Indexer) GetOptions() Options {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	return i.options
}

func (i *Indexer) AddError(path string, err error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.errors = append(i.errors, IndexError{Path: path, Error: err})
	log.Println(path, err)
}

func (i *Indexer) ProgressString() string {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	return fmt.Sprintf("%s: %d files (%v Kbs)", i.progress.RootDir, i.progress.TotalFiles, i.progress.TotalSize/1024)
}
