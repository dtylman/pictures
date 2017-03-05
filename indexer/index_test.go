package indexer

import (
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/indexer/runningindexer"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	t.Log(conf.Options.SourceFolders)
	err := Start(runningindexer.Options{IndexLocation: false, ReIndex: false})
	assert.NoError(t, err)
	assert.True(t, runningindexer.IsRunning())
	err = Start(runningindexer.Options{IndexLocation: false, ReIndex: false})
	assert.Error(t, err)
	for runningindexer.IsRunning() {
		time.Sleep(time.Second)
		t.Log("Still running")
	}
	t.Log("Done")
}
