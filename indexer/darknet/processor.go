package darknet

import (
	"errors"
	"fmt"
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/indexer/picture"
	"github.com/dtylman/pictures/tasklog"
	"time"
)

type Processor struct {
	darknet *Process
}

func NewProcessor() *Processor {
	p := new(Processor)
	p.spawnDarknet()
	return p
}

func (p *Processor) Process(image *picture.Index) error {
	if !db.SetPhase(image.MD5, db.PhaseObjects) {
		return nil
	}
	tasklog.StatusMessage(tasklog.IndexerTask, fmt.Sprintf("Detecing objects for %s", image.Path))
	res, err := p.darknet.Detect(image.Path, time.Duration(conf.Options.DarknetTimeout)*time.Second)
	if err != nil {
		tasklog.Error(err)
		p.spawnDarknet()
		return err
	}
	if res.Result != Success {
		return errors.New(res.Result)
	}
	for _, o := range res.Objects {
		image.Objects += fmt.Sprintf("%d %s with %d %% ", o.Count, o.Name, o.Prob)
	}
	tasklog.StatusMessage(tasklog.IndexerTask, fmt.Sprintf("Found %v", res.Objects))
	return nil
}

func (p *Processor) spawnDarknet() {
	if p.darknet != nil {
		p.darknet.Close()
	}
	var err error
	p.darknet, err = NewProcess()
	if err != nil {
		tasklog.Error(err)
	}
}
