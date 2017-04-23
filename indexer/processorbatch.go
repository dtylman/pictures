package indexer

import (
	"time"

	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/indexer/picture"
)

type processorBatch struct {
	images     *picture.Queue
	commitTime time.Time
}

func newProcessorBatch() *processorBatch {
	pb := new(processorBatch)
	pb.reset()
	return pb
}

func (pb *processorBatch) reset() {
	pb.images = picture.NewQueue()
	pb.commitTime = time.Now()
}

func (pb *processorBatch) add(image *picture.Index) {
	pb.images.PushBack(image)
	if time.Since(pb.commitTime) > time.Second*150 {
		pb.commit()
	}
}

func (pb *processorBatch) commit() {
	db.BatchIndex(pb.images)
	pb.reset()
}
