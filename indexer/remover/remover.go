package remover

import (
	"fmt"
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/indexer/picture"
	"github.com/dtylman/pictures/tasklog"
	"os"
)

type scanner struct {
	itemsToRemove []string
}

func (s *scanner) Remove() error {
	db.WalkImages(s.checkImage)
	return db.Remove(s.itemsToRemove)
}

func (s *scanner) checkImage(key string, image *picture.Index, err error) {
	if err != nil {
		tasklog.Error(err)
		return
	}
	if image == nil {
		tasklog.Println("error! image is null for ", key)
		return
	}
	_, err = os.Stat(image.Path)
	if err != nil {
		tasklog.Println(fmt.Sprintf("image: %v, error: %v", image.Path, err))
		s.itemsToRemove = append(s.itemsToRemove, image.MD5)
	}
}

func newScanner() *scanner {
	s := new(scanner)
	s.itemsToRemove = make([]string, 0)
	return s
}

//Remove removes from db images that no longer exists on disk
func Remove() error {
	s := newScanner()
	return s.Remove()
}
