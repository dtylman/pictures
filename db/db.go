package db

import (
	"github.com/blevesearch/bleve"
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/indexer/picture"
	"log"
	"os"
)

var idx bleve.Index

func init() {
	path, err := conf.DBPath()
	if err != nil {
		log.Fatal(err)
	}
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			idx, err = bleve.New(path, bleve.NewIndexMapping())
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}
	idx, err = bleve.Open(path)
	if err != nil {
		log.Fatal(err)
	}

}

func Index(picture *picture.Index) error {
	return idx.Index(picture.MD5, picture)
}
