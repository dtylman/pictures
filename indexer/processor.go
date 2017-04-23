package indexer

import (
	"fmt"
	"runtime"
	"sync"

	"reflect"

	"log"

	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/indexer/darknet"
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/indexer/location"
	"github.com/dtylman/pictures/indexer/picture"
	"github.com/dtylman/pictures/indexer/thumbs"
	"github.com/dtylman/pictures/tasklog"
	"github.com/jasonwinn/geocoder"
)

type Processor struct {
	images     *picture.Queue
	options    Options
	processors []picture.Processor
}

func newProcessor(options Options) *Processor {
	p := &Processor{
		images:     picture.NewQueue(),
		options:    options,
		processors: make([]picture.Processor, 0),
	}
	db.WalkImages(p.walkImage)
	proc := thumbs.NewProcessor()
	proc.Overwrite = options.DeleteDatabase
	p.processors = append(p.processors, proc)

	if options.WithLocation {
		geocoder.SetAPIKey(conf.Options.MapQuestAPIKey)
		p.processors = append(p.processors, location.NewProcessor())
	}

	if options.WithObjects {
		p.processors = append(p.processors, darknet.NewProcessor())

	}

	return p
}

func (p *Processor) walkImage(key string, image *picture.Index, err error) {
	p.images.PushBack(image)
}

func (p *Processor) worker(wg *sync.WaitGroup, total int) {
	log.Printf("Starting worker (%v)", runtime.NumCPU())
	defer wg.Done()
	var dp *darknet.Process
	var err error

	defer func() {
		if dp != nil {
			dp.Close()
		}
	}()
	left, image := p.images.Pop()
	batch := newProcessorBatch()
	for image != nil {
		before := *image
		tasklog.Status(tasklog.IndexerTask, IsRunning(), total-left, total, fmt.Sprintf("Processing %s...", image.Path))
		for _, processor := range p.processors {
			if !IsRunning() {
				//indexer had stopped.
				return
			}
			err = processor.Process(image)
			if err != nil {
				tasklog.Error(err)
			}
		}
		if !reflect.DeepEqual(before, *image) {
			batch.add(image)
		}
		left, image = p.images.Pop()
	}
	tasklog.Status(tasklog.IndexerTask, IsRunning(), total-left, total, "Saving images...")
	batch.commit()
}

func (p *Processor) update() {
	total := p.images.Length()
	for p.images.Length() > 0 {
		waitGroup := new(sync.WaitGroup)
		for i := 0; i < runtime.NumCPU(); i++ {
			waitGroup.Add(1)
			go p.worker(waitGroup, total)
		}
		waitGroup.Wait()
	}
}
