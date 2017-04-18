package db

import (
	"encoding/json"
	"github.com/blevesearch/bleve"
	"github.com/boltdb/bolt"
	"github.com/dtylman/pictures/indexer/picture"
	"fmt"
)

func BatchIndex(pictures []*picture.Index) error {
	b := idx.NewBatch()
	for _, picture := range pictures {
		err := b.Index(picture.MD5, picture)
		if err != nil {
			return err
		}
	}
	err := idx.Batch(b)
	if err != nil {
		return err
	}
	return bdb.Update(func(tx *bolt.Tx) error {
		for _, picture := range pictures {
			data, err := json.Marshal(picture)
			if err != nil {
				return err
			}
			err = tx.Bucket(imagesBucket).Put([]byte(picture.MD5), data)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

//Index saves one picture into the database
func Index(picture *picture.Index) error {
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

//HasPath returns true if a path exists in the database
func HasPath(path string) (bool, error) {
	q := bleve.NewTermQuery(fmt.Sprintf("path: %s", path))
	sr, err := idx.Search(bleve.NewSearchRequestOptions(q, 1, 0, false))
	if err != nil {
		return false, err
	}
	return sr.Hits.Len() > 0, nil
}

//HasImage returns tue if image wiith ID exists
func HasImage(imageID string) bool {
	var exists bool
	bdb.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(imagesBucket).Get([]byte(imageID))
		exists = data != nil
		return nil
	})
	return exists
}

type WalkImagesFunc func(key string, image *picture.Index, err error)

//WalkImages executes function for all images in the database
func WalkImages(wf WalkImagesFunc) {
	bdb.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(imagesBucket).Cursor()
		key, value := c.First()
		for key != nil && value != nil {
			key, value = c.Next()
			pic := new(picture.Index)
			err := json.Unmarshal(value, pic)
			wf(string(key), pic, err)
		}
		return nil
	})
}

//Remove removes all items in keys
func Remove(keys []string) error {
	err := bdb.Update(func(tx *bolt.Tx) error {
		for _, key := range keys {
			err := tx.Bucket(imagesBucket).Delete([]byte(key))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, key := range keys {
		err = idx.Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}
