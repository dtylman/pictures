package db

import (
	"path"
	"os"
	"github.com/dtylman/pictures/conf"
	"path/filepath"
)

//ImageFolder creates and returns a folder name for an image on the database.
func ImageFolder(md5 string) (string, error) {
	base, err := conf.FilesPath()
	if err != nil {
		return "", err
	}
	size := len(md5)
	if size < 3 {
		// don't know what this is
		return base, nil
	}
	folder := filepath.Join(base, md5[0:2], md5[2:size])
	err = os.MkdirAll(folder, 0755)
	if err != nil {
		return "", err
	}
	return folder, err
}

//ImageFilePath creates folder for image, and returns a file path for a name on that folder.
func ImageFilePath(md5 string, fileName string) (string, error) {
	base, err := ImageFolder(md5)
	if err != nil {
		return "", err
	}
	return path.Join(base, fileName), nil
}