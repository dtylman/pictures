package main

import (
	"flag"
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/server"
	"github.com/dtylman/pictures/server/route"
	"github.com/nightlyone/lockfile"
	"log"
	"os"
	"path/filepath"
)

func main() {
	lf, err := lockfile.New(filepath.Join(os.TempDir(), "pictures.lock"))
	if err != nil {
		log.Fatal(err)
	}
	err = lf.TryLock()
	if err != nil {
		log.Fatal(err)
	}
	defer lf.Unlock()
	err = conf.Load()
	if err != nil {
		log.Fatal(err)
	}
	err = db.Open()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	browser := flag.String("browser", "google-chrome", "the name of the browser to invoke")
	address := flag.String("bind-address", "", "bind to this address")
	flag.Parse()
	err = server.Start(route.LoadHTTP(), *browser, *address)
	if err != nil {
		log.Fatal(err)
	}
}
