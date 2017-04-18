package indexer

import (
	"github.com/dtylman/pictures/indexer/picture"
	"github.com/dtylman/pictures/indexer/db"
	"time"
)

type processorBatch struct {
	images     []*picture.Index
	commitTime time.Time
}

func newProcessorBatch() *processorBatch {
	pb := new(processorBatch)
	pb.reset()
	return pb
}

func (pb*processorBatch) reset() {
	pb.images = make([]*picture.Index, 0)
	pb.commitTime = time.Now()
}

func (pb*processorBatch) add(image *picture.Index) error {
	pb.images = append(pb.images, image)
	if time.Since(pb.commitTime) > time.Second * 150 {
		return pb.commit()
	}
	return nil
}

func (pb*processorBatch) commit() error {
	err := db.BatchIndex(pb.images)
	if err != nil {
		return err
	}
	pb.reset()
	return nil
}