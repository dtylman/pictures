package sourceremove

import (
	"errors"
	"sync"
)

var (
	running bool
	mutex   sync.Mutex
)

//IsRunning returns true if the backuper process is now running
func IsRunning() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return running
}

func Start(path string) error {
	mutex.Lock()
	defer mutex.Unlock()
	if running {
		return errors.New("Already running")
	}
	running = true
	go run(path)
	return nil
}

func Stop() {
	mutex.Lock()
	defer mutex.Unlock()
	running = false
}
