package backuper

import "github.com/dtylman/pictures/progressbar"

//IsRunning returns true if the backuper process is now running
func IsRunning() bool {
	return false
}

var Status progressbar.Status
