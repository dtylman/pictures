package thumbs

import (
	"github.com/dtylman/pictures/indexer/picture"
	"github.com/dtylman/pictures/tasklog"
	"fmt"
	"github.com/dtylman/pictures/indexer/db"
)

type Processor struct {
	Overwrite bool
}

func NewProcessor() *Processor {
	return new(Processor)
}

func (p*Processor) Process(image *picture.Index) error {
	if !db.SetPhase(image.MD5, db.PhaseThumb) {
		return nil
	}
	tasklog.StatusMessage(tasklog.IndexerTask, fmt.Sprintf("Thumbing %s", image.Path))
	_, err := MakeThumb(image.Path, image.MD5, p.Overwrite)
	if err != nil {
		tasklog.Error(err)
	}
	return nil
}