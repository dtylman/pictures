package location

import (
	"github.com/dtylman/pictures/indexer/picture"
	"github.com/dtylman/pictures/tasklog"
	"fmt"
)

type Processor struct {
}

func NewProcessor() *Processor {
	return new(Processor)
}

func (p*Processor) Process(image*picture.Index) error {
	if image.HasPhase(picture.PhaseLocation) {
		return nil
	}
	defer image.SetPhase(picture.PhaseLocation)

	err := PopulateLocation(image)
	if err != nil {
		tasklog.Error(err)
	} else {
		tasklog.StatusMessage(tasklog.IndexerTask, fmt.Sprintf("Found location  %s", image.Location))
	}
	return nil
}