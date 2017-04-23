package db

import (
	"os"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
	"github.com/boltdb/bolt"
	"github.com/dtylman/pictures/conf"
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
			mapping := bleveMapping()
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

func bleveMapping() mapping.IndexMapping {
	mapping := bleve.NewIndexMapping()
	return mapping
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

//Close closes the db
func Close() {
	idx.Close()
	bdb.Close()
}

//DeleteDatabase removes the database and all data
func DeleteDatabase() error {
	err := idx.Close()
	if err != nil {
		return err
	}
	err = bdb.Close()
	if err != nil {
		return err
	}
	path, err := conf.BlevePath()
	if err != nil {
		return err
	}
	err = os.RemoveAll(path)
	if err != nil {
		return err
	}
	path, err = conf.BoltPath()
	if err != nil {
		return err
	}
	err = os.RemoveAll(path)
	if err != nil {
		return err
	}
	path, err = conf.FilesPath()
	if err != nil {
		return err
	}
	return os.RemoveAll(path)
}
