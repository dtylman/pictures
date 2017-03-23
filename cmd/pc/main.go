package main

import (
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/presentation"
	"github.com/dtylman/pictures/presentation/indexview"
	"github.com/dtylman/pictures/webkit"
)

func initViews() {
	indexview.OnIndexerStopped()
	indexview.InputMapQuest.SetValue(conf.Options.MapQuestAPIKey)
	for _, path := range conf.Options.SourceFolders {
		indexview.SourceFolders.AddElement(indexview.NewSourceFolder(path).Element)
	}
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
	err = webkit.Run(presentation.MainView)
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
