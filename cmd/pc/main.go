package main

import (
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/view"
	"github.com/dtylman/pictures/webkit"
)

func initViews() {
	view.OnConfigChanged()
	view.OnIndexerStopped()
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
	initViews()
	err = webkit.Run(view.RootElement())
	if err != nil {
		return err
	}
	return err
}

func main() {
	err := run()
	if err != nil {
		webkit.Error(err)
	}
}
