package backuper

import (
	"github.com/dtylman/pictures/indexer/picture"
	"github.com/dtylman/pictures/indexer/db"
	"log"
	"encoding/json"
	"path/filepath"
	"github.com/dtylman/pictures/conf"
	"io/ioutil"
	"github.com/dtylman/pictures/tasklog"
	"os"
	"io"
	"fmt"
)

type runner struct {
	Running bool
	items   backupItems
}

func (r*runner) checkImage(key string, image *picture.Index, err error) {
	if err != nil {
		log.Println(err)
		return
	}
	err = r.items.Add(image)
	if err != nil {
		log.Println(err)
	}
}

func (r*runner) fileExists(src, dest string) bool {
	in, err := os.Stat(src)
	if err != nil {
		return false
	}
	out, err := os.Stat(dest)
	if err != nil {
		return false
	}
	return in.Size() == out.Size()

}

func (r*runner) copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func (r*runner) copyFiles() error {
	total := len(r.items)
	i := 0
	for _, item := range r.items {
		fileName := item.Sources[0]
		tasklog.Status(tasklog.BackuperTask, true, i, total, fmt.Sprintf("Copying %s", fileName))
		if r.fileExists(fileName, item.Target) {
			continue
		}
		err := r.copyFile(fileName, item.Target)
		if err != nil {
			return err
		}
		i++
	}
	data, err := json.MarshalIndent(r.items, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(conf.Options.BackupFolder, "bome.backup.json"), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (r*runner) run() {
	defer func() {
		Stop()
		log.Println("Backup finished")
		tasklog.Status(tasklog.BackuperTask, false, 0, 0, "Done")
	}()
	log.Printf("Starting backup to %s", conf.Options.BackupFolder)
	tasklog.Status(tasklog.BackuperTask, true, 0, 0, "Backup started...")
	r.items = make(backupItems)
	db.WalkImages(r.checkImage)
	err := r.copyFiles()
	if err != nil {
		log.Println(err)
	}
}

func (r*runner) stop() {
	r.Running = false
}

