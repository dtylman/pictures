package backuper

import (
	"errors"
	"sync"
)

var (
	backupRunner runner
	mutex        sync.Mutex
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
	backupRunner.Running = true
	go backupRunner.run()
	return nil
}

func Stop() {
	mutex.Lock()
	defer mutex.Unlock()
	backupRunner.stop()
}
