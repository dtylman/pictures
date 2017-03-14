package indexer

import (
	"errors"
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/indexer/picture"
	"github.com/dtylman/pictures/indexer/remover"
	"github.com/dtylman/pictures/indexer/thumbs"
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

func Stop() {
	indexer.SetDone()
}

func IsRunning() bool {
	return indexer.IsRunning()
}

func GetProgress() string {
	return indexer.ProgressString()
}

func indexOne(path string, info os.FileInfo, e1 error) error {
	if e1 != nil {
		indexer.AddError(path, e1)
		return nil
	}
	if !indexer.IsRunning() {
		return errors.New("Indexer had stopped")
	}
	if !info.IsDir() {
		if info.Size() > 0 {
			i, err := picture.NewIndex(path, info)
			if err != nil {
				indexer.AddError(path, err)
			} else {
				saveIndex(path, i)
			}
		}
		indexer.progress.incFile(info.Size())
	}
	return nil
}

func index(rootPath string) error {
	err := indexer.progress.init(rootPath)
	if err != nil {
		log.Println(err)
	}
	return filepath.Walk(rootPath, indexOne)
}

func indexPictures() {
	log.Println("Indexing ")
	defer indexer.SetDone()
	for _, folder := range conf.Options.SourceFolders {
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

func saveIndex(path string, i *picture.Index) {
	if indexer.GetOptions().IndexLocation {
		err := i.PopulateLocation()
		if err != nil {
			indexer.AddError(path, err)
		}
	}
	err := db.Index(i)
	if err != nil {
		indexer.AddError(path, err)
	}
	_, err = thumbs.MakeThumb(i.Path, i.MD5, true)
	if err != nil {
		indexer.AddError(path, err)
	}
}
