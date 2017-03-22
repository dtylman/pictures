package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/dtylman/pictures/backuper"
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/indexer"
	"github.com/dtylman/pictures/indexer/db"
	"github.com/dtylman/pictures/lockfile"
	"github.com/dtylman/pictures/server"
	"github.com/dtylman/pictures/server/route"
	"github.com/dtylman/pictures/server/route/middleware/httprouterhandler"
	"github.com/skratchdot/open-golang/open"
)

const nobrowser = "none"

func getLocalAddress() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", err
	}
	defer l.Close()
	return l.Addr().String(), nil
}

func openBrowser(url string, browser string) error {
	if browser == nobrowser {
		return nil
	} else if browser == "" {
		return open.Start(url)
	}
	return open.StartWith(url, browser)
}

func wait() {
	//todo think about "browser.IsRunning()" - to implement with searching for the PID
	for backuper.IsRunning() || indexer.IsRunning() || !httprouterhandler.LastAccess.Elapsed(time.Second*time.Duration(conf.Options.IdleSeconds)) {
		time.Sleep(time.Second * 5)
	}
}

func handleSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	s := <-ch
	log.Printf("Caught signal: %v", s)
	err := lockfile.Delete()
	if err != nil {
		log.Println(err)
	}
	os.Exit(0)
}

func startAndWait(address string, browser string) error {
	info, err := lockfile.Open()
	if err == nil {
		log.Println("Already running in ", info.URL())
		return openBrowser(info.URL(), browser)
	}
	if !os.IsNotExist(err) {
		return err
	}
	info = &lockfile.Info{Address: address, PID: os.Getpid()}
	err = info.Create()
	if err != nil {
		return err
	}
	go handleSignals()
	defer lockfile.Delete()
	err = conf.Load()
	if err != nil {
		return err
	}
	err = db.Open()
	if err != nil {
		return err
	}
	go func() {
		err := server.Start(info.Address, route.LoadHTTP())
		if err != nil {
			log.Fatal(err)
		}
	}()
	if browser == "" {
		err = open.Run(info.URL())
	} else if browser != nobrowser {
		err = open.RunWith(info.URL(), browser)
	}
	if err != nil {
		return err
	}
	wait()
	return nil
}

func main() {
	var err error
	browser := flag.String("browser", "", "the browser to invoke, if set to 'none', no browser will be invoked.")
	address := flag.String("bind", "", "bind to this address")
	flag.Parse()
	if *address == "" {
		*address, err = getLocalAddress()
		if err != nil {
			log.Fatal(err)
		}
	}
	err = startAndWait(*address, *browser)
	if err != nil {
		log.Fatal(err)
	}
}
