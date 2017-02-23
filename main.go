package main

import (
	"bitbucket.org/taruti/mimemagic"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func walky(path string, info os.FileInfo, _ error) error {
	if !info.IsDir() {
		if info.Size() > 0 {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			sig := make([]byte, 1024)
			_, err = file.Read(sig)
			if err != nil {
				return err
			}
			mimeType := mimemagic.Match("", sig)
			mimeType = strings.Split(mimeType, "/")[0]
			if mimeType == "image" || mimeType == "video" {
				_, err = file.Seek(0, 0)
				if err != nil {
					return err
				}
				x, err := exif.Decode(file)
				if err != nil {
					log.Println(path, err)
				} else {
					time, err := x.DateTime()
					if err == nil {
						log.Println(time.Year(), time.Month())
					}
					lat, long, err := x.LatLong()
					if err != nil {
						log.Println(path, err)
					} else {
						log.Println(path, lat, long)
					}
				}
			}
		}
	}
	return nil
}

func worky() error {
	return filepath.Walk("", walky)
}

func main() {
	exif.RegisterParsers(mknote.All...)
	err := worky()
	if err != nil {
		fmt.Println(err)
	}
}
