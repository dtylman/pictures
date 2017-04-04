package indexer

import (
	"errors"
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/tasklog"
	"sync"
)

type ProgressHandler func(progress IndexerProgress)

type Options struct {
	//IndexLocation if true will do include geolocation
	IndexLocation bool
	//ReIndex if true will delete previous results
	ReIndex bool
	//ProgressHandler sets a handler to get information about progress
	ProgressHandler ProgressHandler
}

type IndexError struct {
	Path  string
	Error error
}

type Indexer struct {
	IndexerProgress
	options Options
	mutex   sync.Mutex
}

func (i *Indexer) reset() {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.Running = false
	i.options = Options{}
	i.Errors = make([]IndexError, 0)
	i.CurrentFolder = ""
	i.TotalFiles = 0
	i.TotalSize = 0
	i.TotalRootFolders = 0
	i.CurrentRootFolder = 0
}

func (i *Indexer) IsRunning() bool {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	return i.Running
}

func (i *Indexer) initialize(options Options) error {
	i.reset()
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if i.Running {
		return errors.New("Indexer already running")
	}
	i.Running = true
	i.options = options
	i.TotalRootFolders = len(conf.Options.SourceFolders)
	return nil
}

func (i *Indexer) setDone() {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.Running = false
}

func (i *Indexer) incSourceFolder() {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.CurrentRootFolder++
}

func (i *Indexer) incFile(size int64) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.TotalFiles++
	i.TotalSize += size
}

func (i *Indexer) setCurrentFolder(path string) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.CurrentFolder = path
}

func (i *Indexer) GetOptions() Options {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	return i.options
}

func (i *Indexer) AddError(path string, err error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.Errors = append(i.Errors, IndexError{Path: path, Error: err})
	tasklog.Println(path, err)
}

func (i *Indexer) fireProgressEvent() {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if i.options.ProgressHandler != nil {
		go i.options.ProgressHandler(i.IndexerProgress)
	}
}
