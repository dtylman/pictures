package indexer

import (
	"errors"
	"fmt"
	"github.com/c2h5oh/datasize"
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/progressbar"
	"log"
	"sync"
)

type IndexError struct {
	Path  string
	Error error
}

type Indexer struct {
	running             bool
	options             Options
	errors              []IndexError
	CurrentFolder       string
	TotalFiles          int
	TotalSize           int64
	TotalSourceFolders  int
	CurrentSourceFolder int
	mutex               sync.Mutex
}

func (i *Indexer) reset() {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.running = false
	i.options = Options{}
	i.errors = make([]IndexError, 0)
	i.CurrentFolder = ""
	i.TotalFiles = 0
	i.TotalSize = 0
	i.TotalSourceFolders = 0
	i.CurrentSourceFolder = 0
}

func (i *Indexer) IsRunning() bool {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	return i.running
}

func (i *Indexer) SetStarted(options Options) error {
	i.reset()
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if i.running {
		return errors.New("Indexer already running")
	}
	i.running = true
	i.options = options
	i.TotalSourceFolders = len(conf.Options.SourceFolders)
	return nil
}

func (i *Indexer) setDone() {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.running = false
}

func (i *Indexer) incSourceFolder() {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.CurrentSourceFolder++
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
	i.errors = append(i.errors, IndexError{Path: path, Error: err})
	log.Println(path, err)
}

func (i *Indexer) ProgressStatus() *progressbar.Status {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	return &progressbar.Status{
		Text:       fmt.Sprintf("%s: %v files (%s)", i.CurrentFolder, i.TotalFiles, datasize.ByteSize(i.TotalSize).HumanReadable()),
		Done:       !i.running,
		Percentage: i.CurrentSourceFolder / i.TotalSourceFolders * 100,
	}
}
