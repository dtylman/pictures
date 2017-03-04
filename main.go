package main

import (
	"fmt"

	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/db"
	"github.com/dtylman/pictures/picture"
	"github.com/dtylman/pictures/server"
	"github.com/dtylman/pictures/server/route"
	"github.com/jasonwinn/geocoder"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"log"
	"os"
	"path/filepath"
)

func indexOne(path string, info os.FileInfo, _ error) error {
	if !info.IsDir() {
		if info.Size() > 0 {
			i, err := picture.NewIndex(path, info)
			if err != nil {
				log.Println(err)
			} else {
				err = db.Index(i)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
	return nil
}

func index(rootPath string) error {
	return filepath.Walk(rootPath, indexOne)
}

func indexPictures() {
	geocoder.SetAPIKey("8cCGEGGioKhpCLPjhAG44NfXYaXs9jCk")
	exif.RegisterParsers(mknote.All...)
	err := index("/home/danny/Pictures")
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	err := conf.Load()
	if err != nil {
		log.Fatal(err)
	}
	err = server.Start(route.LoadHTTP())
	if err != nil {
		log.Fatal(err)
	}
}
