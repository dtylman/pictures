package db

import (
	"encoding/json"

	"github.com/blevesearch/bleve"
	"github.com/boltdb/bolt"
	"github.com/dtylman/pictures/indexer/picture"
)

//Index saves one picture into the database
func Index1(picture *picture.Index) error {
	err := idx.Index(picture.MD5, picture)
	if err != nil {
		return err
	}
	return bdb.Update(func(tx *bolt.Tx) error {
		data, err := json.Marshal(picture)
		if err != nil {
			return err
		}
		return tx.Bucket(imagesBucket).Put([]byte(picture.MD5), data)
	})
}

//Search performs bleve search on the pictures index
func Search(req *bleve.SearchRequest) (*bleve.SearchResult, error) {
	return idx.Search(req)
}


