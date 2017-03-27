package indexer

import (
	"fmt"
	"github.com/c2h5oh/datasize"
)

type IndexerProgress struct {
	Running           bool
	CurrentFolder     string
	TotalFiles        int
	TotalSize         int64
	TotalRootFolders  int
	CurrentRootFolder int
	Errors            []IndexError
}

func (ip *IndexerProgress) Percentage() int {
	return ip.CurrentRootFolder / ip.TotalRootFolders * 100
}

func (i *IndexerProgress) Text() string {
	return fmt.Sprintf("Running: %v %s: %v files (%s)", i.Running, i.CurrentFolder, i.TotalFiles, datasize.ByteSize(i.TotalSize).HumanReadable())
}
