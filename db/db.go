package db

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/document"
	"github.com/blevesearch/bleve/search/query"
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/indexer/picture"
	"github.com/syndtr/goleveldb/leveldb/errors"
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
	log.Println(picture.Path)
	return idx.Index(picture.MD5, picture)
}

func QueryAll(from int, size int) (*bleve.SearchResult, error) {
	return Query(bleve.NewMatchAllQuery(), size, from)
}

func Query(q query.Query, from int, size int) (*bleve.SearchResult, error) {
	search := bleve.NewSearchRequestOptions(q, size, from, false)
	return idx.Search(search)
}

func PathForImage(imageID string) (string, error) {
	doc, err := idx.Document(imageID)
	if err != nil {
		return "", err
	}
	for _, field := range doc.Fields {
		if field.Name() == "path" {
			return string(field.Value()), nil
		}
	}
	return "", errors.New("Image do not have a path field")
}

func GetImageDocument(imageID string) (*document.Document, error) {
	return idx.Document(imageID)
}
