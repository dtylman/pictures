package db

import (
	"encoding/json"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/document"
	"github.com/blevesearch/bleve/search/query"
	"github.com/boltdb/bolt"
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/indexer/picture"
	"log"
	"os"
)

var (
	idx bleve.Index
	bdb *bolt.DB
)

func init() {
	path, err := conf.BleveFolder()
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
	boltPath, err := conf.BoltPath()
	if err != nil {
		log.Fatal(err)
	}
	bdb, err = bolt.Open(boltPath, 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = bdb.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("images"))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

}

func Index(picture *picture.Index) error {
	log.Println("Indexing ", picture.Path)
	err := idx.Index(picture.MD5, picture)
	if err != nil {
		return err
	}
	return bdb.Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(picture)
		if err != nil {
			return err
		}
		return tx.Bucket([]byte("images")).Put([]byte(picture.MD5), data)
	})
}

func QueryAll(from int, size int) (*bleve.SearchResult, error) {
	return Query(bleve.NewMatchAllQuery(), size, from)
}

func Query(q query.Query, from int, size int) (*bleve.SearchResult, error) {
	search := bleve.NewSearchRequestOptions(q, size, from, false)
	return idx.Search(search)
}

//GetImage gets image info by image id
func GetImage(imageID string) (*picture.Index, error) {
	index := new(picture.Index)
	return index, bdb.View(func(tx *bolt.Tx) error {
		data := tx.Bucket([]byte("images")).Get([]byte(imageID))
		return json.Unmarshal(data, index)
	})
}

//GetImageDocument get the indexed document from bleve
func GetImageDocument(imageID string) (*document.Document, error) {
	return idx.Document(imageID)
}
