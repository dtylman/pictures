package indexer

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/indexer/picture"
	"github.com/dtylman/pictures/indexer/remover"
	"github.com/dtylman/pictures/tasklog"
)

type Indexer struct {
	options Options
	running bool
	mutex   sync.Mutex
	images  *picture.Queue
}

func (i *Indexer) isRunning() bool {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	return i.running
}

func (i *Indexer) setRunning(value bool) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.running = value

}

func (i *Indexer) indexPictures() {
	defer func() {
		i.setRunning(false)
		tasklog.Status(tasklog.IndexerTask, false, 0, 0, "Done")
	}()
	log.Println("Starting index with options: ", i.options)
	if i.options.DeleteDatabase {
		tasklog.StatusMessage(tasklog.IndexerTask, "Deleting existing database...")
		err := i.deleteDB()
		if err != nil {
			tasklog.StatusMessage(tasklog.IndexerTask, err.Error())
		}
	}
	tasklog.StatusMessage(tasklog.IndexerTask, "Counting files...")
	for _, folder := range conf.Options.SourceFolders {
		err := filepath.Walk(folder, i.processFile)
		if err != nil {
			AddError(folder, err)
		}
	}
	tasklog.StatusMessage(tasklog.IndexerTask, fmt.Sprintf("Saving %d new items...", i.images.Length()))
	db.BatchIndex(i.images)
	tasklog.StatusMessage(tasklog.IndexerTask, "Checking for missing files...")
	err := remover.Remove()
	if err != nil {
		tasklog.Error(err)
	}

	imageProcessor := newProcessor(i.options)
	imageProcessor.update()

}

func (w *Indexer) start(options Options) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.running {
		return errors.New("Already running")
	}
	w.running = true
	w.options = options
	w.images = picture.NewQueue()
	if w.options.WithLocation {
		if conf.Options.MapQuestAPIKey == "" {
			return errors.New("API KEY for map quest is empty")
		}
	}
	go indexer.indexPictures()
	return nil
}

func (w *Indexer) processFile(path string, info os.FileInfo, e1 error) error {
	if e1 != nil {
		AddError(path, e1)
		return nil
	}
	if !IsRunning() {
		return errors.New("Indexer had stopped")
	}
	if info.IsDir() {
		tasklog.StatusMessage(tasklog.IndexerTask, fmt.Sprintf("Scanning folder '%s' (%d new images found)", path, w.images.Length()))
	} else {
		if info.Size() > 0 {
			if w.options.QuickScan {
				exists, err := db.HasPath(path)
				if err != nil {
					return err
				}
				if exists {
					return nil
				}
			}
			i, err := picture.NewIndex(path, info)
			if err != nil {
				AddError(path, err)
				return nil
			}
			if db.HasImage(i.MD5) {
				return nil
			}
			w.images.PushBack(i)
		}
	}
	return nil
}

func (w *Indexer) deleteDB() error {
	err := db.DeleteDatabase()
	if err != nil {
		return err
	}
	return db.Open()
}
