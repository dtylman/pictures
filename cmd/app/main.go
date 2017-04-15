package main

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/pictures/cmd/app/view"
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/indexer/db"
	"net/http"
	_ "net/http/pprof"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"log"
)

func init() {
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()
	log.SetOutput(&lumberjack.Logger{
		Filename:   "bome.log",
		MaxSize:    1, // megabytes
		MaxBackups: 1,
		MaxAge:     7, //days
	})
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
