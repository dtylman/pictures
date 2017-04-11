package indexer

import (
	"sync"
	"github.com/dtylman/pictures/tasklog"
	"github.com/dtylman/pictures/indexer/remover"
	"github.com/dtylman/pictures/conf"
	"os"
	"github.com/dtylman/pictures/indexer/picture"
	"errors"
	"github.com/dtylman/pictures/indexer/db"
	"path/filepath"
	"fmt"
	"github.com/jasonwinn/geocoder"
)

type Walker struct {
	options Options
	running bool
	mutex   sync.Mutex
}

func (w*Walker) isRunning() bool {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.running
}

func (w*Walker) setRunning(value bool) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.running = value
}

func (w*Walker) indexPictures() {
	defer func() {
		w.setRunning(false)
		tasklog.Status(tasklog.IndexerTask, false, 0, 0, "Done")
	}()
	tasklog.Println("Starting index with options: ", w.options)
	if w.options.DeleteDatabase {
		tasklog.StatusMessage(tasklog.IndexerTask, "Deleting existing database...")
		err := w.deleteDB()
		if err != nil {
			tasklog.StatusMessage(tasklog.IndexerTask, err.Error())
		}
	}
	tasklog.StatusMessage(tasklog.IndexerTask, "Indexing files...")
	for _, folder := range conf.Options.SourceFolders {
		err := filepath.Walk(folder, w.indexOne)
		if err != nil {
			AddError(folder, err)
		}
	}
	tasklog.StatusMessage(tasklog.IndexerTask, "Updating indicies...")
	err := remover.Remove()
	if err != nil {
		tasklog.Println(err)
	}

	imagePopulator := NewUpdater(w.options)
	imagePopulator.update()

}

func (w*Walker) start(options Options) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.running {
		return errors.New("Already running")
	}
	w.running = true
	w.options = options
	if w.options.WithLocation {
		if conf.Options.MapQuestAPIKey == "" {
			return errors.New("API KEY for map quest is empty")
		}
		geocoder.SetAPIKey(conf.Options.MapQuestAPIKey)
	}
	go walker.indexPictures()
	return nil
}

func (w*Walker) indexOne(path string, info os.FileInfo, e1 error) error {
	if e1 != nil {
		AddError(path, e1)
		return nil
	}
	if !IsRunning() {
		return errors.New("Indexer had stopped")
	}
	if !info.IsDir() {
		if info.Size() > 0 {
			i, err := picture.NewIndex(path, info)
			if err != nil {
				AddError(path, err)
				return nil
			}
			if !w.options.DeleteDatabase {
				if db.HasImage(i.MD5) {
					return nil
				}
			}
			tasklog.StatusMessage(tasklog.IndexerTask, fmt.Sprintf("Indexing %s", i.Path))
			err = db.Index(i)
			if err != nil {
				AddError(path, err)
			}
		}
	}
	return nil
}

func (w*Walker) deleteDB() error {
	err := db.DeleteDatabase()
	if err != nil {
		return err
	}
	return db.Open()
}
