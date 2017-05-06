package remover

import (
	"os"

	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/tasklog"
)

type scanner struct {
	missingFiles []string
}

func (s *scanner) Remove() error {
	s.missingFiles = make([]string, 0)
	db.WalkFiles(s.checkFile)
	db.RemoveFiles(s.missingFiles)
	return nil
}

func (s *scanner) checkFile(path string, e1 error) error {
	if e1 != nil {
		tasklog.Error(e1)
		return nil
	}
	if path == "" {
		return nil
	}
	_, err := os.Stat(path)
	if err != nil {
		tasklog.ErrorF("path: %v, error: %v", path, err)
		s.missingFiles = append(s.missingFiles, path)
	}
	return nil
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
