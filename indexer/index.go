package indexer

import (
	"github.com/dtylman/pictures/indexer/darknet"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"log"
	"os"
)

func init() {
	exif.RegisterParsers(mknote.All...)
	darknet.DarknetHome, _ = os.Getwd()
}

type Options struct {
	//IndexLocation if true will do include geolocation
	WithLocation bool
	//DeleteDatabase if true will delete previous results
	DeleteDatabase bool
	//WithObjects if true will include objects
	WithObjects bool
	//With faces if true will include faces
	WithFaces bool
	//QuickScan if true, will only compare file name when looking at image
	QuickScan bool
}

var (
	indexer Indexer
)

func IsRunning() bool {
	return indexer.isRunning()
}

//Stop stops the indexer
func Stop() {
	indexer.setRunning(false)
}

//Start starts the indexer
func Start(options Options) error {
	return indexer.start(options)
}

func AddError(path string, err error) {
	log.Printf("%s: %s", path, err.Error())
}
