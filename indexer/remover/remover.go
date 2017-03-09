package remover

import (
	"fmt"
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/indexer/picture"
	"log"
	"os"
)

type scanner struct {
	items []string
}

func (s *scanner) Remove() error {
	db.WalkImages(s.scanPicture)
	return db.Remove(s.items)
}

func (s *scanner) scanPicture(key string, image *picture.Index, err error) {
	if err != nil {
		log.Println(err)
		return
	}
	if image == nil {
		log.Println("error! image is null for ", key)
		return
	}
	_, err = os.Stat(image.Path)
	if err != nil {
		log.Println(fmt.Sprintf("image: %v, error: %v", image.Path, err))
		s.items = append(s.items, image.MD5)
	}
}

func newScanner() *scanner {
	s := new(scanner)
	s.items = make([]string, 0)
	return s
}

//Remove removes from db images that no longer exists on disk
func Remove() error {
	s := newScanner()
	return s.Remove()
}
