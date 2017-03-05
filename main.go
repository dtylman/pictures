package main

import (
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/server"
	"github.com/dtylman/pictures/server/route"
	"log"
)

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
