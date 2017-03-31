package main

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/view"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()
}

func run() error {
	err := conf.Load()
	if err != nil {
		return err
	}
	err = db.Open()
	if err != nil {
		return err
	}
	defer db.Close()
	view.InitializeComponents()
	err = gowd.Run(view.Root.Element)
	if err != nil {
		return err
	}
	return err
}

func main() {
	err := run()
	if err != nil {
		gowd.Error(err)
	}
}
