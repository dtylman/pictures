package backuper

import (
	"github.com/dtylman/pictures/indexer/picture"
	"path/filepath"
	"github.com/dtylman/pictures/conf"
	"os"
)

type backupItem struct {
	Sources []string `json:"sources"`
	Target  string `json:"target"`
}

type backupItems map[string]*backupItem

func (b backupItems) Add(image *picture.Index) error {
	var err error
	item, exists := b[image.MD5]
	if exists {
		item.Sources = append(b[image.MD5].Sources, image.Path)
	} else {
		item = &backupItem{Sources:[]string{image.Path}}
		item.Target, err = targetFileFor(image.MD5)
	}
	b[image.MD5] = item
	return err
}

func targetFileFor(md5 string) (string, error) {
	folder := filepath.Join(conf.Options.BackupFolder, md5[0:2])
	err := os.MkdirAll(folder, 0755)
	if err != nil {
		return "", err
	}
	size := len(md5)
	return filepath.Join(folder, md5[2:size]), nil

}