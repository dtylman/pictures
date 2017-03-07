package indexer

import (
	"sync"
)

type Progress struct {
	RootDir    string
	TotalFiles int
	TotalSize  int64
	mutex      sync.Mutex
}

func (p *Progress) init(path string) error {
	p.reset()
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.RootDir = path
	return nil
}

func (p *Progress) reset() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.TotalSize = 0
	p.TotalFiles = 0
}

func (p *Progress) incFile(size int64) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.TotalFiles++
	p.TotalSize += size
}
