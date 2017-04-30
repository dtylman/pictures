package backuper

import (
	"sync"
	"errors"
)

var (
	backupRunner runner
	mutex sync.Mutex
)

//IsRunning returns true if the backuper process is now running
func IsRunning() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return backupRunner.Running
}

func Start() error {
	mutex.Lock()
	defer mutex.Unlock()
	if backupRunner.Running {
		return errors.New("Already running")
	}
	go backupRunner.run()
	return nil
}

func Stop()  {
	mutex.Lock()
	defer mutex.Unlock()
	backupRunner.stop()
}
