package db

import (
	"github.com/blevesearch/bleve"
	"github.com/dtylman/pictures/picture"
	"log"
)

var idx bleve.Index

func init() {
	var err error
	idx, err = bleve.New("db.bleve", bleve.NewIndexMapping())
	if err != nil {
		log.Println(err)
	}
}

func Index(picture *picture.Index) error {
	return idx.Index(picture.MD5, picture)
}
