package db

import (
	"github.com/blevesearch/bleve"
	"github.com/boltdb/bolt"
	"github.com/dtylman/pictures/conf"
	"os"
)

var (
	idx          bleve.Index
	bdb          *bolt.DB
	imagesBucket = []byte("images")
)

func openBleve() error {
	path, err := conf.BlevePath()
	if err != nil {
		return err
	}
	_, err = os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			mapping := bleve.NewIndexMapping()
			err = mapping.Validate()
			if err != nil {
				return err
			}
			idx, err = bleve.New(path, mapping)
			if err != nil {
				return err
			}
			return nil
		}
	}
	idx, err = bleve.Open(path)
	return err
}

func openBolt() error {
	boltPath, err := conf.BoltPath()
	if err != nil {
		return err
	}
	bdb, err = bolt.Open(boltPath, 0644, nil)
	if err != nil {
		return err
	}
	return bdb.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("images"))
		return err
	})
}

//Open opens or creates the local database
func Open() error {
	err := openBleve()
	if err != nil {
		return err
	}
	return openBolt()
}
