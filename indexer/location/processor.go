package location

import (
	"fmt"
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/indexer/picture"
	"github.com/dtylman/pictures/tasklog"
)

type Processor struct {
}

func NewProcessor() *Processor {
	return new(Processor)
}

func (p *Processor) Process(image *picture.Index) error {
	if !db.SetPhase(image.MD5, db.PhaseLocation) {
		return nil
	}
	err := PopulateLocation(image)
	if err != nil {
		tasklog.Error(err)
	} else {
		tasklog.StatusMessage(tasklog.IndexerTask, fmt.Sprintf("Found location  %s", image.Location))
	}
	return nil
}
