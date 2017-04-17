package indexer

import (
	"github.com/dtylman/pictures/indexer/picture"
	"github.com/dtylman/pictures/tasklog"
	"fmt"
	"github.com/dtylman/pictures/indexer/db"
	"runtime"
	"sync"
	"github.com/dtylman/pictures/indexer/darknet"
	"github.com/dtylman/pictures/indexer/thumbs"
	"github.com/dtylman/pictures/indexer/location"
	"github.com/jasonwinn/geocoder"
	"github.com/dtylman/pictures/conf"
)

type Processor struct {
	images     *picture.Queue
	options    Options
	processors []picture.Processor
}

func newProcessor(options Options) *Processor {
	p := &Processor{
		images: picture.NewQueue(),
		options: options,
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

func (p*Processor) walkImage(key string, image *picture.Index, err error) {
	p.images.PushBack(image)
}

func (p*Processor) worker(wg*sync.WaitGroup, total int) {
	defer wg.Done()
	var dp *darknet.Process
	var err error

	defer func() {
		if dp != nil {
			dp.Close()
		}
	}()
	left, i := p.images.Pop()
	for (i != nil) {
		tasklog.Status(tasklog.IndexerTask, IsRunning(), total - left, total, fmt.Sprintf("Indexing %s", i.Path))
		for _, processor := range p.processors {
			if !IsRunning() {
				//indexer had stopped.
				return
			}
			err = processor.Process(i)
			if err != nil {
				tasklog.Error(err)
			}
		}
		err = db.Index(i)
		if err != nil {
			tasklog.Error(err)
		}
		left, i = p.images.Pop()
	}
}

func (p*Processor) update() {
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
