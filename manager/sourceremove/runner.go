package sourceremove

import (
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/indexer/picture"
	"log"
	"os"
	"path/filepath"
	"github.com/dtylman/pictures/tasklog"
	"fmt"
)

func run(path string) {
	if path == "" {
		return
	}
	log.Println("Starting source remover on ", path)
	tasklog.StatusMessage(tasklog.ManagerTask, fmt.Sprintf("Starting source remover on '%s'", path))

	defer func() {
		log.Printf("source remover on '%s' ended", path)
		tasklog.StatusMessage(tasklog.ManagerTask, "Done")
		Stop()
	}()

	err := filepath.Walk(path, processFile)
	if err != nil {
		log.Println(err)
	}

}

func processFile(path string, info os.FileInfo, e1 error) error {
	if e1 != nil {
		log.Println(e1)
		return nil
	}
	if info.IsDir() {
		return nil
	}

	image, err := picture.NewIndex(path, info)
	if err != nil {
		log.Println(err)
		return nil
	}

	if db.HasImage(image.MD5) {
		log.Printf("'%s' is already indexed, deleting it from source", image.Path)
		err = os.Remove(image.Path)
		if err != nil {
			log.Printf("Failed to remove '%s': %v", image.Path, err)
		} else {
			tasklog.StatusMessage(tasklog.ManagerTask, fmt.Sprintf("%s deleted",image.Path))
		}
	}
	return nil
}
