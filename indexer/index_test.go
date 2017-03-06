package indexer

import (
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/db"
	"github.com/dtylman/pictures/indexer/runningindexer"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	conf.Options.SourceFolders = []string{"/home/danny/Pictures"}
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
	sr, err := db.QueryAll(0, 100)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sr.String())
	for _, facet := range sr.Facets {
		t.Log(facet)
	}
}
