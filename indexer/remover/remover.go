package remover

import (
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/tasklog"
	"os"
)

type scanner struct {
	missingFiles []string
}

func (s *scanner) Remove() error {
	s.missingFiles = make([]string,0)
	db.WalkFiles(s.checkFile)
	db.RemoveFiles(s.missingFiles)
	return nil
}

func (s *scanner) checkFile(path string, e1 error) {
	if e1 != nil {
		tasklog.Error(e1)
		return
	}
	if path == "" {
		return
	}
	fileInfo, err := os.Stat(path)
	if fileInfo.IsDir(){
		return
	}
	if err != nil {
		tasklog.ErrorF("path: %v, error: %v", path, err)
		s.missingFiles=append(s.missingFiles,path)
	}
}

func newScanner() *scanner {
	s := new(scanner)
	return s
}

//Remove removes from db images that no longer exists on disk
func Remove() error {
	s := newScanner()
	return s.Remove()
}
