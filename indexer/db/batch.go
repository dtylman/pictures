package db

import (
	"encoding/json"

	"log"

	"runtime"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/dtylman/pictures/indexer/picture"
)

//DefaultBatchSize is the batch size to use when updaing pictures in batch
var DefaultBatchSize = 100

func batchIndex(pictures []*picture.Index) error {
	log.Printf("Batching %v", pictures)
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

func shardWorker(wg *sync.WaitGroup, images *picture.Queue) {
	log.Println("Starting shard worker")
	defer log.Println("Shard worker ended")
	defer wg.Done()
	batch := make([]*picture.Index, DefaultBatchSize)
	for images.Length() > 0 {
		i := 0
		for ; i < DefaultBatchSize; i++ {
			_, image := images.Pop()
			if image != nil {
				batch[i] = image
			} else {
				break
			}
		}
		err := batchIndex(batch[0:i])
		if err != nil {
			log.Println(err)
		}
	}
}

//BatchIndex updates batch of pictures
func BatchIndex1(images *picture.Queue) {
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go shardWorker(&wg, images)
	}
	wg.Wait()
}
