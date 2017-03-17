package indexer

import (
	"errors"
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/indexer/picture"
	"github.com/dtylman/pictures/indexer/remover"
	"github.com/dtylman/pictures/indexer/thumbs"
	"github.com/dtylman/pictures/progressbar"
	"github.com/jasonwinn/geocoder"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"log"
	"os"
	"path/filepath"
)

var indexer Indexer

func init() {
	exif.RegisterParsers(mknote.All...)
}

//Start starts the indexer
func Start(options Options) error {
	log.Println("Starting index with options: ", options)
	err := indexer.SetStarted(options)
	if err != nil {
		return err
	}
	if options.IndexLocation == true {
		if conf.Options.MapQuestAPIKey == "" {
			return errors.New("API KEY for map quest is empty")
		}
		geocoder.SetAPIKey(conf.Options.MapQuestAPIKey)
	}
	go indexPictures()
	return nil
}

//Stop stops the indexer
func Stop() {
	indexer.setDone()
}

//IsRunning returns true if indexer is running
func IsRunning() bool {
	return indexer.IsRunning()
}

//GetProgressStatus returns status for a progress bar
func GetProgressStatus() *progressbar.Status {
	return indexer.ProgressStatus()
}

func indexOne(path string, info os.FileInfo, e1 error) error {
	if e1 != nil {
		indexer.AddError(path, e1)
		return nil
	}
	if !indexer.IsRunning() {
		return errors.New("Indexer had stopped")
	}
	if info.IsDir() {
		indexer.setCurrentFolder(path)
	} else {
		if info.Size() > 0 {
			i, err := picture.NewIndex(path, info)
			if err != nil {
				indexer.AddError(path, err)
			} else {
				saveIndex(i)
			}
		}
		indexer.incFile(info.Size())
	}
	return nil
}

func index(rootPath string) error {
	return filepath.Walk(rootPath, indexOne)
}

func indexPictures() {
	log.Println("Indexing ")
	defer indexer.setDone()
	for _, folder := range conf.Options.SourceFolders {
		indexer.incSourceFolder()
		err := index(folder)
		if err != nil {
			indexer.AddError(folder, err)
		}
	}
	err := remover.Remove()
	if err != nil {
		log.Println(err)
	}
}

func saveIndex(newIndex *picture.Index) {
	if !indexer.GetOptions().ReIndex {
		if db.HasImage(newIndex.MD5) {
			return
		}
	}
	if indexer.GetOptions().IndexLocation {
		err := newIndex.PopulateLocation()
		if err != nil {
			indexer.AddError(newIndex.Path, err)
		}
	}
	err := db.Index(newIndex)
	if err != nil {
		indexer.AddError(newIndex.Path, err)
	}
	_, err = thumbs.MakeThumb(newIndex.Path, newIndex.MD5, true)
	if err != nil {
		indexer.AddError(newIndex.Path, err)
	}
}
